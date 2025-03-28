package async

//
//import (
//	"context"
//	"gitee.com/zmsoc/gogogo/webook/internal/service/sms"
//)
//
//type SMSService struct {
//	svc  sms.Service
//	repo repository.SMSAsyncReqRepository
//}
//
//func NewSMSService() *SMSService {
//	return &SMSService{}
//}
//
//func (s *SMSService) StartAsync() {
//	go func() {
//		req := s.repo.Find没发出去的请求()
//		for _, req := range reqs {
//			// 在这里发送， 并且控制重试
//			s.svc.Send(ctx, req.biz, req.args, req.numbers...)
//		}
//	}()
//}
//
//func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
//	//
//	err := s.svc.Send(ctx, biz, args, numbers...)
//	if err != nil {
//		// 判断是不是崩溃了
//		if "崩溃了" {
//			s.repo.Store()
//		}
//	}
//	return err
//}
