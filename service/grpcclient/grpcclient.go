package main

import (
	"context"
	"fmt"

	po "yunyuim/grpc/sendmsg"
	"yunyuim/ws"

	"google.golang.org/grpc"
)

func main() {
	// 连接服务器
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := po.NewSendMsgServiceClient(conn)
	// 调用服务端的SayHello
	r, err := c.SendMsg(context.Background(), &po.SendMsgRequest{
		ActionType: ws.Single,
		Sender:     0,
		Data:       map[string]string{"content": "msg"},
		Receiver:   1,
	})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	fmt.Println(r)
}
