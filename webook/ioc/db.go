package ioc

import (
	"gitee.com/zmsoc/gogogo/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:11316)/webook_default",
	}
	err := viper.UnmarshalKey("db", &cfg)
	//if err != nil {
	//	panic(err)
	//}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
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
