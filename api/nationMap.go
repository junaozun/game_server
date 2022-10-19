package api

type ConfigReq struct {
}

type ConfigRsp struct {
	Confs []Conf
}

type Conf struct {
	Type     int8   `json:"type"`
	Level    int8   `json:"level"`
	Name     string `json:"name"`
	Wood     int    `json:"Wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Grain    int    `json:"grain"`
	Durable  int    `json:"durable"`  // 耐久
	Defender int    `json:"defender"` // 防御等级
}

//
// type ScanBlockReq struct {
// 	X      int `json:"x"`
// 	Y      int `json:"y"`
// 	Length int `json:"length"`
// }
//
// type ScanRsp struct {
// 	MRBuilds []MapRoleBuild `json:"mr_builds"` // 角色建筑，包含被占领的基础建筑
// 	MCBuilds []MapRoleCity  `json:"mc_builds"` // 角色城市
// 	Armys    []Army         `json:"armys"`     // 军队
// }
//
// type BuildReq struct {
// 	X    int  `json:"x"`
// 	Y    int  `json:"y"`
// 	Type int8 `json:"type"`
// }
//
// type BuildRsp struct {
// 	X    int  `json:"x"`
// 	Y    int  `json:"y"`
// 	Type int8 `json:"type"`
// }
//
// type GiveUpReq struct {
// 	X int `json:"x"`
// 	Y int `json:"y"`
// }
//
// type GiveUpRsp struct {
// 	X int `json:"x"`
// 	Y int `json:"y"`
// }
//
// type UpBuildReq struct {
// 	X int `json:"x"`
// 	Y int `json:"y"`
// }
//
// type UpBuildRsp struct {
// 	X     int          `json:"x"`
// 	Y     int          `json:"y"`
// 	Build MapRoleBuild `json:"build"`
// }
