package http

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func HandleHttpToWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("{\"error\":10001,\"msg\":\"请求方式错误\"}"))
		return
	}
	var temp [1024]byte
	dataLen, _ := r.Body.Read(temp[:])
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:"+viper.GetString("service.port")+"/?token=0&flag=inside", nil)
	if err != nil {
		w.Write([]byte("{\"error\":10002,\"msg\":\"无法连接socket服务\"}"))
		return
	}
	conn.WriteMessage(1, temp[:dataLen])
	w.Write([]byte("{\"error\":0,\"msg\":\"发送消息成功\"}"))
	conn.Close()
}
