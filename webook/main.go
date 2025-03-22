package main

import (
	"gitee.com/zmsoc/gogogo/webook/config"
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/cache"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/dao"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms/memory"
	"gitee.com/zmsoc/gogogo/webook/internal/web"
	"gitee.com/zmsoc/gogogo/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()

	rdb := initRedis()
	u := initUser(db, rdb)
	u.RegisterRoutes(server)

	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好 你来了")
	})
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(func(ctx *gin.Context) {
		println("这是第一个 middleware")
	})

	server.Use(func(ctx *gin.Context) {
		println("这是第二个 middleware")
	})

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		// ExposeHeaders 不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token"},
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
	}))

	// 步骤1
	//store := cookie.NewStore([]byte("secret"))
	// 单实例部署
	//store := memstore.NewMemStore([]byte("hC2pcTKJUakr7wXNmu2xd4WHxKAJpFDE"),
	//	[]byte("EtreByecTpnpSA5WkwD3Mz5sQQbCnz6R"))

	// 多实例部署
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("hC2pcTKJUakr7wXNmu2xd4WHxKAJpFDE"),
	//	[]byte("EtreByecTpnpSA5WkwD3Mz5sQQbCnz6R"))
	//if err != nil {
	//	panic(err)
	//}

	//server.Use(sessions.Sessions("mysession", store))
	// 步骤3
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		IgnorePaths("/users/login").Build())
	return server
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(webook-live-mysql:11316)/webook"))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return rdb
}
