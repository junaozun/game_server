package pvp

import (
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
	log.Println("[PvpService]..................")
	return nil
}

func (p PvpService) Run() {
}
