//go:build k8s

// 使用k8s这个编译标签
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-live-mysql:11316)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-live-redis:11479",
	},
}
