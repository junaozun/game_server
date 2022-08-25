package cross

import (
	"context"
	"flag"
	"log"

	pkgConfig "github.com/junaozun/game_server/pkg/config"
)

type CrossService struct {
}

func NewCrossService() *CrossService {
	crossService := &CrossService{}
	return crossService
}

func (c *CrossService) ParseFlag(set *flag.FlagSet) {
}

func (c *CrossService) Init(cfg pkgConfig.GameConfig) error {
	log.Println("[CrossService] init successful .....")
	return nil
}

func (c *CrossService) Start(ctx context.Context) error {
	log.Println("[CrossService] start .....")
	return nil
}

func (c *CrossService) Stop(ctx context.Context) error {
	log.Println("[CrossService] stop .....")
	return nil
}
