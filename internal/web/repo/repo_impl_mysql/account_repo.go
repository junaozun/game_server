package repo_impl_mysql

import (
	"context"

	"github.com/junaozun/game_server/component"
	"github.com/junaozun/game_server/internal/web/repo"
	common_model "github.com/junaozun/game_server/model"
)

type AccountRepo struct {
	commponent *component.Component
}

func NewAccountRepo(component *component.Component) repo.IAccountRepo {
	return &AccountRepo{
		commponent: component,
	}
}

func (a *AccountRepo) GetAccountByName(ctx context.Context, name string) (*common_model.User, error) {
	res := &common_model.User{}
	err := a.commponent.Dao.DB.WithContext(ctx).Where(&common_model.User{Username: name}).Limit(1).Find(res).Error
	return res, err
}

func (a *AccountRepo) AddAccount(ctx context.Context, user *common_model.User) error {
	return a.commponent.Dao.DB.WithContext(ctx).Create(user).Error
}
