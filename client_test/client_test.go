package client_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/junaozun/game_server/common"
	"github.com/junaozun/game_server/pkg/utils"
	"github.com/junaozun/game_server/pkg/ws"
	"github.com/mitchellh/mapstructure"
)

var (
	Ws        *websocket.Conn // 客户端与gateway的长链接
	Secretkey string
)

func TestMain(m *testing.M) {
	var err error
	url := "ws://localhost:8004" // gateway地址
	Ws, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	for {
		_, data, err := Ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式,需要json数据：", err)
			return
		}
		handshake := &ws.Handshake{}
		body := &ws.RespBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Fatal(err)
		}
		if body.Router == common.HandshakeMsg {
			err := mapstructure.Decode(body.Msg, handshake)
			if err != nil {
				log.Println("Decode err")
				return
			}
			Secretkey = handshake.Key
			break
		}
	}
	os.Exit(m.Run())
}

func SendWsData(req *ws.ReqBody) *ws.RespBody {
	data, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	// 对数据加密
	encryptData, err := utils.AesCBCEncrypt(data, []byte(Secretkey), []byte(Secretkey), openssl.ZEROS_PADDING)
	if err != nil {
		log.Println(err)
		return nil
	}
	// 再对数据进行压缩
	zipData, err := utils.Zip(encryptData)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = Ws.WriteMessage(websocket.BinaryMessage, zipData)
	if err != nil {
		log.Println(err)
		return nil
	}
	_, res, err := Ws.ReadMessage()
	if err != nil {
		log.Println(err)
		return nil
	}
	// 对返回数据解压
	resData, err := utils.UnZip(res)
	if err != nil {
		log.Println("解压数据出错，非法格式,需要json数据：", err)
		return nil
	}
	realData, err := utils.AesCBCDecrypt(resData, []byte(Secretkey), []byte(Secretkey), openssl.ZEROS_PADDING)
	if err != nil {
		log.Println("数据格式有误，解密失败：", err)
		return nil
	}
	respBody := &ws.RespBody{}
	err = json.Unmarshal(realData, respBody)
	if err != nil {
		log.Println("数据解析失败", err)
		return nil
	}
	return respBody
}
