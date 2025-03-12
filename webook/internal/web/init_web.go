package web

import "github.com/gin-gonic/gin"

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	u := &UserHandler{}
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/login", u.Login)
	server.POST("/users/edit", u.Edit)
	server.GET("/users/profile", u.Profile)
	return server
}

/*func registerUserRoutes(server *gin.Engine) {


/* //分组注册
ug := server.Group("/users")
ug.GET("/profile", u.Profile)
ug.POST("/signup", u.SignUp)
ug.POST("/login", u.Login)
ug.POST("/edit", u.Edit)*/

//// 这是 REST 风格
//server.PUT("/user", func(context *gin.Context) {
//
//})

//// 这是 REST 风格
//server.POST("/users/:id", func(context *gin.Context) {
//
//})

// REST 风格
//server.GET("/users/:id", func(context *gin.Context) {
//
//})
