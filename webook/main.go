package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"net/http"
)

func main() {
	//db := initDB()
	//rdb := initRedis()
	//server := initWebServer()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)

	initViperV1()
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好 你来了")
	})
	server.Run(":8080")
}

func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:11316)/webook"

redis:
  addr: "localhost:6379"
`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

// 看起来，remote 不支持 key 的切割
func initViperRemote() {
	viper.SetConfigType("yaml")
	// 通过 webook 和其他使用 etcd 的区别出来
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379",
		"/webook")
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml",
		"指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	//viper.SetDefault("db.mysql.dsn",
	//	"root:root@tcp(localhost:3306)/webook")
	//viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViper() {
	// 配置文件的名字，但是不包含扩展名
	// 不包含 .go, .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置文件的格式
	// 现实中，有很多格式，JSON，XML，YAML...
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./config")
	// 读取配置到 viper 里面，或者你可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//otherViper := viper.New()
	//otherViper.SetConfigName("myjson")
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigType("json")
}
