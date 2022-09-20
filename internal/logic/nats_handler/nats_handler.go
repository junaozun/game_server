package nats_handler

import (
	"context"
	"fmt"

	"github.com/junaozun/game_server/pkg/natsx"
	"github.com/junaozun/game_server/pkg/natsx/testdata"
)

type TestReqServer struct {
	total int64
}

func (a *TestReqServer) AddTestMine(ctx context.Context, req *testdata.TestMine) (*testdata.TestMineResp, error) {
	a.total += req.Id
	fmt.Println(a.total)
	repl := &testdata.TestMineResp{
		Id:      a.total,
		Brother: req.Name,
		Childs:  []int64{1, 2, 3, 4, 5, 6},
	}
	return repl, nil
}

func RegisterHandler(natsxSrv *natsx.NatsxServer, serverId string) {
	natsxSrv.Register(natsxSrv.ServerName, &TestReqServer{}, natsx.WithServiceID(serverId))
}
