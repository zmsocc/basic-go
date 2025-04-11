package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms/memory"
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms/service"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) service.Service {
	// 换内存还是换别的
	//svc := ratelimit.NewRatelimitSMSService(memory.NewService(),
	//	limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	//return retryable.NewService(svc, 3)
	return memory.NewService()
}
