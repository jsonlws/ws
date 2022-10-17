package main

import (
	"context"
	"fmt"

	po "yunyuim/grpc/sendmsg"
	"yunyuim/ws"

	"google.golang.org/grpc"
)

func main() {
	// 连接服务器47.107.32.231
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := po.NewSendMsgServiceClient(conn)
	// 调用服务端的SayHello
	r, err := c.SendMsg(context.Background(), &po.SendMsgRequest{
		RequestType: ws.Single,
		Receiver:    "323278010",
		Data:        map[string]string{"order_sn": "2022021015561395204583", "pay_status": "2"},
		Msg:         "支付提醒测试",
		UpdateType:  "PayNotice",
	})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	fmt.Println(r)
}
