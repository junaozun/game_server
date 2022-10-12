package rank_server

import (
	"context"

	"github.com/junaozun/game_server/api"
	"github.com/junaozun/gogopkg/natsx/testdata"
)

type RankHandler struct {
	Rank *Rank
}

func (r *RankHandler) OnGetRank(ctx context.Context, req *api.GetRankReq) (*api.GetRankResp, error) {
	resp := &api.GetRankResp{}
	var err error
	ch := make(chan struct{})
	r.Rank.GetRank(req.GetMe(), req.GetRankKey(), 0, 10, func(res *RankResult) {
		resRankItem := make([]*api.RankItem, 0, len(res.RankList))
		for _, v := range res.RankList {
			resRankItem = append(resRankItem, &api.RankItem{
				Id:      v.id,
				Score:   v.score,
				Rank:    v.rank,
				OldRank: v.oldRank,
			})
		}
		resp.RankItem = resRankItem
		resp.Me = &api.RankItem{
			Id:      res.Me.id,
			Score:   res.Me.score,
			Rank:    res.Me.rank,
			OldRank: res.Me.oldRank,
		}
		resp.Total = int64(res.Total)
		ch <- struct{}{}
	})
	<-ch
	return resp, err
}

func (r *RankHandler) OnAddRankScore(ctx context.Context, req *testdata.TestMine) {

}

func (r *RankHandler) OnUpdateRankScore(ctx context.Context, req *testdata.TestMine) {

}

func (r *RankHandler) OnDeleteRankData(ctx context.Context, req *testdata.TestMine) {

}
