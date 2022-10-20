package game

import (
	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/internal/logic/data"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/game_server/ret"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/mitchellh/mapstructure"
)

type MyProperty struct {
	game *Game
}

func NewMyProperty(game *Game) *MyProperty {
	return &MyProperty{
		game: game,
	}
}

func (m *MyProperty) RegisterRouter(cb func(command ExecCommand)) {
	cb(ExecCommand{
		group:    "role",
		name:     "myProperty",
		execFunc: m.myProperty,
	})
}

// 获取角色拥有的属性
// 资源、池、建筑、部队、武将
func (m *MyProperty) myProperty(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {

	rsp.Body.Router = req.Body.Router
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Code = ret.OK.Code

	reqParams := &api.MyRolePropertyReq{}
	respParams := &api.MyRolePropertyRsp{}

	err := mapstructure.Decode(req.Body.Msg, reqParams)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[MyProperty] myProperty 解析请求参数错误")
		rsp.Body.Code = ret.Err_Param.Code
		return
	}

	// 获取角色
	role, ok := req.Conn.GetProperty("role")
	if !ok { // 未登录
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[MyProperty] myProperty session不合法")
		rsp.Body.Code = ret.Err_SessionInvalid.Code
		return
	}
	// 根据角色id查询角色的资源、城池、建筑、部队、武将
	rid := role.(*data.Role).Id

	// 1.查资源
	roleRes, code := m.game.Role.GetRoleRes(rid)
	if code != ret.OK.Code {
		rsp.Body.Code = code
		return
	}
	respParams.RoleRes = roleRes

	// 2.查城池
	roleCitys, code := m.game.Role.GetRoleCitys(rid)
	if code != ret.OK.Code {
		rsp.Body.Code = code
		return
	}
	respParams.Citys = roleCitys

	// 3.查建筑
	roleBuilds, code := m.game.Role.GetBuilds(rid)
	if code != ret.OK.Code {
		rsp.Body.Code = code
		return
	}
	respParams.MRBuilds = roleBuilds

	// 4.查部队
	roleArmys, code := m.game.Role.GetArmys(rid)
	if code != ret.OK.Code {
		rsp.Body.Code = code
		return
	}
	respParams.Armys = roleArmys

	// 5.查武将
	roleCommanders, code := m.game.Role.GetCommanders(rid)
	if code != ret.OK.Code {
		rsp.Body.Code = code
		return
	}
	respParams.Commander = roleCommanders

	rsp.Body.Msg = respParams
}
