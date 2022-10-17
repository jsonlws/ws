package ws

/*
消息结构体相关定义
*/

//通讯类型定义
const (
	Single    = "single"
	Group     = "group"
	All       = "all"
	Ping      = "ping"
	PayNotice = "payNotice"
)

//定义消息结构体
type MsgBody struct {
	RequestType string                 `json:"request_type"` //通讯类型
	Content     string                 `json:"content"`      //消息内容
	MessageId   string                 `json:"message_id"`   //消息id
	Accepter    string                 `json:"accepter"`     //消息接收者 当request-type 为 single时就是uid  group时就是groupid 广播时此参数忽略
	Data        map[string]interface{} `json:"data"`         //数据
}

//updateType常量定义
const (
	AuthSuccess   = "AuthenticationSuccess" //鉴权成功
	AuthFailure   = "AuthenticationFailure" //鉴权失败
	PingResponse  = "ping"                  //心跳回复
	RemoteLogin   = "re_login"              //异地登录
	ErrorResponse = "error"                 //错误标识
)

//回复消息结构体定义
type ResponseMsgBody struct {
	Code       int         `json:"code"`
	Data       interface{} `json:"data"`
	Msg        string      `json:"msg"`
	UpdateType string      `json:"updateType"`
}

//服务器接收数据结构体定义
type ReqMsgBody struct {
	RequestType string      `json:"requestType"` //客户端定义操作类型
	OldTime     uint        `json:"old_time"`    //服务器返回，客户端只需要原值返回
	MsgBody     interface{} `json:"content"`     //消息体
	Receiver    string      `json:"receiver"`    //接收者
	MessageId   string      `json:"messageId"`   //消息id
}
