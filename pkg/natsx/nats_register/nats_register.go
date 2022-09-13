package nats_register

import (
	"github.com/junaozun/game_server/pkg/config"
	"github.com/junaozun/game_server/pkg/natsx"
)

type RegisterNatsFunc func(serverName string, svc interface{}) error

func RegisterNats(natsCfg *config.NatsConfig) RegisterNatsFunc {
	connEnc, err := natsx.NewNatsJSONEnc(natsCfg.Server) // nats.MaxReconnects(int(natsCfg.MaxReconnects)),
	// nats.ReconnectWait(time.Duration(natsCfg.ReconnectWait)),
	// nats.Timeout(time.Duration(natsCfg.RequestTimeout)),

	if err != nil {
		panic(err)
	}
	server, err := natsx.NewServer(connEnc)
	if err != nil {
		panic(err)
	}
	return func(serverName string, svc interface{}) error {
		return server.Register(serverName, svc)
	}
}
