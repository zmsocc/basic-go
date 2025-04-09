package ioc

import "gitee.com/zmsoc/gogogo/webook/pkg/logger"

//func InitLogger() logger.LoggerV1 {
//	l, err := logger.NewZapLogger()
//	if err != nil {
//		panic(err)
//	}
//	return logger.NewZapLogger(l)
//}

func InitLogger() logger.LoggerV1 {
	return &logger.NopLogger{}
}
