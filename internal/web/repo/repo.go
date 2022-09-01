package repo

import (
	"context"

	common_model "github.com/junaozun/game_server/model"
)

type IAccountRepo interface {
	GetAccountByName(ctx context.Context, name string) (*common_model.User, error)
	AddAccount(ctx context.Context, user *common_model.User) error
}

type IGMRepo interface {
}
