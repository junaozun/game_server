package web

import (
	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/dao"
	"github.com/junaozun/game_server/pkg/httpx"
)

func NewHttpServiceWithConfig(webCfg *config.ServerConfig, commonCfg *config.CommonConfig) *httpx.HttpxServer {
	service := NewWebService(webCfg, commonCfg)
	routers := RegisterRouters(service)
	httpServer, err := httpx.New(routers, httpx.WithAddress("0.0.0.0:"+webCfg.Port))
	if err != nil {
		panic(err)
	}
	return httpServer
}

type WebService struct {
	Component *component.Component
}

func NewWebService(webCfg *config.ServerConfig, commonCfg *config.CommonConfig) *WebService {
	// 初始化数据访问层
	dao, err := dao.NewDao([]interface{}{webCfg.Mysql, commonCfg.Etcd, commonCfg.Cache})
	if err != nil {
		panic(err)
	}
	return &WebService{
		Component: component.NewComponent(dao),
	}
}
