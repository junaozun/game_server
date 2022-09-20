package natsx

import (
	"context"
	"log"

	"github.com/junaozun/game_server/pkg/config"
)

type NatsxServer struct {
	ServerName string
	*Server
}

func New(natsCfg *config.NatsConfig, serverName string) *NatsxServer {
	connEnc, err := NewNatsPBEnc(natsCfg.Server) // nats.MaxReconnects(int(natsCfg.MaxReconnects)),
	// nats.ReconnectWait(time.Duration(natsCfg.ReconnectWait)),
	// nats.Timeout(time.Duration(natsCfg.RequestTimeout)),

	if err != nil {
		panic(err)
	}
	server, err := NewServer(connEnc)
	if err != nil {
		panic(err)
	}
	return &NatsxServer{
		Server:     server,
		ServerName: serverName,
	}
}

func (n *NatsxServer) Start(ctx context.Context) error {
	log.Printf("[NatsxServer] %s Start success", n.ServerName)
	select {
	case <-ctx.Done():
		return nil
	}
	return nil
}

func (n *NatsxServer) Stop(ctx context.Context) error {
	return n.Close(ctx)
}
