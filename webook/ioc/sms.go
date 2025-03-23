package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms"
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
