package api

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hardware string `json:"hardware"`
}
