package httpx

import (
	"net"

	"github.com/gin-gonic/gin"
)

type Option struct {
	Address         string          // listen地址 ":80"
	Listener        net.Listener    // listener，非空使用这个listener
	NoMethodHandler gin.HandlerFunc // 无该方法handler
	NoRouterHandler gin.HandlerFunc // 无该路由handler
	RecoverHandler  gin.HandlerFunc // recover handler
	PProf           bool            // 是否开启pprof
}

type ServerOption func(option *Option)

func WithAddress(address string) ServerOption {
	return func(option *Option) {
		option.Address = address
	}
}

// WithListener listener
func WithListener(ls net.Listener) ServerOption {
	return func(option *Option) {
		option.Listener = ls
	}
}

// WithNoMethodHandler 无该方法的handler
func WithNoMethodHandler(h gin.HandlerFunc) ServerOption {
	return func(option *Option) {
		option.NoMethodHandler = h
	}
}

// WithNoRouteHandler 无路由的handler
func WithNoRouteHandler(h gin.HandlerFunc) ServerOption {
	return func(option *Option) {
		option.NoRouterHandler = h
	}
}

// WithRecoverHandler 无路由的handler
func WithRecoverHandler(h gin.HandlerFunc) ServerOption {
	return func(option *Option) {
		option.RecoverHandler = h
	}
}

// WithPProf 是否开启PProf
func WithPProf(enablePProf bool) ServerOption {
	return func(option *Option) {
		option.PProf = enablePProf
	}
}
