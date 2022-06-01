package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//定义一个升级器，将普通的http连接升级为websocket连接
var wsUpgrader = &websocket.Upgrader{
	HandshakeTimeout:  5 * time.Second,
	EnableCompression: true,
	//定义读写缓冲区大小
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	//校验请求
	CheckOrigin: func(r *http.Request) bool {
		//如果不是get请求，返回错误
		if r.Method != "GET" {
			log.Println("请求方式错误")
			return false
		}

		if len(r.URL.Query()["uid"]) == 0 {
			log.Println("请求PATH中缺少uid参数")
			return false
		}

		//还可以根据其他需求定制校验规则
		return true
	},
}

//处理ws连接
func WsHandle(w http.ResponseWriter, r *http.Request, hub *Hub, heartbeat *Bucket) {
	//通过升级后的升级器得到链接
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade http conn err:", err)
		return
	}

	var uid, groupId string
	//通信标识
	txFlag := r.URL.Query()["flag"]
	//走内部rpc通信使用
	if len(txFlag) == 1 && txFlag[0] == "inside" {
		uid, groupId = handleInsideCorrespond(r)
	} else {
		uid, groupId = handleExternalCorrespond(conn, r, hub, heartbeat)
		if uid == "" {
			return
		}
	}

	//实例化用户信息
	user := &User{
		conn:     conn,
		msg:      make(chan []byte, 256),
		Uid:      uid,
		GroupId:  groupId,
		UserLock: new(sync.Mutex),
	}

	//进行用户信息放入内存当中
	hub.register <- user

	uidint, _ := strconv.Atoi(uid)
	oldTime := heartbeat.FirstHeartHandler(uint(uidint), user)

	//发送鉴权成功
	sendLoginNoticeMsg(conn, true, oldTime)

	//得到连接后，就可以开始读写数据了
	go user.readPump(hub, heartbeat)

}

//处理外部部通信
func handleExternalCorrespond(conn *websocket.Conn, r *http.Request, hub *Hub, heartbeat *Bucket) (uid, groupId string) {
	uids := r.URL.Query()["uid"][0]
	//账号重复登录得问题
	if oldConn, ok := hub.wsUserToConn[uids]; ok {
		rspData := fmt.Sprintf(`{"action_type":"distance_login","err_msg":%s}`, "账号异地登录")
		oldConn.conn.WriteMessage(1, []byte(rspData))
		oldConn.conn.Close() //关闭连接
	}
	return uids, "0"
}

//处理内部通信
func handleInsideCorrespond(r *http.Request) (uid, groupId string) {
	return "0", "0"
}

//读取发送消息
func (c *User) readPump(hub *Hub, heartbeat *Bucket) {
	//创建两个协程一个协程读取数据，一个协程进行写回数据给客户端
	dataChannel := make(chan []byte, 1024)

	//读协程
	go func() {
		for {
			msgType, msg, err := c.conn.ReadMessage()
			if err != nil || msgType == -1 {
				log.Println("用户退出:", c.Uid, "原因:", err.Error())
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
				sendErrorMsg(c, "数据异常")
				hub.unregister <- c
				return
			}

			switch requestData.ActionType {
			//发送广播消息
			case All:
				//将读取到的信息传入websocket处理器中的broadcast中，
				hub.broadcast <- readData
				//发送给某个用户消息
			case Single:
				uid := strconv.Itoa(int(requestData.Receiver))
				if conn, ok := hub.wsUserToConn[uid]; ok {
					sendSingleMsg(conn, uid, readData)
				}
				//发送群组消息
			case Group:
				sendGroupMsg(hub, strconv.Itoa(int(requestData.Receiver)), readData)
				//回复心态消息
			case Ping:

				oldTime, err := heartbeat.FutureHeartHandler(requestData.Sender, requestData.OldIndex, c)
				if err != nil {
					sendErrorMsg(c, err.Error())
				} else {
					sendPingMsg(c, oldTime)
				}
			default:
				sendErrorMsg(c, "不支持的消息类型")
			}

		}
	}()
}
