package service

import (
	"context"

	"github.com/junaozun/game_server/internal/web/repo"
	"github.com/junaozun/gogopkg/natsx"
	"github.com/junaozun/gogopkg/natsx/testdata"
)

type NatsService struct {
	Repo repo.INatsRepo
}

func NewNatsService(natsRepo repo.INatsRepo) *NatsService {
	return &NatsService{
		Repo: natsRepo,
	}
}

func (n *NatsService) UseNatsTest(ctx context.Context, name string) (*testdata.TestMineResp, error) {
	req := &testdata.TestMine{
		Id:   1001,
		Name: name,
		Sex:  "woman",
	}
	resp := &testdata.TestMineResp{}
	err := n.Repo.Request(ctx, "TestReqServer", "AddTestMine", req, resp, natsx.WithCallID(100001))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
