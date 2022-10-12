package gate

import (
	"sync"

	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/gogopkg/logrusx"
)

type Handler struct {
	mutex      sync.RWMutex
	proxyMap   map[string]map[int]*ProxyClient // 代理地址（logic or login)->客户端ID->该客户端的链接通道
	loginProxy string
	logicProxy string
}

func NewHandler() *Handler {
	return &Handler{
		proxyMap: make(map[string]map[int]*ProxyClient),
	}
}

func (h *Handler) SetLogicProxy(proxy string) {
	h.logicProxy = proxy
}

func (h *Handler) SetLoginProxy(proxy string) {
	h.loginProxy = proxy
}

func (h *Handler) GetLoginProxy() string {
	return h.loginProxy
}

func (h *Handler) GetLogicProxy() string {
	return h.logicProxy
}

func (h *Handler) GetProxyMap(proxyAddr string) map[int]*ProxyClient {
	h.mutex.Lock()
	v, ok := h.proxyMap[proxyAddr]
	if !ok {
		v = make(map[int]*ProxyClient)
		h.proxyMap[proxyAddr] = v
	}
	h.mutex.Unlock()
	return v
}

func (h *Handler) SetProxyMapValue(proxyAddr string, fieldKey int, fieldValue *ProxyClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.proxyMap[proxyAddr][fieldKey] = fieldValue
}

func (h *Handler) DeleteCid(proxyAddr string, cid int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.proxyMap[proxyAddr], cid)
}

func (h *Handler) OnPush(conn *ClientConn, body *ws.RespBody) {
	clientConn, ok := conn.GetProperty("clientConn")
	if !ok {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[handler] OnPush getProperty clientConn not found")
		return
	}
	cc, ok := clientConn.(ws.IWsConn)
	cc.Push(body.Router, body.Msg)
}
