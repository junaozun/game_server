package client_test

import (
	"fmt"
	"testing"

	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/mitchellh/mapstructure"
)

func TestLogin(t *testing.T) {
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

func TestRegister(t *testing.T) {

}
