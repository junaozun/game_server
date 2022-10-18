package wsMgr

import (
	"sync"

	"github.com/junaozun/game_server/pkg/ws"
)

type WsMgr struct {
	mu        sync.RWMutex
	userCache map[uint64]ws.IWsConn // key:uid ;value:用户连接
}

func NewWsMgr() *WsMgr {
	return &WsMgr{
		userCache: make(map[uint64]ws.IWsConn),
	}
}

func (w *WsMgr) UserLogin(conn ws.IWsConn, uid uint64, token string) {
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
