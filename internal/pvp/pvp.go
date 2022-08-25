package pvp

import (
	"context"
	"flag"
	"log"

	pkgConfig "github.com/junaozun/game_server/pkg/config"
)

type PvpService struct {
}

func NewPvpService() *PvpService {
	pvpService := &PvpService{}
	return pvpService
}

func (p PvpService) ParseFlag(set *flag.FlagSet) {
}

func (p PvpService) Init(cfg pkgConfig.GameConfig) error {
	log.Println("[PvpService] init successful ....")
	return nil
}

func (p PvpService) Start(ctx context.Context) error {
	log.Println("[PvpService] start ....")
	return nil
}

func (p *PvpService) Stop(ctx context.Context) error {
	log.Println("[PvpService] stop ....")
	return nil
}
