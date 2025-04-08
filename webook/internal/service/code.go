package service

import (
	"context"
	"fmt"
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms"
	"go.uber.org/atomic"
	"math/rand"
)

var codeTplId atomic.String = atomic.String{}

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	codeTplId.Store("1877556")
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	codeTplId.Store(viper.GetString("code.tpl.id"))
	//})

	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// biz 区别业务场景
func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成一个验证码
	code := svc.generateCode()
	// 塞进去 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}
	// 发送出去
	err = svc.smsSvc.Send(ctx, codeTplId.Load(), []string{code}, phone)
	if err != nil {
		err = fmt.Errorf("发送短信出现异常 %w", err)
	}

	return err
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	// 六位数， num 在0， 1000000之间，包含0，不包含1000000
	num := rand.Intn(1000000)
	// 不够六位的，前面加上0，补够六位
	// 000001
	return fmt.Sprintf("%06d", num)
}
