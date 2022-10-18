package client_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/mitchellh/mapstructure"
)

func TestLogin(t *testing.T) {
	var count int
	for {
		time.Sleep(2 * time.Second)
		loginTest()
		count++
		if count > 10 {
			break
		}
	}
}

func loginTest() {
	loginReq := &api.LoginReq{
		Username: "suxuefeng",
		Password: "123456",
		Ip:       "127.0.0.1",
		Hardware: "mac",
	}
	var m map[string]interface{}
	mapstructure.Decode(loginReq, &m)
	req := &ws.ReqBody{
		Seq:    1,
		Router: "account.login",
		Msg:    m,
	}
	res := SendWsData(req)
	fmt.Println("receive: ", res)
}
