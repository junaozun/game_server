package rank_server

import (
	"context"

	"github.com/junaozun/game_server/internal/rank/rank_api"
	"github.com/junaozun/gogopkg/natsx/testdata"
)

type RankHandler struct {
	Rank *Rank
}

func (r *RankHandler) OnGetRank(ctx context.Context, req *rank_api.GetRankReq) (*rank_api.GetRankResp, error) {
	resp := &rank_api.GetRankResp{}
	var err error
	ch := make(chan struct{})
	r.Rank.GetRank(req.GetMe(), req.GetRankKey(), int64(req.GetBeginRank()), int64(req.GetCount()), func(res *RankResult) {
		resRankItem := make([]*rank_api.RankItem, 0, len(res.RankList))
		for _, v := range res.RankList {
			resRankItem = append(resRankItem, &rank_api.RankItem{
				Id:      v.id,
				Score:   v.score,
				Rank:    v.rank,
				OldRank: v.oldRank,
			})
		}
		resp.RankItem = resRankItem
		resp.Me = &rank_api.RankItem{
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
