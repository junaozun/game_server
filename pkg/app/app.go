package app

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// App 应用
type App struct {
	ctx         context.Context
	ctxCancel   func()
	name        string
	version     string
	onExitHook  func()
	stopTimeout time.Duration
	sigs        []os.Signal
	runners     []Runner
}

func New(opts ...Option) *App {
	app := &App{
		sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP},
		stopTimeout: time.Second * 30,
	}
	for _, opt := range opts {
		opt(app)
	}
	app.ctx, app.ctxCancel = context.WithCancel(context.Background())
	return app
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Version() string {
	return a.version
}

func (a *App) Run() error {
	if a.onExitHook != nil {
		defer a.onExitHook()
	}
	eg, errCtx := errgroup.WithContext(a.ctx)
	wg := sync.WaitGroup{}
	for _, v := range a.runners {
		srv := v
		eg.Go(func() error {
			// 监听到cancel 停止指令
			<-errCtx.Done()
			// 用一个新的ctx来作用于超时逻辑
			sctx, cancel := context.WithTimeout(context.Background(), a.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(errCtx)
		})
	}
	wg.Wait()

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, a.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-signalChan:
				err := a.Stop()
				if err != nil {
					log.Printf("failed to stop app: %v", err)
					return err
				}
				log.Println("app run interrupt.....")
			}
		}
	})
	// 捕获err
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
	return nil
}
