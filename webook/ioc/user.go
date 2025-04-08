package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"gitee.com/zmsoc/gogogo/webook/internal/service"
	"go.uber.org/zap"
)

func InitUserHander(repo repository.UserRepository) service.UserService {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return service.NewUserService(repo, l)
}
