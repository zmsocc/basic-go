package ioc

import (
	"context"
	"gitee.com/zmsoc/gogogo/webook/internal/web"
	ijwt "gitee.com/zmsoc/gogogo/webook/internal/web/jwt"
	"gitee.com/zmsoc/gogogo/webook/internal/web/middleware"
	"gitee.com/zmsoc/gogogo/webook/pkg/ginx/middlewares/logger"
	logger2 "gitee.com/zmsoc/gogogo/webook/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable,
	l logger2.LoggerV1, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
		l.Debug("HTTP请求", logger2.Field{
			Key:   "al",
			Value: al,
		})
	}).AllowReqBody(true).AllowRespBody()
	viper.OnConfigChange(func(in fsnotify.Event) {
		ok := viper.GetBool("web.logreq")
		bd.AllowReqBody(ok)
	})
	return []gin.HandlerFunc{
		corsHdl(),
		bd.Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/users/login").
			Build(),
		//ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		// ExposeHeaders 不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// AllowCredentials 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
