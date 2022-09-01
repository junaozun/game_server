package repo_impl_mysql

import (
	"context"

	"github.com/junaozun/game_server/internal/web/repo"
	common_model "github.com/junaozun/game_server/model"
	"gorm.io/gorm"
)

type AccountRepo struct {
	DB *gorm.DB
}

func NewAccountRepo(DB *gorm.DB) repo.IAccountRepo {
	return &AccountRepo{
		DB: DB,
	}
}

func (a *AccountRepo) GetAccountByName(ctx context.Context, name string) (*common_model.User, error) {
	res := &common_model.User{}
	err := a.DB.WithContext(ctx).Where(&common_model.User{Username: name}).Limit(1).Find(res).Error
	return res, err
}

func (a *AccountRepo) AddAccount(ctx context.Context, user *common_model.User) error {
	return a.DB.WithContext(ctx).Create(user).Error
}
