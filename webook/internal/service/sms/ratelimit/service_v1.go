package ratelimit

import (
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms/service"
	"gitee.com/zmsoc/gogogo/webook/pkg/ratelimit"
)

type RatelimitSMSServiceV1 struct {
	// 这样的话，注释 Send 方法不会报错
	service.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSServiceV1(svc service.Service, limiter ratelimit.Limiter) service.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

//func (s *RatelimitSMSServiceV1) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
//	limited, err := s.limiter.Limit(ctx, "sms:tencent")
//	if err != nil {
//		// 系统错误
//		// 可以限流：保守策略，你的下游很坑的时候，
//		// 可以不限: 你的下游很强，业务可用性要求很高，尽量容错策略
//		// 包一下这个错误
//		return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)
//	}
//	if limited {
//		return errLimited
//	}
//	// 你在这里加一些代码，新特性
//	err = s.Service.Send(ctx, tpl, args, numbers...)
//	// 你在这里也可以加一些代码，新特性
//	return err
//}
