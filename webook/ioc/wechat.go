package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service/oauth2/wechat"
	"gitee.com/zmsoc/gogogo/webook/internal/web"
	"net/http"
	"os"
)

func InitWechatService() wechat.Service {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID env variable is not set")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET env variable is not set")
	}
	return wechat.NewService(appId, appKey, http.DefaultClient)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
