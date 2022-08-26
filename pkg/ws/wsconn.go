package ws

type ReqBody struct {
	Seq    int64       `json:"seq"`
	Router string      `json:"router"` // account.login
	Msg    interface{} `json:"msg"`
	Proxy  string      `json:"proxy"` // 多进程，服务之间调用
}

type RespBody struct {
	Seq    int64       `json:"seq"`
	Router string      `json:"router"`
	Code   int         `json:"code"`
	Msg    interface{} `json:"msg"`
}

type WsMsgReq struct {
	Body *ReqBody
	Conn IWsConn
}

type WsMsgResp struct {
	Body *RespBody
}

type IWsConn interface {
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
