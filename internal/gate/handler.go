package gate

import (
	"sync"

	"github.com/junaozun/game_server/pkg/ws"
	"github.com/junaozun/gogopkg/logrusx"
)

type Handler struct {
	mutex      sync.RWMutex
	proxyMap   map[string]map[int64]*ProxyClient // 代理地址（logic or login)->客户端ID(cid)->该客户端的链接通道
	loginProxy string
	logicProxy string
}

func NewHandler() *Handler {
	return &Handler{
		proxyMap: make(map[string]map[int64]*ProxyClient),
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

func (h *Handler) GetProxyMap(proxyAddr string) map[int64]*ProxyClient {
	h.mutex.Lock()
	v, ok := h.proxyMap[proxyAddr]
	if !ok {
		v = make(map[int64]*ProxyClient)
		h.proxyMap[proxyAddr] = v
	}
	h.mutex.Unlock()
	return v
}

func (h *Handler) SetProxyMapValue(proxyAddr string, fieldKey int64, fieldValue *ProxyClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.proxyMap[proxyAddr][fieldKey] = fieldValue
}

func (h *Handler) DeleteCid(proxyAddr string, cid int64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.proxyMap[proxyAddr], cid)
}

func (h *Handler) OnPushClient(conn *ClientConn, body *ws.RespBody) {
	clientConn, ok := conn.GetProperty("clientConn")
	if !ok {
		logrusx.Log.WithFields(logrusx.Fields{}).Error("[handler] OnPushClient getProperty clientConn not found")
		return
	}
	client, ok := clientConn.(ws.IWsConn)
	if ok {
		client.Push(body.Router, body.Msg)
	}
}
