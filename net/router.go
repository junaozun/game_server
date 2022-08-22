package net

import (
	"strings"
)

type Router struct {
	group map[string]*group // key:组标识，即prefix
}

func NewRouter() *Router {
	return &Router{
		group: make(map[string]*group),
	}
}

func (r *Router) Group(prefix string) *group {
	if v, ok := r.group[prefix]; ok {
		return v
	}
	g := &group{
		handlerMap: make(map[string]HandlerFunc),
	}
	r.group[prefix] = g
	return g
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgResp) {
	// req.Body.Name 路径 登录业务 account.login (account 组标识)(login 路由标识)
	strs := strings.Split(req.Body.Router, ".")
	if len(strs) != 2 {
		return
	}
	prefix := strs[0]
	name := strs[1]
	g, ok := r.group[prefix]
	if !ok {
		return
	}
	g.exec(name, req, rsp)
}

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgResp)

type group struct {
	handlerMap map[string]HandlerFunc // key:路由标识
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgResp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req, rsp)
	}
}

func (g *group) AddRouter(name string, handlerFunc HandlerFunc) {
	g.handlerMap[name] = handlerFunc
}
