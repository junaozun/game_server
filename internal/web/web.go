package web

import (
	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/web/wire"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/dao"
	"github.com/junaozun/gogopkg/httpx"
	"github.com/junaozun/gogopkg/logrusx"
)

type WebApp struct {
	httpxServer *httpx.HttpxServer
}

func NewWebApp(cfg config.GameConfig) *WebApp {
	dao, err := dao.NewDao([]interface{}{cfg.Web.Mysql, cfg.Common.Etcd, cfg.Common.Redis})
	if err != nil {
		panic(err)
	}
	component := component.NewComponent(dao, cfg)
	routers := wire.NewWebRouterMgr(component)
	httpServer, err := httpx.New(routers, httpx.WithAddress("0.0.0.0:"+cfg.Web.Port), httpx.WithPProf(false))
	if err != nil {
		panic(err)
	}
	return &WebApp{
		httpxServer: httpServer,
	}
}

func (w *WebApp) Run() error {
	web := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": w.httpxServer.Addr,
			}).Info("web app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.WithFields(logrusx.Fields{
				"addr": w.httpxServer.Addr,
			}).Info("web app exit .....")
		}),
		app.Name("web"),
		app.Runners(w.httpxServer),
	)
	if err := web.Run(); err != nil {
		return err
	}
	return nil
}
