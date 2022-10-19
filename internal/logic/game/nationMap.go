package game

import (
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/game_server/ret"
)

type NationMap struct {
	game *Game
}

func NewNationMap(game *Game) *NationMap {
	return &NationMap{
		game: game,
	}
}

func (n *NationMap) RegisterRouter(cb func(command ExecCommand)) {
	cb(ExecCommand{
		group:    "nationMap",
		name:     "config",
		execFunc: n.nationMapConfig,
	})
}

//
func (n *NationMap) nationMapConfig(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {

	rsp.Body.Router = req.Body.Router
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Code = ret.OK.Code

}
