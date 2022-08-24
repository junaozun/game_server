package cross

import (
	"errors"
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

func (c CrossService) ParseFlag(set *flag.FlagSet) {
}

func (c CrossService) Init(cfg pkgConfig.GameConfig) error {
	log.Println("[CrossService]..................")
	return errors.New("niaho")
}

func (c CrossService) Run() {
}
