package game

import (
	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/global"
	"github.com/junaozun/game_server/internal/logic/game_config"
	"github.com/junaozun/game_server/internal/logic/model"
	"github.com/junaozun/game_server/pkg/utils"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/game_server/ret"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/mitchellh/mapstructure"
)

type Role struct {
	game *Game
}

func NewRole(game *Game) *Role {
	return &Role{
		game: game,
	}
}

func (a *Role) RegisterRouter(cb func(command ExecCommand)) {
	cb(ExecCommand{
		group:    "role",
		name:     "enterServer",
		execFunc: a.enterServer,
	})
}

// 进入游戏的逻辑
func (a *Role) enterServer(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {
	// session 需要验证是否合法，合法的情况下可以取出登录的用户ID
	// 根据用户id,去查询对应的游戏角色，如果有，就继续，没有提示无角色
	// 根据角色id，查询角色拥有的资源roleRes，如果资源有，返回。没有，初始化资源

	rsp.Body.Router = req.Body.Router
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Code = ret.OK.Code

	reqParam := &api.EnterServerReq{}
	respParam := &api.EnterServerRsp{}

	err := mapstructure.Decode(req.Body.Msg, reqParam)
	if err != nil {
		rsp.Body.Code = ret.Err_Param.Code
		return
	}
	session := reqParam.Session

	claim, err := utils.CheckToken(session)
	if err != nil {
		rsp.Body.Code = ret.Err_SessionInvalid.Code
		return
	}

	db := a.game.Dao.DB
	uid := claim.Uid

	role := &model.Role{}
	err = db.Where("uid=?", uid).Find(role).Error
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[role] enterServer 查询角色失败")
		rsp.Body.Code = ret.Err_DB.Code
		return
	}

	// 角色不存在
	if role.Id == 0 {
		// todo 创建一个
		err := db.Create(&model.Role{
			UId:        int(uid),
			NickName:   "牛油果树",
			Balance:    123,
			HeadId:     331,
			Sex:        1,
			Profile:    "ooooooo",
			LoginTime:  global.Now(),
			LogoutTime: global.Now(),
			CreatedAt:  global.Now(),
		}).Error
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{"err": err}).Error("[role] enterServer 创建角色失败")
			rsp.Body.Code = ret.Err_RoleNotExist.Code
			return
		}
	}
	rid := role.Id
	roleRes := &model.RoleRes{}
	err = db.Where("rid=?", rid).Find(roleRes).Error
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[role] enterServer查询角色资源失败")
		rsp.Body.Code = ret.Err_DB.Code
		return
	}

	// 角色还没有资源,则初始化资源
	if roleRes.Id == 0 {
		roleRes.RId = rid
		roleRes.Wood = game_config.Base.Role.Wood
		roleRes.Iron = game_config.Base.Role.Iron
		roleRes.Stone = game_config.Base.Role.Stone
		roleRes.Grain = game_config.Base.Role.Grain
		roleRes.Gold = game_config.Base.Role.Gold
		roleRes.Decree = game_config.Base.Role.Decree
		err = db.Create(roleRes).Error
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{}).Error("[role] enterServer 初始化资源插入数据库失败")
			rsp.Body.Code = ret.Err_DB.Code
			return
		}
	}

	respParam.RoleRes = roleRes.ToClient()
	respParam.Role = role.ToClient()
	respParam.Time = global.Now().Unix()
	// 生成一个角色id的token
	token, err := utils.SetToken(uint64(rid))
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[role] enterServer 生成角色token失败")
		rsp.Body.Code = ret.Err_TokenGenERR.Code
		return
	}
	respParam.Token = token

	rsp.Body.Msg = respParam
}
