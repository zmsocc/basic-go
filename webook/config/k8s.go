//go:build k8s
// 使用k8s这个编译标签
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-live-mysql:11309)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	}
}
