package client_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/junaozun/game_server/net"
	"github.com/junaozun/game_server/utils"
	"github.com/mitchellh/mapstructure"
)

var (
	Ws   *websocket.Conn
	Skey string
)

func TestMain(m *testing.M) {
	var err error
	url := "ws://localhost:8003" // 服务器地址
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
		handshake := &net.Handshake{}
		body := &net.RespBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Fatal(err)
		}
		if body.Router == net.HandshakeMsg {
			err := mapstructure.Decode(body.Msg, handshake)
			if err != nil {
				log.Println("Decode err")
				return
			}
			Skey = handshake.Key
			break
		}
	}
	os.Exit(m.Run())
}

func SendWsData(req *net.ReqBody) *net.RespBody {
	data, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	// 对数据加密
	encryptData, err := utils.AesCBCEncrypt(data, []byte(Skey), []byte(Skey), openssl.ZEROS_PADDING)
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
	realData, err := utils.AesCBCDecrypt(resData, []byte(Skey), []byte(Skey), openssl.ZEROS_PADDING)
	if err != nil {
		log.Println("数据格式有误，解密失败：", err)
		return nil
	}
	respBody := &net.RespBody{}
	err = json.Unmarshal(realData, respBody)
	if err != nil {
		log.Println("数据解析失败", err)
		return nil
	}
	return respBody
}
