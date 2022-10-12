package rank

import (
	"context"

	"github.com/junaozun/gogopkg/redisx"
)

type Dao struct {
	asyncRedis *redisx.AsyncClient
}

func NewDao() *Dao {
	cfg := redisx.Config{
		Server: "127.0.0.1:6379",
		Index:  0,
	}
	var err error
	client, err := redisx.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// callBack need
	pushFunc := func(ctx context.Context, f func()) error {
		f()
		return nil
	}
	asyncClient := redisx.NewAsync(client, redisx.WithPushFunc(pushFunc))
	return &Dao{
		asyncRedis: asyncClient,
	}
}

func (d *Dao) LoadHashAllStringBytesSync(key string) (map[string]string, error) {
	res, err := d.asyncRedis.Sync().HGetAll(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}
