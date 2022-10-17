package ws

import (
	"encoding/json"
	"fmt"
	"log"
)

/**
websocket消息读取写回
*/

//读取发送消息
func (c *User) readPump(hub *Hub, heartbeat *Bucket) {
	//创建两个协程一个协程读取数据，一个协程进行写回数据给客户端
	dataChannel := make(chan []byte, 1024)

	//读协程
	go func() {
		for {
			msgType, msg, err := c.conn.ReadMessage()
			if err != nil || msgType == -1 {
				log.Println("Type", msgType, "用户退出:", c.Uid, "原因:", err.Error())
				hub.unregister <- c
				break
			}

			dataChannel <- msg
		}
	}()

	//写协程
	go func() {
		for {
			readData := <-dataChannel
			var requestData ReqMsgBody
			err := json.Unmarshal(readData, &requestData)
			//说明客户端乱发送数据直接给断开
			if err != nil {
				errMsgStr := `{"code":1002,"updateType":"err","msg":"数据异常"}`
				sendSingleMsg(hub, c.conn, []byte(errMsgStr))
				return
			}

			switch requestData.RequestType {
			//发送广播消息
			case All:
				//将读取到的信息传入websocket处理器中的broadcast中，
				hub.broadcast <- readData
				//发送给某个用户消息
			case Single:
				uid := requestData.Receiver
				if conn, ok := hub.wsUserToConn[uid]; ok {
					sendSingleMsg(hub, conn.conn, readData)
				} else {
					log.Printf("用户%s不在线\n", uid)
				}
				//发送群组消息
			case Group:
				sendGroupMsg(hub, requestData.Receiver, readData)
				//回复心态消息
			case Ping:
				oldTime, err := heartbeat.FutureHeartHandler(c.Uid, requestData.OldTime, c, &Lock)
				if err != nil {
					errMsgStr := fmt.Sprintf(`{"code":1002,"updateType":"err","msg":%s}`, err.Error())

					sendSingleMsg(hub, c.conn, []byte(errMsgStr))

				} else {
					sendPingMsg(hub, c.conn, oldTime, requestData.MessageId)
				}
			case PayNotice:
			default:
				errMsgStr := `{"code":1002,"updateType":"err","msg":"不支持的消息类型"}`
				sendSingleMsg(hub, c.conn, []byte(errMsgStr))
			}

		}
	}()
}
