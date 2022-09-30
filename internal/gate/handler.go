package gate

import (
	"sync"
)

type Handler struct {
	mutex      sync.Mutex
	proxyMap   map[string]map[int]*ProxyClient // 代理地址（ws:0.0.0.0:8003）:客户端ID：该客户端的链接通道
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
