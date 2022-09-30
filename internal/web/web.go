package web

import (
	"log"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/web/wire"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/dao"
	"github.com/junaozun/gogopkg/httpx"
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
	httpServer, err := httpx.New(routers, httpx.WithAddress("0.0.0.0:"+cfg.Web.Port))
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
			log.Printf("web app start addr:%s ....", w.httpxServer.Addr)
		}),
		app.OnExitHook(func() {
			log.Printf("web app exit addr:%s ....", w.httpxServer.Addr)
		}),
		app.Name("web"),
		app.Runners(w.httpxServer),
	)
	if err := web.Run(); err != nil {
		return err
	}
	return nil
}
