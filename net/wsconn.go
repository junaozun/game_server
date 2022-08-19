package net

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/forgoer/openssl"
	"github.com/junaozun/game_server/utils"
)

type ReqBody struct {
	Seq   int64       `json:"seq"`
	Name  string      `json:"name"`
	Msg   interface{} `json:"msg"`
	Proxy string      `json:"proxy"` // 多进程，服务之间调用
}

type RespBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

type WsMsgReq struct {
	Body *ReqBody
	Conn WsConn
}

type WsMsgResp struct {
	Body *RespBody
}

type WsConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, bool)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

func (w *wsServer) Start() {
	// 启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			fmt.Println(msg)
		}
	}
}

func (w *wsServer) readMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			w.Close()
			log.Fatal(err)
		}
	}()

	// 先读到客户端发送过来的数据，然后进行处理，再发送回给客户端消息
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		// 收到消息后要解析消息 就是json格式
		// 1 data解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式：", err)
			continue
		}
		// 2 前端的消息  加密消息 进行解密
		secretKey, ok := w.GetProperty("secretKey")
		if !ok {
			log.Println("未设置secretKey值")
			continue
		}
		key := secretKey.(string)
		// 客户端传过来的数据是加密的，需要解密
		d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			log.Println("数据格式有误，解密失败：", err)
			// 出错后发起握手
			// w.Handshake()
			continue
		}
		data = d
		// 3.data 转为body
		reqBody := &ReqBody{}
		err = json.Unmarshal(data, reqBody)
		if err != nil {
			log.Println("数据解析失败", err)
			continue
		}
		// 获取到前端传递的数据了，拿上这些数据 去具体的业务进行处理
		req := &WsMsgReq{
			Body: reqBody,
			Conn: w,
		}
		rsp := &WsMsgResp{
			Body: &RespBody{
				Seq:  reqBody.Seq,
				Name: reqBody.Name,
			},
		}
		//  交给router处理业务
		w.router.Run(req, rsp)
		// 将结果返回给客户端
		w.outChan <- rsp
	}
	w.Close()
}

func (w *wsServer) Close() {
	_ = w.wsConn.Close()
}
