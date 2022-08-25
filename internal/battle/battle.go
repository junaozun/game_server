package battle

import (
	"context"
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

func (b *BattleService) ParseFlag(set *flag.FlagSet) {
}

func (b *BattleService) Init(cfg pkgConfig.GameConfig) error {
	log.Println("[battleService] init successful ......")
	return nil
}

func (b *BattleService) Start(ctx context.Context) error {
	log.Println("[battleService] start .......")
	return nil
}

func (c *BattleService) Stop(ctx context.Context) error {
	log.Println("[battleService] stop .......")
	return nil
}
