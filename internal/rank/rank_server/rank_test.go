package rank_server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/junaozun/gogopkg/logrusx"
)

var r *Rank

const (
	rankKey = "area"
)

func TestMain(m *testing.M) {
	r = NewRank()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	ticker := time.NewTimer(time.Minute)
	defer ticker.Stop()
	// 主循环
QUIT:
	for {
		select {
		case sig := <-sigs:
			r.Stop(context.Background())
			log.Printf("Signal: %s", sig.String())
			break QUIT
		case <-ticker.C:
			logrusx.Log.WithFields(logrusx.Fields{
				"goroutine count": runtime.NumGoroutine(),
			}).Info("协程数量")
		}
	}
	logrusx.Log.Info("[main] quiting......")
	os.Exit(m.Run())
}

func TestGetRankData(t *testing.T) {
	r.AddRankScore(rankKey, map[string]int{
		"nihao1":  20,
		"nihao2":  1,
		"nihao3":  4,
		"nihao4":  98,
		"nihao99": 48,
	})
	r.GetRank("nihao3", rankKey, 0, 10, func(result *RankResult) {
		fmt.Println(result)
	})
}
