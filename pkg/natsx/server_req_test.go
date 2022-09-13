package natsx

import (
	"context"
	"fmt"
	"testing"

	"github.com/junaozun/game_server/pkg/natsx/testdata"
)

type TestReqServer struct {
	total int32
}

func (a *TestReqServer) AddTestMine(ctx context.Context, req *testdata.TestMine) (*testdata.TestMineResp, error) {
	a.total++
	fmt.Println(a.total)
	repl := &testdata.TestMineResp{
		Id:      req.Id + req.Id,
		Brother: "sister",
		Childs:  []int64{1, 2, 3, 4, 5, 6},
	}
	return repl, nil
}

type TestReqServer2 struct {
	total int32
}

func (a *TestReqServer2) AddTestMine(ctx context.Context, req *testdata.TestMine) (*testdata.TestMineResp, error) {
	a.total++
	fmt.Println(a.total)
	repl := &testdata.TestMineResp{
		Id:      req.Id + req.Id,
		Brother: "brother",
		Childs:  []int64{1, 2, 3, 4, 5, 6},
	}
	return repl, nil
}

func TestStartServer(t *testing.T) {
	connEnc, err := NewNatsJSONEnc("nats://0.0.0.0:4222")
	if err != nil {
		t.Error(err)
		return
	}
	server, err := NewServer(connEnc)
	if err != nil {
		t.Error(err)
		return
	}
	err = server.Register("chess", &TestReqServer{})
	err = server.Register("chess", &TestReqServer2{})
	for {

	}
}

func TestStartClient(t *testing.T) {
	connEnc, err := NewNatsJSONEnc("nats://0.0.0.0:4222")
	if err != nil {
		t.Error(err)
		return
	}
	chessClient, err := NewClient(connEnc, "chess")
	if err != nil {
		t.Error(err)
		return
	}
	req := &testdata.TestMine{
		Id:   512,
		Name: "houwenwen",
		Sex:  "man",
	}
	resp := &testdata.TestMineResp{}
	err = chessClient.Request(context.Background(), "TestReqServer2", "AddTestMine", req, resp)
	if err != nil {
		t.Error(err)
		return
	}
}
