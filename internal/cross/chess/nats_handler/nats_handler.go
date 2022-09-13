package nats_handler

import (
	"context"
	"fmt"

	"github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/natsx/nats_register"
	"github.com/junaozun/game_server/pkg/natsx/testdata"
)

type TestChess struct {
	total int32
}

func (a *TestChess) TestChessFunc(ctx context.Context, req *testdata.TestMine) (*testdata.TestMineResp, error) {
	a.total++
	fmt.Println(a.total)
	repl := &testdata.TestMineResp{
		Id:      req.Id + req.Id,
		Brother: "chess",
		Childs:  []int64{1, 2, 3, 4, 5, 6},
	}
	return repl, nil
}

func RegisterHandler(serverName string, natsCfg *config.NatsConfig) {
	natsHandler := nats_register.RegisterNats(natsCfg)
	natsHandler(serverName, &TestChess{})
}
