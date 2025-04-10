// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/article"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/cache"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/dao"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	"gitee.com/zmsoc/gogogo/webook/internal/web"
	"gitee.com/zmsoc/gogogo/webook/internal/web/jwt"
	"gitee.com/zmsoc/gogogo/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	loggerV1 := InitLog()
	handler := jwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitMiddlewares(cmdable, loggerV1, handler)
	gormDB := InitTestDB()
	userDAO := dao.NewUserDAO(gormDB)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository, loggerV1)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService(cmdable)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := InitPhantomWechatService(loggerV1)
	wechatHandlerConfig := InitWechatHandlerConfig()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, wechatHandlerConfig, handler)
	aticleDAO := dao.NewGORMArticleDAO(gormDB)
	articleRepository := article.NewArticleRepository(aticleDAO)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	return engine
}

func InitArticleHandler() *web.ArticleHandler {
	gormDB := InitTestDB()
	aticleDAO := dao.NewGORMArticleDAO(gormDB)
	articleRepository := article.NewArticleRepository(aticleDAO)
	articleService := service.NewArticleService(articleRepository)
	loggerV1 := InitLog()
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	return articleHandler
}

func InitUserSvc() service.UserService {
	gormDB := InitTestDB()
	userDAO := dao.NewUserDAO(gormDB)
	cmdable := InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	loggerV1 := InitLog()
	userService := service.NewUserService(userRepository, loggerV1)
	return userService
}

func InitJwtHdl() jwt.Handler {
	cmdable := InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	return handler
}

// wire.go:

var thirdProvider = wire.NewSet(
	InitRedis,
	InitTestDB,
	InitLog,
)

var userSvcProvider = wire.NewSet(dao.NewUserDAO, cache.NewUserCache, repository.NewUserRepository, service.NewUserService)
