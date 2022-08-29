package httpx

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/junaozun/game_server/pkg/httpx/middleware"
)

type HttpxServer struct {
	*http.Server
	opt *Option
	ls  net.Listener
}

// New 构造http server
func New(handler func(engine *gin.Engine), opts ...ServerOption) (*HttpxServer, error) {
	var (
		ls  net.Listener
		err error
	)
	opt := &Option{
		Address: "0.0.0.0:80",
		NoMethodHandler: func(c *gin.Context) {
			c.JSON(http.StatusNotFound,
				gin.H{
					"code":    404,
					"message": "no method found",
				})
		},
		NoRouterHandler: func(c *gin.Context) {
			c.JSON(http.StatusNotFound,
				gin.H{
					"code":    404,
					"message": "no route found",
				})
		},
		RecoverHandler: middleware.Recover(),
		PProf:          true,
	}
	for _, v := range opts {
		v(opt)
	}
	if opt.Listener == nil {
		ls, err = net.Listen("tcp", opt.Address)
		if err != nil {
			return nil, err
		}
	} else {
		ls = opt.Listener
	}
	router := gin.New()
	if gin.Mode() == gin.DebugMode {
		router.Use(gin.Logger())
		gin.ForceConsoleColor()
	}
	router.NoRoute(opt.NoRouterHandler)
	router.NoMethod(opt.NoMethodHandler)
	router.Use(opt.RecoverHandler)
	if opt.PProf {
		pprof.Register(router)
	}
	handler(router)
	h := &HttpxServer{
		Server: &http.Server{
			Handler: router,
		},
		opt: opt,
		ls:  ls,
	}
	return h, nil
}

// Start 运行
func (h *HttpxServer) Start(ctx context.Context) error {
	log.Printf("[HttpxServer] addr:%s,start .....", h.opt.Address)
	h.Server.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}
	// 不负责外面传来的listener的关闭
	if h.opt.Listener == nil {
		defer h.ls.Close()
	}
	err := h.Serve(h.ls)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop 停止
func (h *HttpxServer) Stop(ctx context.Context) error {
	log.Printf("[HttpxServer] addr:%s,Stop .....", h.opt.Address)
	return h.Shutdown(ctx)
}
