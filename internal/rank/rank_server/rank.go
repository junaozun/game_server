package rank

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/junaozun/game_server/internal/rank/nats_handler"
	"github.com/junaozun/gogopkg/app"
	"github.com/junaozun/gogopkg/config"
	"github.com/junaozun/gogopkg/logrusx"
	"github.com/junaozun/gogopkg/natsx"
	"google.golang.org/protobuf/proto"
)

const (
	rankSnapshotKey = "cross:rank:snapshot" // 排行榜快照数据，用于对比昨日排名
)

type Rank struct {
	ServerName        string
	snapshotRank      map[string]*RankSnapshot // 排行榜快照 rankKey -> id -> score
	dirtySnapshotRank map[string]struct{}      // 脏标记
	db                *Dao                     // redis 数据存储
	rpc
}

func NewRank() *Rank {
	r := &Rank{
		ServerName:        "rank",
		snapshotRank:      make(map[string]*RankSnapshot),
		dirtySnapshotRank: make(map[string]struct{}),
		db:                NewDao(),
	}
	return r
}

func (c *Rank) Run(cfg config.GameConfig) error {
	runners := make([]app.Runner, 0)
	rankServer := NewRank()
	rankServer.dirtySnapshotRank
	natsxServer := natsx.New(cfg.Common.NATS, c.ServerName)
	// 注册nats
	nats_handler.RegisterHandler(natsxServer)
	runners = append(runners, natsxServer)
	rank := app.New(
		app.OnBeginHook(func() {
			logrusx.Log.Info("rank app start .....")
		}),
		app.OnExitHook(func() {
			logrusx.Log.Info("rank app exit .....")
		}),
		app.Name(c.ServerName),
		app.Runners(runners...),
	)
	if err := rank.Run(); err != nil {
		return err
	}
	return nil
}

// Init 启动
func (r *Rank) Init(ctx context.Context) error {
	return r.loadSnapshotData()
}

func (r *Rank) Stop(ctx context.Context) error {
	return r.db.asyncRedis.Stop(ctx)
}

// 加载昨日排行榜快照
func (r *Rank) loadSnapshotData() error {
	saveSnapshotRanks, err := r.db.LoadHashAllStringBytesSync(rankSnapshotKey)
	if err != nil {
		return err
	}

	for rankKey, data := range saveSnapshotRanks {
		snapshotData := new(RankSnapshot)
		err := proto.Unmarshal([]byte(data), snapshotData)
		if err != nil {
			return err
		}
		r.snapshotRank[rankKey] = snapshotData
	}
	return nil
}

// AddRankScore 更新单个排行榜，score 增加
func (r *Rank) AddRankScore(rankKey string, rankData map[string]int) {
	r.db.asyncRedis.ZIncrby(rankKey, rankData)
}

// UpdateRankScore 更新单个排行榜，score直接替换
func (r *Rank) UpdateRankScore(rankKey string, rankData map[string]float64) {
	r.db.asyncRedis.ZAdd(rankKey, rankData)
}

// DeleteRankData 删除排行榜中指定成员
func (r *Rank) DeleteRankData(rankKey string, mems ...string) {
	r.db.asyncRedis.ZRem(rankKey, mems...)
}

type RankResult struct {
	RankList []*RankItem
	Me       *RankItem
	Total    int
}

// GetRank 获取排行榜
func (r *Rank) GetRank(me string, rankKey string, start int64, count int64, cb func(*RankResult)) {
	meScore, err := r.db.asyncRedis.Sync().ZRevRank(rankKey, me)
	if err != nil {
		return
	}
	meRank, err := r.db.asyncRedis.Sync().ZRevRank(rankKey, me)
	if err != nil {
		return
	}

	result := &RankResult{
		Me: &RankItem{
			Id:    me,
			Score: meScore,
			Rank:  uint32(meRank),
		},
	}
	r.db.asyncRedis.ZRevRange(rankKey, start, start+count-1, func(res []redis.Z, err error) {
		for rank, data := range res {
			result.RankList = append(result.RankList, &RankItem{
				Id:    data.Member.(string),
				Score: int64(data.Score),
				Rank:  uint32(rank),
			})
		}
		cb(result)
		return
	})
}
