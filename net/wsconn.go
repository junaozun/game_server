package net

type ReqBody struct {
	Seq   int64       `json:"seq"`
	Name  string      `json:"name"`
	Msg   interface{} `json:"msg"`
	Proxy string      `json:"proxy"` // 多进程，服务之间调用
}

type RespBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

type WsMsgReq struct {
	Body *ReqBody
	Conn WsConn
}

type WsMsgResp struct {
	Body *RespBody
}

type WsConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, bool)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

const HandshakeMsg = "handshake"

type Handshake struct {
	Key string `json:"key"`
}
