package ws

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"yunyuim/lib"

	"github.com/gorilla/websocket"
)

/**
websocket连接层处理，这里不能有太多逻辑处理不然并发会急剧下降
*/

var Lock sync.Mutex

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

		if len(r.URL.Query()["token"]) == 0 {
			log.Println("请求PATH中缺少token参数")
			return false
		}

		//还可以根据其他需求定制校验规则
		return true
	},
}

//处理ws连接
func WsConnectHandler(w http.ResponseWriter, r *http.Request, hub *Hub, heartbeat *Bucket) {
	//通过升级后的升级器得到链接
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade http conn err:", err)
		return
	}

	var (
		uid     string
		groupId string
	)
	//通信标识
	txFlag := r.URL.Query()["flag"]
	//走内部rpc通信使用
	if len(txFlag) == 1 && txFlag[0] == "inside" {
		uid, groupId = handleInsideCorrespond(r)
	} else {
		uid, groupId = handleExternalCorrespond(conn, r, hub, heartbeat)
		if uid == "" {
			sendLoginNoticeMsg(conn, false, 0)
			conn.Close()
			return
		}
	}

	//实例化用户信息
	user := &User{
		conn:    conn,
		Uid:     uid,
		GroupId: groupId,
	}

	//进行用户信息放入内存当中
	hub.register <- user

	//系统不用注册到心跳处理中
	//if uid != "0" {
	oldTime := heartbeat.FirstHeartHandler(uid, user, &Lock)
	//发送鉴权成功
	sendLoginNoticeMsg(conn, true, oldTime)
	//}

	//得到连接后，就可以开始读写数据了
	go user.readPump(hub, heartbeat)

	//发送广播消息
	//sendBroadcastMsg(user)

}

//处理外部部通信
func handleExternalCorrespond(conn *websocket.Conn, r *http.Request, hub *Hub, heartbeat *Bucket) (uid, groupId string) {
	tokenStr := r.URL.Query()["token"][0]
	userInfo, err := lib.ParseToken(tokenStr)
	if err != nil {
		log.Println("解析token失败", err)
		return "", ""
	}

	return strconv.Itoa(userInfo.UserToken), strconv.Itoa(userInfo.ShopToken)
}

//处理内部通信
func handleInsideCorrespond(r *http.Request) (uid, groupId string) {
	return "0", "0"
}
