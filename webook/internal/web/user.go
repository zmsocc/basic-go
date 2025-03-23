package web

import (
	"fmt"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const biz = "login"

// 确保 UserHandler 上实现了 handler 接口
var _ handler = &UserHandler{}

// 这个更优雅
var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         *service.UserService
	codeSvc     *service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
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
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.POST("/login_sms/code/send", u.SendLoginSMACode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}

	// 我这个手机号，会不会是一个新用户
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这边要怎么办？
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
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
		ctx.JSON(http.StatusOK, Result{Msg: "发送次数太频繁，请稍后再试"})
	default:
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

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
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
	// 我可以随便设置值了
	// 步骤2
	// 在这里用 JWT 设置登录态
	// 生成一个 token
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(user)
	ctx.String(http.StatusOK, "登陆成功")
	return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("hC2pcTKJUakr7wXNmu2xd4WHxKAJpFDE"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
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

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	// 你可以断定， 必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	claims, ok := c.(*UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
	ctx.String(http.StatusOK, "你的 profile")
	// 这边就是你补充 profile 的其他代码
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的 Profile")
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	Uid int64
	// 自己随便加
	UserAgent string
}
