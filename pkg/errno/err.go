package errno

type err struct {
	Code int         `json:"code"` // 业务编码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 成功时返回的数据
}

func NewError(code int, msg string) err {
	return err{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func (e err) WithData(data interface{}) err {
	e.Data = data
	return e
}
