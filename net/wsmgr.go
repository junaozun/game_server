package net

import "sync"

type WsMgr struct {
	mu        sync.RWMutex
	userCache map[int]WsConn // key:uid ;value:用户连接
}

func NewWsMgr() *WsMgr {
	return &WsMgr{
		userCache: make(map[int]WsConn),
	}
}

func (w *WsMgr) UserLogin(conn WsConn, uid int, token string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	oldConn := w.userCache[uid]
	// 用户还在登录
	if oldConn != nil {
		if conn != oldConn { // 用户重新登录了,通知客户端
			oldConn.Push("robLogin", nil)
		}
	}
	w.userCache[uid] = conn
	conn.SetProperty("uid", uid)
	conn.SetProperty("token", token)
}
