package errno

type Err struct {
	Code int         `json:"code"` // 业务编码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 成功时返回的数据
}

func NewError(code int, msg string) Err {
	return Err{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func (e Err) WithData(data interface{}) Err {
	e.Data = data
	return e
}
