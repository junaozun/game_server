package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaozun/game_server/internal/web/controller/account"
	"github.com/junaozun/gogopkg/httpx/middleware"
)

type WebRouter func(engine *gin.Engine)

func NewWebRouter(accountCtl *account.AccountCtl) WebRouter {
	gin.SetMode(gin.ReleaseMode)
	return func(g *gin.Engine) {
		g.Use(middleware.Cors())
		g.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})

		account := g.Group("/account")
		{
			account.Any("/register", accountCtl.Register)
			account.Any("/nats_rpc", accountCtl.UseNatsTest)
		}
	}
}
