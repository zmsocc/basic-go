package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//db := initDB()
	//rdb := initRedis()
	//server := initWebServer()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)

	server := initWebServer()
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

	server.Use()

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
	server.Use()
	return server
}
