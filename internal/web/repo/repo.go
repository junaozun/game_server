package repo

import (
	"context"

	common_model "github.com/junaozun/game_server/model"
	"github.com/junaozun/game_server/pkg/natsx"
)

type IAccountRepo interface {
	GetAccountByName(ctx context.Context, name string) (*common_model.User, error)
	AddAccount(ctx context.Context, user *common_model.User) error
}

type INatsRepo interface {
	Publish(objectName string, methodName string, req interface{}, opt ...natsx.CallOption) error
	Request(ctx context.Context, objectName string, methodName string, req interface{}, resp interface{}, opt ...natsx.CallOption) error
	NatsI()
}
