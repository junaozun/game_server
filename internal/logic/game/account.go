package game

import (
	"log"

	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/global"
	"github.com/junaozun/game_server/internal/logic/model"
	"github.com/junaozun/game_server/net"
	"github.com/junaozun/game_server/ret"
	"github.com/junaozun/game_server/utils"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm/clause"
)

type Account struct {
	game *Game
}

func NewAccount(game *Game) *Account {
	return &Account{
		game: game,
	}
}

func (a *Account) RegisterRouter(cb func(command ExecCommand)) {
	cb(ExecCommand{
		group:    "account",
		name:     "login",
		execFunc: a.login,
	})
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	loginReq := &api.LoginReq{}
	err := mapstructure.Decode(req.Body.Msg, loginReq)
	if err != nil {
		log.Println("[Account] login err")
		return
	}

	user := &model.User{}
	db := a.game.Dao.DB
	err = db.Where(&model.User{Username: loginReq.Username}).Find(user).Error
	if err != nil {
		return
	}

	if user.UId == 0 { // 用户不存在
		rsp.Body.Code = ret.Err_UserNotFound
		return
	}

	// 检查密码
	if utils.ScryptPasswd(loginReq.Password) != user.Passwd {
		rsp.Body.Code = ret.Err_PasswdNotRight
		return
	}

	token, err := utils.SetToken(user.UId)
	if err != nil {
		return
	}

	rsp.Body.Code = ret.Ok
	loginResp := &api.LoginRsp{
		Username: user.Username,
		Session:  token,
		UId:      user.UId,
	}
	rsp.Body.Msg = loginResp

	// 保存用户登录记录
	loginRecord := &model.LoginHistory{
		UId:       user.UId,
		UserName:  user.Username,
		LoginTime: global.Now(),
		Ip:        loginReq.Ip,
		State:     model.UserStatus_Login,
		Hardware:  loginReq.Hardware,
	}
	err = db.Create(loginRecord).Error
	if err != nil {
		log.Println("[Account] save loginRecord error")
		return
	}

	// 保存用户的最后一次登录
	// upsert 没有插入，有则更新
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},                                                             //  这里的列必须是唯一的，比如主键或是唯一索引
		DoUpdates: clause.AssignmentColumns([]string{"login_time", "ip", "session", "is_logout", "hardware"}), // 更新哪些字段
	}).Create(&model.LoginLast{
		UId:       user.UId,
		LoginTime: global.Now(),
		Ip:        loginReq.Ip,
		Session:   token,
		IsLogout:  model.UserStatus_Login,
		Hardware:  loginReq.Hardware,
	})
	// 缓存一下此用户和当前的ws连接
	a.game.UserLogin(req.Conn, user.UId, token)
}
