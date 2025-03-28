package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms"
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms/memory"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	return memory.NewService()
}
