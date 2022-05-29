package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	po "yunyuim/grpc/sendmsg"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
}

func (this *server) SendMsg(ctx context.Context, req *po.SendMsgRequest) (*po.SendMsgResponse, error) {
	//等到websocket服务连接
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:"+viper.GetString("service.port")+"/?token=0&flag=inside", nil)
	if err != nil {
		return &po.SendMsgResponse{
			Code: 0,
			Msg:  "不能连接websocket服务",
		}, nil
	}
	sendByte, _ := json.Marshal(req)
	log.Println("接收消息", string(sendByte))
	conn.WriteMessage(1, sendByte)

	return &po.SendMsgResponse{
		Code: 0,
		Msg:  "发送消息成功",
	}, nil
}

//读取配置文件
func initConfingFile() {
	//读取配置文件
	viper.SetConfigFile("./../../config/config.json")

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
