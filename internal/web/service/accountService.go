package service

import (
	"context"

	"github.com/junaozun/game_server/internal/web/repo"
	common_model "github.com/junaozun/game_server/model"
	"github.com/junaozun/game_server/pkg/errno"
	"github.com/junaozun/game_server/ret"
)

type AccountService struct {
	Repo repo.IAccountRepo
}

func NewAccountService(accountRepo repo.IAccountRepo) *AccountService {
	return &AccountService{
		Repo: accountRepo,
	}
}

func (a *AccountService) GetAccountByName(ctx context.Context, name string) (*common_model.User, errno.Err) {
	res, err := a.Repo.GetAccountByName(ctx, name)
	if err != nil {
		return nil, ret.Err_DB
	}
	return res, ret.OK
}

func (a *AccountService) AddAccount(ctx context.Context, user *common_model.User) errno.Err {
	account, errno := a.GetAccountByName(ctx, user.Username)
	if errno.Code != ret.OK.Code {
		return errno
	}
	// 用户已存在
	if account.UId != 0 {
		return ret.Err_UserExist
	}
	err := a.Repo.AddAccount(ctx, user)
	if err != nil {
		return ret.Err_DB
	}
	return ret.OK
}
