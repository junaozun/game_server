package app

import (
	"context"
	"os"
)

type Runner interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Option func(o *App)

func Name(name string) Option {
	return func(o *App) {
		o.name = name
	}
}

// OnExitHook 全部退出钩子函数
func OnExitHook(hook func()) Option {
	return func(o *App) { o.onExitHook = hook }
}

func Version(version string) Option {
	return func(o *App) {
		o.version = version
	}
}

func Context(ctx context.Context) Option {
	return func(o *App) {
		o.ctx = ctx
	}
}

// Signal 信号
func Signal(sigs ...os.Signal) Option {
	return func(o *App) { o.sigs = sigs }
}

func Runners(servers ...Runner) Option {
	return func(o *App) {
		o.runners = servers
	}
}
