package ws

/*
消息结构体相关定义
*/

//通讯类型定义
const (
	Single = "single"
	Group  = "group"
	All    = "all"
	Ping   = "ping"
)

//服务器接收数据结构体定义
type ReqMsgBody struct {
	Sender     uint        `json:"sender"`      //发送者
	ActionType string      `json:"action_type"` //操作类型
	OldIndex   uint        `json:"old_index"`   //服务器返回，客户端只需要原值返回
	MsgBody    interface{} `json:"msg_body"`    //消息体
	Receiver   uint        `json:"receiver"`    //接收者
}
