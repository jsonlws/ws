package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

/*
这里为websockt核心代码
websocket理解为http的升级版本即可
将所有用户抽象成对象User，User中应当包括一个连接和一个消息信道
数据处理器Hub:用于获取到某个用户发送的数据推送给每个用户
*/

type User struct {
	conn     *websocket.Conn
	msg      chan []byte
	Uid      string
	GroupId  string
	UserLock *sync.Mutex
}

type Hub struct {
	//实现组模式的消息发送
	groupList    map[string][]*User
	wsUserToConn map[string]*User
	//注册chan，用户注册时添加到chan中
	register chan *User
	//注销chan，用户退出时添加到chan中，再从map中删除
	unregister chan *User
	//广播消息，将消息广播给所有连接
	broadcast chan []byte
}

//实例化数据处理器
func NewHub() *Hub {
	return &Hub{
		groupList:    make(map[string][]*User),
		wsUserToConn: make(map[string]*User),
		register:     make(chan *User),
		unregister:   make(chan *User),
		broadcast:    make(chan []byte),
	}
}

//处理中心处理获取到的信息
func (h *Hub) Run() {
	for {
		select {
		//从注册chan中取数据
		case user := <-h.register:
			//绑定组关系
			h.groupList[user.GroupId] = append(h.groupList[user.GroupId], user)
			h.wsUserToConn[user.Uid] = user
		//处理用户下线
		case user := <-h.unregister:
			//删除链接信息
			delete(h.wsUserToConn, user.Uid)
			user.conn.Close()
		case data := <-h.broadcast:
			//从广播chan中取消息，然后遍历给每个用户，发送到用户的msg中
			for _, u := range h.wsUserToConn {
				go u.conn.WriteMessage(1, data)
			}
		}
	}
}
