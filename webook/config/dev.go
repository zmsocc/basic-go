//go:build !k8s

// 没有k8s这个编译标签
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-live-mysql:11316)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
