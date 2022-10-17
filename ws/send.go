package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//发送心态回复
func sendPingMsg(hub *Hub, connInfo *websocket.Conn, oldTime uint, msgId string) {
	rspData := fmt.Sprintf(`{"updateType":"ping","old_time":%d,"messageId":%s}`, oldTime, msgId)
	hub.singleMsg <- &SingleChanDef{
		MsgByte: []byte(rspData),
		WsConn:  connInfo,
	}
}

//发送登录成功或失败回复
func sendLoginNoticeMsg(conn *websocket.Conn, isSuccess bool, oldTime uint) {
	rspData := `{"updateType":"AuthenticationFailure"}`
	if isSuccess {
		rspData = fmt.Sprintf(`{"updateType":"AuthenticationSuccess","old_time":%d}`, oldTime)
	}
	conn.WriteMessage(1, []byte(rspData))
}

//发送消息给某一个用户
func sendSingleMsg(hub *Hub, connInfo *websocket.Conn, respData []byte) {
	hub.singleMsg <- &SingleChanDef{
		MsgByte: respData,
		WsConn:  connInfo,
	}
}

//发送群组消息
func sendGroupMsg(hub *Hub, groupId string, respData []byte) {
	hub.groupMsg <- &GroupChanDef{
		MsgByte: respData,
		GroupId: groupId,
	}
}
