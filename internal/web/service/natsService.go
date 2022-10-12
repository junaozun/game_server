package service

import (
	"context"

	"github.com/junaozun/game_server/api"
	"github.com/junaozun/game_server/common"
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
	err := n.Repo.Request(ctx, common.ServerName_Logic, "TestReqServer", "AddTestMine", req, resp, natsx.WithCallID(100001))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n *NatsService) GetRankTest(ctx context.Context, rankKey string) (*api.GetRankResp, error) {
	req := &api.GetRankReq{
		RankKey:   "area",
		Me:        "nihao3",
		BeginRank: 0,
		Count:     10,
	}
	resp := &api.GetRankResp{}
	err := n.Repo.Request(ctx, common.ServerName_Rank, "RankHandler", "OnGetRank", req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
