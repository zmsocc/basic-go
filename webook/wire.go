//go:build wireinject

package main

import (
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/cache"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/dao"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	"gitee.com/zmsoc/gogogo/webook/internal/web"
	"gitee.com/zmsoc/gogogo/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		// 初始化 DAO
		dao.NewUserDAO,
		// 初始化 缓存
		cache.NewUserCache,
		cache.NewCodeCache,
		//
		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		// 直接基于内存实现
		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
