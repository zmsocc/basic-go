package main

import "github.com/gin-gonic/gin"

func main() {
	server := gin.Default()

	server.POST("/users/signup", func(context *gin.Context) {

	})
	//// 这是 REST 风格
	//server.PUT("/user", func(context *gin.Context) {
	//
	//})

	server.POST("/users/login", func(context *gin.Context) {

	})

	server.POST("/users/edit", func(context *gin.Context) {

	})
	//// 这是 REST 风格
	//server.POST("/users/:id", func(context *gin.Context) {
	//
	//})

	server.GET("/users/profile", func(context *gin.Context) {

	})
	// REST 风格
	server.GET("/users/:id", func(context *gin.Context) {

	})
	server.Run(":8080")
}
