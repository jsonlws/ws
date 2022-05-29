package main

import (
	"flag"
	"log"
	"net/http"
	myHttp "yunyuim/http"
	"yunyuim/ws"

	"github.com/spf13/viper"
)

var configFile = flag.String("f", "", "提供配置文件")

//读取配置文件
func initConfingFile() {
	flag.Parse()

	if (*configFile) == "" {
		panic("请提供服务启动配置文件")
	}
	//读取配置文件
	viper.SetConfigFile(*configFile)

	// 处理读取配置文件的错误
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	//实现动态读取
	viper.WatchConfig()
}

func init() {
	initConfingFile()
}

func main() {
	servicePort := viper.GetString("service.port")
	log.Println("服务启动成功,运行端口为[", servicePort, "]")
	hub := ws.NewHub()
	go hub.Run()
	//启动心跳检测器
	checkHeart := ws.NewHeartBeat()
	go checkHeart.TurnClockwise(hub)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.WsHandle(w, r, hub, checkHeart)
	})
	//处理http消息转发到ws
	http.HandleFunc("/httpToWs", myHttp.HandleHttpToWs)
	http.ListenAndServe(":"+servicePort, nil)

}
