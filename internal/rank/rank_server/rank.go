package rank_server

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/junaozun/game_server/internal/rank/rank_data"
	"github.com/junaozun/gogopkg/logrusx"
	"google.golang.org/protobuf/proto"
)

const (
	rankSnapshotKey = "rank:snapshot" // 排行榜快照数据，用于对比昨日排名
)

type Rank struct {
	snapshotRank      map[string]*rank_data.RankSnapshot // 排行榜快照 rankKey -> id -> score
	dirtySnapshotRank map[string]struct{}                // 脏标记
	db                *Dao                               // redis 数据存储
}

func NewRank() *Rank {
	r := &Rank{
		snapshotRank:      make(map[string]*rank_data.RankSnapshot),
		dirtySnapshotRank: make(map[string]struct{}),
		db:                NewDao(),
	}
	return r
}

// Init 启动
func (r *Rank) Start(ctx context.Context) error {
	return r.loadSnapshotData()
}

func (r *Rank) Stop(ctx context.Context) error {
	return r.db.asyncClient.Stop(ctx)
}

// 加载昨日排行榜快照
func (r *Rank) loadSnapshotData() error {
	saveSnapshotRanks, err := r.db.LoadHashAllStringBytesSync(rankSnapshotKey)
	if err != nil {
		return err
	}

	for rankKey, data := range saveSnapshotRanks {
		snapshotData := new(rank_data.RankSnapshot)
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
	r.db.asyncClient.ZIncrby(rankKey, rankData)
}

// UpdateRankScore 更新单个排行榜，score直接替换
func (r *Rank) UpdateRankScore(rankKey string, rankData map[string]float64) {
	r.db.asyncClient.ZAdd(rankKey, rankData)
}

// DeleteRankData 删除排行榜中指定成员
func (r *Rank) DeleteRankData(rankKey string, mems ...string) {
	r.db.asyncClient.ZRem(rankKey, mems...)
}

type RankResult struct {
	RankList []*rankItem
	Me       *rankItem
	Total    int
}

type rankItem struct {
	id      string
	score   int64
	rank    uint32
	oldRank uint32
}

// GetRank 获取排行榜
func (r *Rank) GetRank(me string, rankKey string, start int64, count int64, cb func(*RankResult)) {
	result := &RankResult{}
	r.db.asyncClient.ZRevRank(rankKey, me, func(meRank int64, err error) {
		if err != nil {
			logrusx.Log.WithFields(logrusx.Fields{
				"me": me,
			}).Error("zrevrank err")
			cb(nil)
			return
		}
		merank := meRank
		r.db.asyncClient.ZScore(rankKey, me, func(meSocre float64, err error) {
			if err != nil {
				logrusx.Log.WithFields(logrusx.Fields{
					"me": me,
				}).Error("zscore err")
				cb(nil)
				return
			}
			mescore := meSocre
			r.db.asyncClient.ZRevRange(rankKey, start, start+count-1, func(res []redis.Z, err error) {
				for rank, data := range res {
					result.RankList = append(result.RankList, &rankItem{
						id:    data.Member.(string),
						score: int64(data.Score),
						rank:  uint32(rank),
					})
				}
				result.Me = &rankItem{
					id:    me,
					score: int64(mescore),
					rank:  uint32(merank),
				}
				cb(result)
				return
			})
		})
	})
}

// SaveSnapshotRank 存排行榜快照
func (r *Rank) SaveSnapshotRank() {

}
