package server

import (
	"github.com/junaozun/game_server/net"
	"github.com/junaozun/game_server/server/controller"
)

func InitRouter(router *net.Router) {
	controller.DefaultAccount.Router(router)
}
