package repo_impl_nats

import (
	"context"
	"fmt"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/web/repo"
	"github.com/junaozun/game_server/pkg/natsx"
)

type NatsRepo struct {
	commponent *component.Component
}

func NewNatsRepo(component *component.Component) repo.INatsRepo {
	return &NatsRepo{
		commponent: component,
	}
}
func (n *NatsRepo) Publish(objectName string, methodName string, req interface{}, opt ...natsx.CallOption) error {
	return n.commponent.LogicClient.Publish(objectName, methodName, req, opt...)
}
func (n *NatsRepo) Request(ctx context.Context, objectName string, methodName string, req interface{}, resp interface{}, opt ...natsx.CallOption) error {
	return n.commponent.LogicClient.Request(ctx, objectName, methodName, req, resp, opt...)
}

func (n *NatsRepo) NatsI() {
	fmt.Println("唯一实现，不可调用")
}
