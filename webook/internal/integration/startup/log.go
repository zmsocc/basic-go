package startup

import (
	"gitee.com/zmsoc/gogogo/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
