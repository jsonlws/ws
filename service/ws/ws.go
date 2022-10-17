package main

import (
	"flag"
	myHttp "jsonlwsim/http"
	"jsonlwsim/ws"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var configFile = flag.String("f", "./../../config/config.json", "提供配置文件")

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
	log.Printf("Service startup succeeded,running on the %s \n", servicePort)
	//启动webscoket服务
	hub := ws.NewHub()
	go hub.Run()

	//启动心跳检测器
	checkHeart := ws.NewHeartBeat()
	go checkHeart.TurnClockwise(hub)

	//websocket处理类
	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.WsConnectHandler(w, r, hub, checkHeart)
	})

	//处理http消息转发到ws
	go http.HandleFunc("/httpToWs", myHttp.HandleHttpToWs)

	err := http.ListenAndServe(":"+servicePort, nil)
	if err != nil {
		log.Println("Service startup failure,reason:", err)
	}
}
