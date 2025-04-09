package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service/oauth2/wechat"
	"gitee.com/zmsoc/gogogo/webook/internal/web"
	logger2 "gitee.com/zmsoc/gogogo/webook/pkg/logger"
	"os"
)

func InitWechatService(l logger2.LoggerV1) wechat.Service {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID env variable is not set")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET env variable is not set")
	}
	return wechat.NewService(appId, appKey, l)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
