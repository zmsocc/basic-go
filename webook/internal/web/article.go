package web

import (
	"fmt"
	intrv1 "gitee.com/zmsoc/gogogo/webook/api/proto/gen/intr/v1"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	ijwt "gitee.com/zmsoc/gogogo/webook/internal/web/jwt"
	"gitee.com/zmsoc/gogogo/webook/pkg/ginx"
	"gitee.com/zmsoc/gogogo/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc     service.ArticleService
	l       logger.LoggerV1
	intrSvc intrv1.InteractiveServiceClient
	biz     string
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1,
	// 你可以注入真的 grpc 客户端，也可以注入那个本地调用伪装的
	intrSvc intrv1.InteractiveServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		l:       l,
		biz:     "article",
		intrSvc: intrSvc,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	// 修改
	//g.PUT("/")
	// 新增
	//g.POST("/")

	g.POST("/edit", h.Edit)
	g.POST("/withdraw", h.Withdraw)
	g.POST("/publish", h.Publish)
	// 创作者的查询接口
	// 这个是获取数据的接口，理论上来说(遵循 RESTful 规范)，应该是用 GET 方法
	// GET localhost/articles => List 接口
	g.POST("/list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](h.List))
	g.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](h.Detail))
	g.GET("/:id", ginx.WrapToken[ijwt.UserClaims](h.PubDetail), func(ctx *gin.Context) {
		// // 增加阅读计数。
		//		//go func() {
		//		//	// 开一个 goroutine，异步去执行
		//		//	er := a.intrSvc.IncrReadCnt(ctx, a.biz, art.Id)
		//		//	if er != nil {
		//		//		a.l.Error("增加阅读计数失败",
		//		//			logger.Int64("aid", art.Id),
		//		//			logger.Error(err))
		//		//	}
		//		//}()
	})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context, uc ijwt.UserClaims) (Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return Result{Code: 4, Msg: "参数错误"}, fmt.Errorf("前端输入的 ID 不对 %w", err)
	}

	uc = ctx.MustGet("user").(ijwt.UserClaims)
	var eg errgroup.Group
	var art domain.Article
	eg.Go(func() error {
		art, err = h.svc.GetPublishedById(ctx, id, uc.Uid)
		return err
	})

	//var getResp *intrv1.GetResponse
	eg.Go(func() error {
		// 这个地方可以容忍错误
		_, err := h.intrSvc.Get(ctx, &intrv1.GetRequest{
			Biz: h.biz, BizId: id, Uid: uc.Uid,
		})
		// 这种是容错的写法
		// if err != nil {
		// 记录日志
		// }
		// return nil
		return err
	})
	// 在这儿等，要保证前面两个
	err = eg.Wait()
	if err != nil {
		// 代表查询出错了
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	// 增加阅读计数
	go func() {
		// 你都异步了，怎么还说有巨大的压力呢？
		// 开一个 goroutine，异步去执行
		_, er := h.intrSvc.IncrReadCnt(ctx, &intrv1.IncrReadCntRequest{
			Biz: h.biz, BizId: art.Id,
		})
		if er != nil {
			h.l.Error("增加阅读计数失败",
				logger.Int64("aid", art.Id),
				logger.Error(err))
		}
	}()

	// ctx.Set("art", art)
	//intr := getResp.Intr

	// 这个功能是不是可以让前端，主动发一个 HTTP 请求，来增加一个计数？
	return Result{
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 要把作者信息带出去
			Author: art.Author.Name,
			Ctime:  art.Ctime.Format(time.DateTime),
			Utime:  art.Utime.Format(time.DateTime),
		},
	}, nil
}

func (h *ArticleHandler) Detail(ctx *gin.Context, usr ijwt.UserClaims) (Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return ginx.Result{Code: 4, Msg: "参数错误"}, fmt.Errorf("前端输入的 ID 不对 %w", logger.Error(err))
	}
	usr, ok := ctx.MustGet("user").(ijwt.UserClaims)
	if !ok {
		return ginx.Result{Code: 5, Msg: "系统错误"}, fmt.Errorf("获得用户会话信息失败")
	}
	art, err := h.svc.GetById(ctx, id)
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, fmt.Errorf("获得文章信息失败 %w", logger.Error(err))
	}
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != usr.Uid {
		h.l.Error("非法访问文章，创作者 Id 不匹配", logger.Int64("uid", usr.Uid))
		return ginx.Result{Code: 4, Msg: ""}, nil
	}
	return ginx.Result{
		Data: ArticleVO{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: src.Author,
			Ctime: art.Ctime.Format(time.DateTime),
			Utime: art.Utime.Format(time.DateTime),
		},
	}, nil
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 要打日志
		h.l.Error("未发现用户的 session 信息")
		return
	}
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 要打日志
		h.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "发表帖子成功",
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 要打日志
		h.l.Error("未发现用户的 session 信息")
		return
	}
	// 检测输入，跳过这一步
	// 调用 service 代码
	err := h.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 要打日志
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "保存帖子成功",
	})
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 要打日志
		h.l.Error("未发现用户的 session 信息")
		return
	}

	// 检测输入，跳过这一步
	// 调用 service 代码
	id, err := h.svc.Save(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 要打日志
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "保存帖子成功",
		Data: id,
	})
}

func (h *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (ginx.Result, error) {
	res, err := h.svc.List(ctx, uc.Uid, req.Offset, req.Limit)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	// 在列表页，不显示全文，只显示一个"摘要"
	// 比如说，简单的摘要就是前几句话
	// 强大的摘要是 AI 帮你生成的
	return ginx.Result{
		Data: slice.Map[domain.Article, ArticleVO](res,
			func(idx int, src domain.Article) ArticleVO {
				return ArticleVO{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Status:   src.Status.ToUint8(),
					// 这个列表请求，不需要返回内容
					// Content: src.Content,
					// 这个是创作者看自己的文章列表，也不需要这个字段
					//Author: src.Author,
					Ctime: src.Ctime.Format(time.DateTime),
					Utime: src.Utime.Format(time.DateTime),
				}
			}),
	}, nil
}
