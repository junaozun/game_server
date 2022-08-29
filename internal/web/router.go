package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaozun/game_server/internal/web/controller/account"
	"github.com/junaozun/game_server/pkg/httpx/middleware"
)

func RegisterRouters(webSrv *WebService) func(g *gin.Engine) {
	return func(g *gin.Engine) {
		g.Use(middleware.Cors())
		g.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})
		accountCtl := account.NewAccount(webSrv.Component.Dao)
		account := g.Group("/account")
		{
			account.Any("/register", accountCtl.Register)
		}
	}
}
