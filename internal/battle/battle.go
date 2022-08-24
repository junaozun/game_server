package battle

import (
	"flag"
	"log"

	pkgConfig "github.com/junaozun/game_server/pkg/config"
)

type BattleService struct {
}

func NewBattleService() *BattleService {
	battleService := &BattleService{}
	return battleService
}

func (b BattleService) ParseFlag(set *flag.FlagSet) {
}

func (b BattleService) Init(cfg pkgConfig.GameConfig) error {
	log.Println("[battleService]..................")
	return nil
}

func (b BattleService) Run() {
}
