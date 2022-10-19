package game

import (
	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/internal/logic/game_config"
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

// 加载地图建筑
func (n *NationMap) nationMapConfig(req *ws.WsMsgReq, rsp *ws.WsMsgResp) {

	rsp.Body.Router = req.Body.Router
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Code = ret.OK.Code

	// reqParams := &api.ConfigReq{}
	respParams := &api.ConfigRsp{}

	cfgs := game_config.MapBuildConf.Cfg

	for _, v := range cfgs {
		respParams.Confs = append(respParams.Confs, api.Conf{
			Type:     v.Type,
			Level:    v.Level,
			Name:     v.Name,
			Wood:     v.Wood,
			Iron:     v.Iron,
			Stone:    v.Stone,
			Grain:    v.Grain,
			Durable:  v.Durable,
			Defender: v.Defender,
		})
	}
	rsp.Body.Msg = respParams
}
