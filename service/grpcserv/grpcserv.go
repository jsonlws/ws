package main

import (
	context "context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	po "yunyuim/grpc/sendmsg"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
}

var WsConn *websocket.Conn

var WsError error

var WsRsp *http.Response

var MsgChan chan []byte

var PingChan chan uint

func (this *server) SendMsg(ctx context.Context, req *po.SendMsgRequest) (*po.SendMsgResponse, error) {
	sendByte, err := json.Marshal(req)

	if err != nil {
		return &po.SendMsgResponse{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}
	log.Println("接收消息", string(sendByte))

	MsgChan <- sendByte

	return &po.SendMsgResponse{
		Code: 200,
		Msg:  "发送消息成功",
	}, nil
}

//读取配置文件
var configFile = flag.String("f", "./../../config/config.json", "提供配置文件")

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
	log.Println("rpc服务启动成功，运行端口为", viper.GetString("rpc.port"))

	go connectWs()

	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":"+viper.GetString("rpc.port"))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer() // 创建gRPC服务器
	po.RegisterSendMsgServiceServer(s, &server{})

	reflection.Register(s) //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}

func SendMsgToWs() {
	for {
		select {
		case data := <-MsgChan:
			WsConn.WriteMessage(1, data)
		}
	}
}

func SendPingToWs() {
	for {
		select {
		case data := <-PingChan:
			time.Sleep(35 * time.Second) //心跳间隔发送时间
			pingContent := fmt.Sprintf(`{"content": "ping","old_time": %d,"requestType": "ping","messageId": "1"}`, data)
			log.Println("进行心跳回复处理", pingContent)
			WsConn.WriteMessage(1, []byte(pingContent))
		}
	}
}

type MsgDef struct {
	UpdateType string `json:"updateType"`
	OldTime    uint   `json:"old_time"`
}

func connectWs() {
	MsgChan = make(chan []byte, 1024)
	PingChan = make(chan uint, 1)
	//已经实现重连逻辑
	reconnectionHandler()
	//读取消息协程
	go func() {
		for {
			defer func() {
				if err := recover(); err != nil {
					log.Println("ws服务端关闭,err:", err)
					reconnectionHandler()
				}
			}()
			_, msg, _ := WsConn.ReadMessage()
			var msgbodyStruct MsgDef
			json.Unmarshal(msg, &msgbodyStruct)
			if msgbodyStruct.UpdateType == "AuthenticationSuccess" || msgbodyStruct.UpdateType == "ping" {
				log.Println("收到心跳回复包数据", string(msg))
				PingChan <- msgbodyStruct.OldTime
			}
		}
	}()
	//发送消息协程
	go SendMsgToWs()
	//心跳专用协程
	go SendPingToWs()
}

func reconnectionHandler() {
	WsConn, WsRsp, WsError = websocket.DefaultDialer.Dial("ws://127.0.0.1:"+viper.GetString("service.port")+"?token=0&flag=inside", nil)
	if WsError != nil {
		log.Println("进入重连逻辑,原因socket连接失败err:", WsError)
		//若连接失败进入每隔两秒进行重连
		time.Sleep(2 * time.Second)
		reconnectionHandler()
		return
	}
}
