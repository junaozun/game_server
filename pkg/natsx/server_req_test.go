package natsx

import (
	"context"
	"fmt"
	"testing"

	"github.com/junaozun/game_server/pkg/natsx/testdata"
)

type BenchNotifyService struct {
}

func (a *BenchNotifyService) Notify(ctx context.Context, req *testdata.TestMine) {
	fmt.Println("BenchNotifyService Notify success")
}

func (a *BenchNotifyService) Test(ctx context.Context, req *testdata.TestMine) {
	fmt.Println("BenchNotifyService Test success")
}

type BeginTime struct {
}

func (a *BeginTime) Calculate(ctx context.Context, req *testdata.TestMine) {
	fmt.Println("BeginTime Calculate success")
}

func (a *BeginTime) BeginEnd(ctx context.Context, req *testdata.TestMine) {
	fmt.Println("BeginTime BeginEnd success")
}

type Logic struct {
}

func (l *Logic) Suxuefeng(ctx context.Context, req *testdata.TestMine) {
	fmt.Println("Logic Suxuefeng success")
}

func TestServer(t *testing.T) {
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
	err = server.Register("chess", &BenchNotifyService{})
	err = server.Register("chess", &BeginTime{})
	err = server.Register("logic", &Logic{})
	for {

	}
}

func TestClient(t *testing.T) {
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
		Id:   7128,
		Name: "game_server",
		Sex:  "man",
	}
	err = chessClient.Publish("BeginTime", "Calculate", req)
	if err != nil {
		t.Error(err)
		return
	}
	chessClient.Publish("")
}
