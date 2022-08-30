package app

import (
	"context"
	"errors"

	pkgConfig "github.com/junaozun/game_server/pkg/config"
	"golang.org/x/sync/errgroup"
)

type IApp interface {
	Run(ctx context.Context, cfg pkgConfig.GameConfig) error
}

type AppManager struct {
	Apps []IApp
}

func NewAppMgr(apps ...IApp) *AppManager {
	return &AppManager{
		Apps: apps,
	}
}

func (a *AppManager) Runs(cfg pkgConfig.GameConfig) {
	eg, errCtx := errgroup.WithContext(context.Background())
	// 将所有的App实例运行起来
	for _, v := range a.Apps {
		app := v
		eg.Go(func() error {
			return app.Run(errCtx, cfg)
		})
	}
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}
