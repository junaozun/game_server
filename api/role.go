package api

// EnterServerReq 进入游戏
type EnterServerReq struct {
	Session string `json:"session"`
}

type EnterServerRsp struct {
	Role    RoleBase `json:"role"`
	RoleRes RoleRes  `json:"role_res"`
	Time    int64    `json:"time"`
	Token   string   `json:"token"`
}

// 数据库的字段 不一定是客户端需要的字段，做业务逻辑的时候 会将数据库的结果 映射到客户端需要的结果上
// 其中 可能会做一些转换
// dto data trasfer object business object
type RoleBase struct {
	RId      int    `json:"rid"`
	UId      int    `json:"uid"`
	NickName string `json:"nickName"`
	Sex      int8   `json:"sex"`
	Balance  int    `json:"balance"`
	HeadId   int16  `json:"headId"`
	Profile  string `json:"profile"`
}

type RoleRes struct {
	Wood          int `json:"wood"`
	Iron          int `json:"iron"`
	Stone         int `json:"stone"`
	Grain         int `json:"grain"`
	Gold          int `json:"gold"`
	Decree        int `json:"decree"` // 令牌
	WoodYield     int `json:"wood_yield"`
	IronYield     int `json:"iron_yield"`
	StoneYield    int `json:"stone_yield"`
	GrainYield    int `json:"grain_yield"`
	GoldYield     int `json:"gold_yield"`
	DepotCapacity int `json:"depot_capacity"` // 仓库容量
}
