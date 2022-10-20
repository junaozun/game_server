package login

import (
	"log"

	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/global"
	"github.com/junaozun/game_server/internal/login/data"
	common_model "github.com/junaozun/game_server/model"
	"github.com/junaozun/game_server/pkg/utils"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/game_server/ret"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm/clause"
)

func (l *LoginApp) InitLogin() {
	l.initTable()
	l.initRouter()
}

func (l *LoginApp) initRouter() {
	l.Router.Group("account").AddRouter("login", l.login)
}

func (g *LoginApp) initTable() {
	err := g.Dao.DB.AutoMigrate(
		new(common_model.User),
		new(data.LoginHistory),
		new(data.LoginLast),
	)
	if err != nil {
		panic(err)
	}
}

func (l *LoginApp) login(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {
	loginReq := &api.LoginReq{}
	err := mapstructure.Decode(req.Body.Msg, loginReq)
	if err != nil {
		log.Println("[Account] login err")
		return
	}

	user := &common_model.User{}
	db := l.Dao.DB
	err = db.Where(&common_model.User{Username: loginReq.Username}).Find(user).Error
	if err != nil {
		return
	}

	if user.ID == 0 { // 用户不存在
		rsp.Body.Code = ret.Err_UserNotFound.Code
		return
	}

	// 检查密码
	if utils.ScryptPasswd(loginReq.Password) != user.Passwd {
		rsp.Body.Code = ret.Err_PasswdNotRight.Code
		return
	}

	token, err := utils.SetToken(user.ID)
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[Account] login token生成错误")
		rsp.Body.Code = ret.Err_TokenGenERR.Code
		return
	}

	rsp.Body.Code = ret.OK.Code
	loginResp := &api.LoginRsp{
		Username: user.Username,
		Session:  token,
		UId:      user.ID,
	}
	rsp.Body.Msg = loginResp

	// 保存用户登录记录
	loginRecord := &data.LoginHistory{
		UId:       user.ID,
		UserName:  user.Username,
		LoginTime: global.Now(),
		Ip:        loginReq.Ip,
		State:     data.UserStatus_Login,
		Hardware:  loginReq.Hardware,
	}
	err = db.Create(loginRecord).Error
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[Account] save loginRecord error")
		return
	}

	// 保存用户的最后一次登录
	// upsert 没有插入，有则更新
	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},                                                           //  这里的列必须是唯一的，比如主键或是唯一索引
		DoUpdates: clause.AssignmentColumns([]string{"loginTime", "ip", "session", "isLogout", "hardware"}), // 更新哪些字段
	}).Create(&data.LoginLast{
		UId:        user.ID,
		LoginTime:  global.Now(),
		LogoutTime: global.Now(),
		Ip:         loginReq.Ip,
		Session:    token,
		IsLogout:   data.UserStatus_Login,
		Hardware:   loginReq.Hardware,
	}).Error
	if err != nil {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[Account] save 保存用户的最后一次登录错误")
		return
	}
	// 缓存一下此用户和当前的ws连接
	l.onLineUser.UserLogin(req.Conn, user.ID, token)
}
