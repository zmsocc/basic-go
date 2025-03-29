package web

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	"gitee.com/zmsoc/gogogo/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSVC service.UserService
	jwtHandler
}

func NewOAuth2WechatHandler(svc wechat.Service, userSVC service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     svc,
		userSVC: userSVC,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登陆URL失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	info, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 从 userService 里面拿 uid
	u, err := h.userSVC.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	// 验证微信的code
}

//type OAuth2Handler struct {
//	wechatService string
//	svcs          map[string]OAuth2Service
//}
//
//func (h *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
//	// 统一处理所有的 OAuth2 的
//	g := server.Group("/oauth2")
//	g.GET("/:platform/authurl", h.AuthURL)
//	g.GET("/:platform/callback", h.Callback)
//}
//
//func (h *OAuth2Handler) AuthURL(ctx *gin.Context) {
//	platform := ctx.Param("platform")
//	switch platform {
//	case "wechat":
//		h.wechatService.AuthURL
//	}
//
//	svc := h.svcs
//}
//
//func (h *OAuth2Handler) Callback(ctx *gin.Context) {
//
//}
