package natsx

import (
	"context"
	"fmt"
	"testing"

	"github.com/junaozun/game_server/pkg/natsx/testdata"
)

type ChessService1 struct {
}

func (a *ChessService1) Notify(ctx context.Context, req *testdata.TestMine) {
	fmt.Printf("ChessService1 Notify success req:%v", req)
}

func (a *ChessService1) Test(ctx context.Context, req *testdata.TestMine) {
	fmt.Printf("ChessService1 Test success req:%v", req)
}

type ChessService2 struct {
}

func (a *ChessService2) Calculate(ctx context.Context, req *testdata.TestMine) {
	fmt.Printf("ChessService2 Calculate success :%v", req)
}

func (a *ChessService2) BeginEnd(ctx context.Context, req *testdata.TestMine) {
	fmt.Printf("ChessService2 BeginEnd success :%v", req)
}

type LogicService1 struct {
}

func (l *LogicService1) Suxuefeng(ctx context.Context, req *testdata.TestMine) {
	fmt.Printf("LogicService1 Suxuefeng success req:%v", req)
}

func TestServer(t *testing.T) {
	connEnc, err := NewNatsPBEnc("nats://0.0.0.0:4222")
	if err != nil {
		t.Error(err)
		return
	}
	server, err := NewServer(connEnc)
	if err != nil {
		t.Error(err)
		return
	}
	server.Register("chess", &ChessService1{})
	server.Register("chess", &ChessService2{})
	server.Register("logic", &LogicService1{}, WithServiceID(100001))
	for {

	}
}

func TestClient(t *testing.T) {
	connEnc, err := NewNatsPBEnc("nats://0.0.0.0:4222")
	if err != nil {
		t.Error(err)
		return
	}
	/************************chess client *********************/
	chessClient, err := NewClient(connEnc, "chess")
	if err != nil {
		t.Error(err)
		return
	}
	req := &testdata.TestMine{
		Id:   7128,
		Name: "chess",
		Sex:  "man",
	}
	err = chessClient.Publish("ChessService1", "Notify", req)
	if err != nil {
		t.Error(err)
		return
	}
	/*****************************logic clinet ********************/
	logicClient, err := NewClient(connEnc, "logic")
	if err != nil {
		t.Error(err)
		return
	}
	req2 := &testdata.TestMine{
		Id:   123,
		Name: "logic",
		Sex:  "man",
	}
	err = logicClient.Publish("LogicService1", "Suxuefeng", req2, WithCallID(100001))
	if err != nil {
		t.Error(err)
		return
	}
}
