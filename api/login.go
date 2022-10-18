package api

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Hardware string `json:"hardware"`
}

type LoginRsp struct {
	Username string `json:"username"`
	// Password string `json:"password"`
	Session string `json:"session"`
	UId     uint64 `json:"uid"`
}

type LogoutReq struct {
	UId int `json:"uid"`
}

type LogoutRsp struct {
	UId int `json:"uid"`
}

type ReLoginReq struct {
	Session  string `json:"session"`
	Ip       string `json:"ip"`
	Hardware string `json:"hardware"`
}

type ReLoginRsp struct {
	Session string `json:"session"`
}
