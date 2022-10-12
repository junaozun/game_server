package repo_impl_nats

import (
	"context"
	"fmt"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/web/repo"
	"github.com/junaozun/gogopkg/natsx"
)

type NatsRepo struct {
	commponent *component.Component
}

func NewNatsRepo(component *component.Component) repo.INatsRepo {
	return &NatsRepo{
		commponent: component,
	}
}

func (n *NatsRepo) Publish(objectName string, serverName string, methodName string, req interface{}, opt ...natsx.CallOption) error {
	client, ok := n.commponent.NatsClient[serverName]
	if !ok {
		return fmt.Errorf("not fount natsClient %s", serverName)
	}
	return client.Publish(objectName, methodName, req, opt...)
}

func (n *NatsRepo) Request(ctx context.Context, serverName string, objectName string, methodName string, req interface{}, resp interface{}, opt ...natsx.CallOption) error {
	client, ok := n.commponent.NatsClient[serverName]
	if !ok {
		return fmt.Errorf("not fount natsClient %s", serverName)
	}
	return client.Request(ctx, objectName, methodName, req, resp, opt...)
}
