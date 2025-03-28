package web

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms"
	"github.com/gin-gonic/gin"
)

type OAuth2WechatHandler struct {
	svc sms.Service
}

func NewOAuth2WechatHandler() {

}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {

}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {

}

//type OAuth2Handler struct {
//	wechatService string
//	svcs          map[string]OAuth2Service
//}
//
//func (h *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
//	// 统一处理所有的 OAuth2 的
//	g := server.Group("/oauth2")
//	g.GET("/:platform/authurl", h.AuthURL)
//	g.GET("/:platform/callback", h.Callback)
//}
//
//func (h *OAuth2Handler) AuthURL(ctx *gin.Context) {
//	platform := ctx.Param("platform")
//	switch platform {
//	case "wechat":
//		h.wechatService.AuthURL
//	}
//
//	svc := h.svcs
//}
//
//func (h *OAuth2Handler) Callback(ctx *gin.Context) {
//
//}
