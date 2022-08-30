package web

import (
	"context"
	"log"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/pkg/app"
	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/dao"
	"github.com/junaozun/game_server/pkg/httpx"
)

type WebApp struct {
	httpxServer *httpx.HttpxServer
}

func NewWebApp(cfg pkgConfig.GameConfig) *WebApp {
	dao, err := dao.NewDao([]interface{}{cfg.Web.Mysql, cfg.Common.Etcd, cfg.Common.Cache})
	if err != nil {
		panic(err)
	}
	routers := RegisterRouters(component.NewComponent(dao))
	httpServer, err := httpx.New(routers, httpx.WithAddress("0.0.0.0:"+cfg.Web.Port))
	if err != nil {
		panic(err)
	}
	return &WebApp{
		httpxServer: httpServer,
	}
}

func (w *WebApp) Run(ctx context.Context, cfg pkgConfig.GameConfig) error {
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
