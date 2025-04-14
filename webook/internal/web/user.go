package web

import (
	"fmt"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	ijwt "gitee.com/zmsoc/gogogo/webook/internal/web/jwt"
	"gitee.com/zmsoc/gogogo/webook/pkg/ginx"
	"gitee.com/zmsoc/gogogo/webook/pkg/logger"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const biz = "login"

// 确保 UserHandler 上实现了 handler 接口
var _ handler = &UserHandler{}

// 这个更优雅
var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	ijwt.Handler
	cmd redis.Cmdable
	l   logger.LoggerV1
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService,
	jwtHdl ijwt.Handler, l logger.LoggerV1) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		codeSvc:     codeSvc,
		Handler:     jwtHdl,
	}
}

//func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
//	ug.GET("/profile", u.Profile)
//	ug.POST("/signup", u.SignUp)
//	ug.POST("/login", u.Login)
//	ug.POST("/edit", u.Edit)
//}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", ginx.WrapBodyV1(u.LoginJWT))
	ug.POST("/logout", u.LogoutJWT)
	ug.POST("/edit", u.Edit)
	ug.POST("/login_sms/code/send", u.SendLoginSMACode)
	ug.POST("/login_sms", ginx.WrapBody[LoginSMSReq](u.l.With(logger.String("method", "login_sms")), u.LoginSMS))
	ug.POST("/refresh_token", u.RefreshToken)
}

// RefreshToken 可以同时刷新长短 token， 用 redis 来记录是否有效，即 refresh_token
// 参考登录校验部分，比较 User-Agent 来增强安全性
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 只有这个接口，拿出来的才是 refresh_token，其他地方都是 acess token
	refreshToken := u.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = u.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 要么 redis 有问题， 要么已经退出登录
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 搞个新的 access_token
	err = u.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		// 信息量不足，无效日志
		zap.L().Error("系统异常", zap.Error(err))
		// 正常来讲， msg 的部分就应该包含足够的定位信息
		zap.L().Error("设置 JWT token 出现异常",
			zap.Error(err),
			zap.String("服务", "UserHandler:RefreshToken"))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})

}

type LoginSMSReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (u *UserHandler) LoginSMS(ctx *gin.Context, req LoginSMSReq) (Result, error) {
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		//ctx.JSON(http.StatusOK, Result{
		//	Code: 5,
		//	Msg:  "系统错误",
		//})
		//zap.L().Error("校验验证码出错", zap.Error(err)) // 不能这样打，因为手机号码是敏感数据，你不能打到日志里面
		////zap.String("手机号码", req.Phone)
		//
		//// 最多这样打日志，要非常小心
		//zap.L().Debug("", zap.String("手机号码", req.Phone))
		//return Result{Code: 5, Msg: "系统异常"}, err
		return Result{Code: 5, Msg: "系统异常"}, fmt.Errorf("用户手机号码登录失败 %w", err)
	}
	if !ok {
		//ctx.JSON(http.StatusOK, Result{
		//	Code: 4,
		//	Msg:  "验证码有误",
		//})
		return Result{Code: 4, Msg: "验证码有误"}, nil
	}

	// 我这个手机号，会不会是一个新用户
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		//ctx.JSON(http.StatusOK, Result{
		//	Code: 5,
		//	Msg:  "系统错误",
		//})
		return Result{Code: 5, Msg: "系统错误"}, fmt.Errorf("登录或者注册用户失败 %w", err)
	}
	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		// 记录日志
		//ctx.JSON(http.StatusOK, Result{
		//	Code: 5,
		//	Msg:  "系统错误",
		//})
		return Result{Code: 5, Msg: "系统错误"}, fmt.Errorf("登录或者注册用户失败 %w", err)
	}

	//ctx.JSON(http.StatusOK, Result{
	//	Msg: "验证码校验通过",
	//})
	return Result{Msg: "登陆成功"}, nil
}

func (u *UserHandler) SendLoginSMACode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 是不是一个合法的手机号码
	// 考虑正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{Msg: "发送成功"})
	case service.ErrCodeSendTooMany:
		zap.L().Warn("短信发送太频繁",
			zap.Error(err))
		ctx.JSON(http.StatusOK, Result{Msg: "发送次数太频繁，请稍后再试"})
	default:
		zap.L().Warn("发送短信失败",
			zap.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
	}
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	// 只有超时才会有 err
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 按道理这里应该记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、英文和特殊字符")
		return
	}

	// 调用一下 svc
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常 ")
		return
	}
	ctx.String(http.StatusOK, "注册成功 ")
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserHandler) LoginJWT(ctx *gin.Context, req LoginReq) (Result, error) {
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		//ctx.String(http.StatusOK, "用户名或密码不对")
		return Result{Code: 4, Msg: "用户名或密码不对"}, nil
	}
	if err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		return Result{Code: 5, Msg: "系统错误"}, nil
	}

	// 在这里登陆成功了
	// 我可以随便设置值了
	// 步骤2
	// 在这里用 JWT 设置登录态
	// 生成一个 token
	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		return Result{Msg: "系统错误"}, nil
	}

	//ctx.String(http.StatusOK, "登陆成功")
	return Result{Msg: "登陆成功"}, nil
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 在这里登陆成功了
	// 设置 session
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 要放在 session 里面的值
	// 步骤2
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		// 生产环境再设置，开发环境就不用了
		//Secure: true,
		//HttpOnly: true,
		MaxAge: 60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登陆成功")
	return
}

func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登陆失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登陆OK",
	})
}

//func (u *UserHandler) Do(fn func(ctx context.Context)(any, error))  {
//	data, err := fn(ctx)
//	if err != nil {
//		// 在这里打日志
//	}
//}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 要放在 session 里面的值
	// 步骤2
	sess.Options(sessions.Options{
		// 生产环境再设置，开发环境就不用了
		//Secure: true,
		//HttpOnly: true,
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "成功退出登陆")
}

func (c *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		// 注意，其他字段，尤其是密码，邮箱和手机
		// 修改都要通过别的手段
		// 修改邮箱和手机都要验证
		// 密码就更加不用说了
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
	}
	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "昵称不能为空",
		})
		return
	}
	if len(req.AboutMe) > 1024 {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "关于我过长",
		})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		// 也就是说，我们其实并没有直接校验具体的格式
		// 而是如果你能转化过来，那就说明没问题
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "日期格式不对",
		})
		return
	}
	uc := ctx.MustGet("user").(*ijwt.UserClaims)
	err = c.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		AboutMe:  req.AboutMe,
		Birthday: birthday,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})
}

func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		AboutMe  string
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	u, err := c.svc.Profile(ctx, uc.Uid)
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Birthday: u.Birthday.Format(time.DateOnly),
		AboutMe:  u.AboutMe,
	})
	//c, _ := ctx.Get("claims")
	// 你可以断定， 必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	//claims, ok := c.(*UserClaims)
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	//println(claims.Uid)
	//ctx.String(http.StatusOK, "你的 profile")
	//// 这边就是你补充 profile 的其他代码
}

func (c *UserHandler) Profile(ctx *gin.Context) {
	//type Profile struct {
	//	Email string
	//}
	//sess := sessions.Default(ctx)
	//id := sess.Get(userIdKey).(int64)
	//u, err := c.svc.Profile(ctx, id)
	//if err != nil {
	//	// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
	//	// 那就说明是系统出了问题
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	//ctx.JSON(http.StatusOK, Profile{
	//	Email: u.Email,
	//})
}
