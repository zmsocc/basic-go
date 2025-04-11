//go:build wireinject

package main

import (
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/article"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/cache"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/dao"
	article2 "gitee.com/zmsoc/gogogo/webook/internal/repository/dao/article"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	"gitee.com/zmsoc/gogogo/webook/internal/web"
	ijwt "gitee.com/zmsoc/gogogo/webook/internal/web/jwt"
	"gitee.com/zmsoc/gogogo/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		// 初始化 DAO
		dao.NewUserDAO,
		article2.NewGORMArticleDAO,
		// 初始化 缓存
		cache.NewUserCache,
		cache.NewCodeCache,
		//
		repository.NewUserRepository,
		repository.NewCodeRepository,
		article.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		// 直接基于内存实现
		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ioc.NewWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
