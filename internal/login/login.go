package login

import (
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/internal/login/wsMgr"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/dao"
	"github.com/junaozun/gogopkg/logrusx"
)

const (
	host     = "0.0.0.0:"
	ServerId = "100001"
)

type LoginApp struct {
	onLineUser *wsMgr.WsMgr
	Router     *ws.Router
	Dao        *dao.Dao
	ServerName string
}

func NewLoginApp() *LoginApp {
	return &LoginApp{
		onLineUser: wsMgr.NewWsMgr(),
		Router:     ws.NewRouter(),
		ServerName: common.ServerName_Login,
	}
}

func (l *LoginApp) Run(cfg config.GameConfig) error {
	dao, err := dao.NewDao([]interface{}{cfg.Common.Mysql})
	if err != nil {
		panic(err)
	}
	l.Dao = dao
	l.InitLogin()
	wsServer := ws.NewWsServer(host+cfg.Logic.Port, l.Router, false)
	login := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": wsServer.Addr,
			}).Info("login app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": wsServer.Addr,
			}).Info("login app exit .....")
		}),
		app.Name(l.ServerName+ServerId),
		app.Runners(wsServer),
	)
	if err := login.Run(); err != nil {
		return err
	}
	return nil
}
