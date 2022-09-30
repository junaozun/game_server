package nats_handler

import (
	"context"
	"fmt"

	"github.com/junaozun/gogopkg/natsx"
	"github.com/junaozun/gogopkg/natsx/testdata"
)

type TestGvg struct {
	total int32
}

func (a *TestGvg) TestGvgFunc(ctx context.Context, req *testdata.TestMine) (*testdata.TestMineResp, error) {
	a.total++
	fmt.Println(a.total)
	repl := &testdata.TestMineResp{
		Id:      req.Id + req.Id,
		Brother: "gvg",
		Childs:  []int64{1, 2, 3, 4, 5, 6},
	}
	return repl, nil
}

func RegisterHandler(natsxSrv *natsx.NatsxServer) {
	natsxSrv.Register(natsxSrv.ServerName, &TestGvg{})
}
