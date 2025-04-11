//go:build !k8s

// 没有k8s这个编译标签
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:11326)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
