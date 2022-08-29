package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaozun/game_server/internal/web/service"
	"github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/httpx"
	"github.com/junaozun/game_server/pkg/httpx/middleware"
)

type WebService struct {
	*httpx.HttpxServer
}

func NewWebServiceWithConfig(cfg *config.WebConfig) *WebService {
	service := service.NewService()
	routers := RegisterRouters(service)
	httpServer, err := httpx.New(routers, httpx.WithAddress("0.0.0.0:"+cfg.Port))
	if err != nil {
		panic(err)
	}
	return &WebService{
		HttpxServer: httpServer,
	}
}

func RegisterRouters(s *service.Service) func(g *gin.Engine) {
	return func(g *gin.Engine) {
		g.Use(middleware.Cors())
		g.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})
	}
}
