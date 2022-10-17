package ws

import "encoding/json"

//构建发送消息
func BuildSendToClientMsg(code int, msg, updateType string, data interface{}) []byte {
	if data == nil {
		data = map[string]interface{}{}
	}
	msgBody := ResponseMsgBody{
		Code:       code,
		UpdateType: updateType,
		Msg:        msg,
		Data:       data,
	}
	jsonByte, _ := json.Marshal(msgBody)
	return jsonByte
}
