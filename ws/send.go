package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//发送错误消息提示
func sendErrorMsg(user *User, msg string) {
	user.UserLock.Lock()
	defer user.UserLock.Unlock()

	rspData := fmt.Sprintf(`{"action_type":"login_success","err_msg":%s}`, msg)
	user.conn.WriteMessage(1, []byte(rspData))
}

//发送心态回复
func sendPingMsg(user *User, oldTime uint) {
	user.UserLock.Lock()
	defer user.UserLock.Unlock()

	rspData := fmt.Sprintf(`{"action_type":"ping","old_index":%d}`, oldTime)

	user.conn.WriteMessage(1, []byte(rspData))
}

//发送登录成功或失败回复
func sendLoginNoticeMsg(conn *websocket.Conn, isSuccess bool, oldIndex uint) {

	rspData := fmt.Sprintf(`{"action_type":"login_success","old_index":%d}`, oldIndex)
	conn.WriteMessage(1, []byte(rspData))
}

//发送消息给某一个用户
func sendSingleMsg(userConn *User, receiver string, respData []byte) {
	userConn.UserLock.Lock()
	defer userConn.UserLock.Unlock()
	userConn.conn.WriteMessage(1, respData)
}

//发送消息一个组中的所有连接对象
func sendGroupMsg(hub *Hub, groupId string, respData []byte) {
	if connList, ok := hub.groupList[groupId]; ok {
		if len(connList) == 0 {
			return
		}
		for _, connInfo := range connList {
			go sendSingleMsg(connInfo, connInfo.Uid, respData)
		}
	}
}
