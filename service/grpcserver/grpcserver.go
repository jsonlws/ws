package main

import (
	context "context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
	po "yunyuim/grpc/sendmsg"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
}

//服务器接收数据结构体定义
type ReqMsgBody struct {
	ActionType string `json:"action_type"` //操作类型
	OldIndex   uint   `json:"old_index"`   //服务器返回，客户端只需要原值返回
}

//声明一个全局的消息传递协程通道
var GloabChan = make(chan []byte, 1024)

func (this *server) SendMsg(ctx context.Context, req *po.SendMsgRequest) (*po.SendMsgResponse, error) {
	// conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:"+viper.GetString("service.port")+"/?uid=0&flag=inside", nil)

	// defer conn.Close()
	sendByte, _ := json.Marshal(req)
	log.Println("接收消息", string(sendByte))

	//conn.WriteMessage(1, sendByte)

	GloabChan <- sendByte
	// go this.sendHeartHandler()

	return &po.SendMsgResponse{
		Code: 0,
		Msg:  "发送消息成功",
	}, nil
}

//websocket链接处理
func WsHandler() {
	conn, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:"+viper.GetString("service.port")+"/?uid=0&flag=inside", nil)
	reqData := make(chan []byte, 64)
	pingData := make(chan uint, 8)
	//读取消息协程
	go func() {
		for {
			select {
			case msg := <-GloabChan:
				conn.WriteMessage(1, msg)
			default:
				msgType, data, err := conn.ReadMessage()
				if err != nil || msgType == -1 {
					log.Println(err.Error())
				}
				reqData <- data
			}
		}
	}()
	//写处理消息协程
	go func() {
		for {
			var msgBody ReqMsgBody
			json.Unmarshal(<-reqData, &msgBody)
			log.Println("接收数据", msgBody)
			if msgBody.ActionType == "ping" || msgBody.ActionType == "login_success" {
				pingData <- msgBody.OldIndex
			}
		}
	}()

	//专门处理心跳的协程
	go func() {
		for {
			//每20秒执行一次
			time.Sleep(20 * time.Second)
			pingData := fmt.Sprintf(`{"sender":0,"action_type":"ping","old_index":%d}`, <-pingData)
			log.Println("发送心跳数据", pingData)
			conn.WriteMessage(1, []byte(pingData))

		}
	}()
}

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
	log.Println("rpc服务启动成功，运行端口为", viper.GetString("rpc.port"))
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":"+viper.GetString("rpc.port"))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	go WsHandler()
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
