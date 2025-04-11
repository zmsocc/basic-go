package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/repository/dao"
	"gitee.com/zmsoc/gogogo/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:11326)/webook_default",
	}
	err := viper.UnmarshalKey("db", &cfg)
	//if err != nil {
	//	panic(err)
	//}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		// 缺了一个 writer
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值，只有执行时间超过这个阈值，才会使用
			// 50ms， 100ms
			// SQL 查询必然会要求命中索引，最好就是走一次磁盘 IO
			// 一次磁盘 IO 是不到 10ms
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	//dao.NewUserDAOV1(func() *gorm.DB {
	//	viper.OnConfigChange(func(in fsnotify.Event) {
	//		db, err = gorm.Open(mysql.Open())
	//		pt := unsafe.Pointer(&db)
	//		atomic.StorePointer(&pt, unsafe.Pointer(&db))
	//	})
	//	// 要用原子操作
	//	return db
	//})

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{
		Key:   "args",
		Value: args,
	})
}

type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
