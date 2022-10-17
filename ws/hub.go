package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

/*
这里为websockt核心代码
websocket理解为http的升级版本即可
将所有用户抽象成对象User，User中应当包括一个连接和一个消息信道
数据处理器Hub:用于获取到某个用户发送的数据推送给每个用户
*/

//用户连接信息定义
type User struct {
	conn    *websocket.Conn //用户在ws中唯一连接对象
	Uid     string          //用户唯一标识
	GroupId string          //群组模式
}

//整个ws中心对象的定义
type Hub struct {
	groupList map[string][]*User //实现组模式的消息发送

	wsUserToConn map[string]*User //单个用户的连接信息

	register chan *User //注册chan，用户注册时添加到chan中

	unregister chan *User //注销chan，用户退出时添加到chan中，再从map中删除

	singleMsg chan *SingleChanDef //单聊消息发送通道

	groupMsg chan *GroupChanDef //组发送消息发送通道

	broadcast chan []byte //广播消息，将消息广播给所有连接
}

//单聊发送消息通道
type SingleChanDef struct {
	MsgByte []byte
	WsConn  *websocket.Conn
}

//群组发送消息通道
type GroupChanDef struct {
	MsgByte []byte
	GroupId string
}

//实例化数据处理器
func NewHub() *Hub {
	return &Hub{
		groupList:    make(map[string][]*User),
		wsUserToConn: make(map[string]*User),
		register:     make(chan *User, 1024),
		unregister:   make(chan *User, 1024),
		broadcast:    make(chan []byte, 1024),
		singleMsg:    make(chan *SingleChanDef, 1024),
		groupMsg:     make(chan *GroupChanDef, 1024),
	}
}

//处理中心处理获取到的信息
func (h *Hub) Run() {
	go h.registerUser()

	go h.eliminateUser()

	go h.sendSingleMsg()

	go h.sendRadioMsg()

	go h.sendGroupMsg()
}

//注册用户
func (h *Hub) registerUser() {
	for {
		select {
		//从注册chan中取数据
		case user := <-h.register:
			//绑定组关系
			h.groupList[user.GroupId] = append(h.groupList[user.GroupId], user)
			h.wsUserToConn[user.Uid] = user
		}
	}
}

//踢用户(即下线操作)
func (h *Hub) eliminateUser() {
	for {
		select {
		//从注册chan中取数据
		case user := <-h.unregister:
			log.Println("用户", user.Uid, "下线")
			//删除链接信息
			user.conn.Close()
			delete(h.wsUserToConn, user.Uid)
		}
	}
}

//发送给单用户消息
func (h *Hub) sendSingleMsg() {
	for {
		select {
		//从注册chan中取数据
		case data := <-h.singleMsg:
			data.WsConn.WriteMessage(1, data.MsgByte)
		}
	}
}

//发送群组消息
func (h *Hub) sendGroupMsg() {
	for {
		select {
		//从注册chan中取数据
		case data := <-h.groupMsg:
			//这里采用并发控制一次最多使用5个协程，防止协程泄漏
			ch := make(chan uint, 10)
			for _, v := range h.groupList[data.GroupId] {
				ch <- 0
				go func(conn *User) {
					conn.conn.WriteMessage(1, data.MsgByte)
					<-ch
				}(v)
			}
		}
	}
}

//发送广播消息
func (h *Hub) sendRadioMsg() {
	for {
		select {
		case data := <-h.broadcast:
			//这里采用并发控制一次最多使用5个协程，防止协程泄漏
			ch := make(chan uint, 10)
			//从广播chan中取消息，然后遍历给每个用户，发送到用户的msg中
			for _, v := range h.wsUserToConn {
				ch <- 0
				go func(conn *User) {
					conn.conn.WriteMessage(1, data)
					<-ch
				}(v)
			}
		}
	}
}
