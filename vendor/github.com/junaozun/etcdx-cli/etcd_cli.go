package etcdx_cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var _ IClient = (*Client)(nil)

// Client etcd 客户端
type Client struct {
	once      sync.Once
	closeChan chan struct{}
	client    *clientv3.Client // etcd v3 client
	cfg       Config
	wg        sync.WaitGroup
}

// Config 配置
type Config struct {
	Servers        string `xml:"servers" yaml:"servers"`
	DialTimeout    int64  `xml:"dial_timeout" yaml:"dial_timeout"`
	RequestTimeout int64  `xml:"request_timeout" yaml:"request_timeout"`
}

func defaultConfig(cfg *Config) {
	if 0 == cfg.RequestTimeout {
		cfg.RequestTimeout = 1
	}
	if 0 == cfg.DialTimeout {
		cfg.DialTimeout = 1
	}
}

// NewClient 构造一个注册服务
func NewClient(client *clientv3.Client) (*Client, error) {
	cfg := Config{}
	defaultConfig(&cfg)
	s := &Client{
		closeChan: make(chan struct{}),
		client:    client,
		cfg:       cfg,
	}
	return s, nil
}

func NewClientWithConfig(cfg Config) (*Client, error) {
	defaultConfig(&cfg)
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(cfg.Servers, ","),
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
	})

	if nil != err {
		return nil, err
	}
	cli, err := NewClient(client)
	if nil != err {
		return nil, err
	}
	cli.cfg = cfg
	return cli, nil
}

type KV struct {
	Key   string
	Value string
}

func NewKV(k string, v string) KV {
	return KV{
		Key:   k,
		Value: v,
	}
}

func (c *Client) KV(ctx context.Context, key string) (map[string][]byte, error) {
	resp, err := c.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	ret := make(map[string][]byte, len(resp.Kvs))
	for _, v := range resp.Kvs {
		ret[string(v.Key)] = v.Value
	}
	return ret, nil
}

func (c *Client) Put(ctx context.Context, key string, value string) error {
	return c.PutWithTTL(ctx, key, value, 0)
}

func (c *Client) PutJSON(ctx context.Context, key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Put(ctx, key, string(b))
}
func (c *Client) PutWithTTL(ctx context.Context, key, value string, ttl int64) error {
	var opt []clientv3.OpOption
	if ttl > 0 {
		lease := clientv3.NewLease(c.client)
		grantResp, err := lease.Grant(ctx, ttl)
		if err != nil {
			return err
		}
		opt = append(opt, clientv3.WithLease(grantResp.ID))
	}

	_, err := c.client.Put(ctx, key, value, opt...)
	return err
}

func (c *Client) PutJSONWithTTL(ctx context.Context, key string, data interface{}, ttl int64) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.PutWithTTL(ctx, key, string(b), ttl)
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	ret, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(ret.Kvs) == 0 {
		return nil, nil
	}
	return ret.Kvs[0].Value, nil
}

func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	ret, err := c.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func (c *Client) GetAllKey(ctx context.Context, key string) ([]string, error) {
	ret, err := c.KV(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	tmp := make([]string, 0, len(ret))
	for k, _ := range ret {
		tmp = append(tmp, k)
	}
	sort.Strings(tmp)
	return tmp, nil
}

func (c *Client) GetAllValue(ctx context.Context, key string) ([][]byte, error) {
	ret, err := c.KV(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	tmp := make([][]byte, 0, len(ret))
	for _, v := range ret {
		tmp = append(tmp, v)
	}
	return tmp, nil
}

func (c *Client) GetJSON(ctx context.Context, key string, data interface{}) error {
	b, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, data)
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.client.Delete(ctx, key)
	return err
}

// PutIfNotExist 如果create不存在,放入键值create跟kv
func (c *Client) PutIfNotExist(ctx context.Context, create KV, kv ...KV) (bool, error) {
	// 比较Revision, 当key不存在时，createRevision是0
	keyNotExist := clientv3.Compare(clientv3.CreateRevision(create.Key), "=", 0)

	var puts = make([]clientv3.Op, 0, len(kv)+1)
	puts = append(puts, clientv3.OpPut(create.Key, create.Value))
	for _, v := range kv {
		put := clientv3.OpPut(v.Key, v.Value)
		puts = append(puts, put)
	}

	resp, err := c.client.Txn(ctx).If(keyNotExist).Then(puts...).Commit()
	if err != nil {
		return false, err
	}
	if !resp.Succeeded {
		return false, nil
	}
	return true, nil
}

// UpdateTx 更新事务
func (c *Client) UpdateTx(ctx context.Context, key, val string) error {
	// 比较Revision, 当key不存在时，createRevision是0
	keyExist := clientv3.Compare(clientv3.CreateRevision(key), "!=", 0)
	put := clientv3.OpPut(key, val)
	resp, err := c.client.Txn(ctx).If(keyExist).Then(put).Commit()
	if err != nil {
		return err
	}
	if !resp.Succeeded {
		return errors.New("update failed")
	}
	return nil
}

// KeepAlive 保活
func (c *Client) KeepAlive(ctx context.Context, key string, val string, ttl int64) (err error) {
	c.wg.Add(1)
	defer c.wg.Done()

	ctxGrant, cancelGrant := context.WithTimeout(ctx, time.Second*time.Duration(c.cfg.RequestTimeout))
	defer cancelGrant()
	leaseResp, err := c.client.Grant(ctxGrant, ttl)
	if err != nil {
		return err
	}

	id := leaseResp.ID

	defer func() {
		select {
		case <-c.client.Ctx().Done():
			return
		default:
		}
		ctxRevoke, cancelGranRevoke := context.WithTimeout(context.Background(), time.Second*time.Duration(c.cfg.RequestTimeout))
		_, err = c.client.Revoke(ctxRevoke, id)
		cancelGranRevoke()
	}()

	ctxPut, cancelPut := context.WithTimeout(ctx, time.Second*time.Duration(c.cfg.RequestTimeout))
	defer cancelPut()
	_, err = c.client.Put(ctxPut, key, val, clientv3.WithLease(id))
	if err != nil {
		return err
	}

	laChan, err := c.client.KeepAlive(ctx, id)
	if err != nil {
		return err
	}
	for {
		select {
		case _, ok := <-laChan:
			if !ok {
				return
			}
		case <-c.client.Ctx().Done():
			return
		case <-ctx.Done():
			return
		case <-c.closeChan:
			return
		}
	}
}

// WatchEventType 事件类型
type WatchEventType int

const (
	Init   WatchEventType = iota // 初始
	Create                       // 创建
	Modify                       // 修改
	Delete                       // 删除
)

type WatchEvent struct {
	Revision int64
	KV       map[WatchEventType]map[string][]byte
}

// Watch 监听
func (c *Client) Watch(ctx context.Context, key string, prefix bool) (<-chan *WatchEvent, error) {
	var ops []clientv3.OpOption
	if prefix {
		ops = append(ops, clientv3.WithPrefix())
	}
	resp, err := c.client.Get(ctx, key, ops...)
	if err != nil {
		return nil, err
	}

	c.wg.Add(1)
	defer c.wg.Done()

	retChan := make(chan *WatchEvent)
	var initWe *WatchEvent
	// init
	if len(resp.Kvs) > 0 {
		initWe = &WatchEvent{
			Revision: resp.Header.Revision,
		}
		kvs := make(map[string][]byte, len(resp.Kvs))
		for _, ev := range resp.Kvs {
			kvs[string(ev.Key)] = ev.Value
		}
		initWe.KV = map[WatchEventType]map[string][]byte{Init: kvs}
	}

	newRevision := resp.Header.Revision + 1

	go func() {
		defer close(retChan)
		if initWe != nil {
			retChan <- initWe
		}
		for {
			watchOps := append(ops, clientv3.WithRev(newRevision))
			rch := c.client.Watch(ctx, key, watchOps...)

		LOOP:
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("[etcdx] watch %s ctx.Done", key)
					return
				case <-c.closeChan:
					fmt.Printf("[etcdx] watch %s close", key)
					return
				case <-c.client.Ctx().Done():
					return
				case watchEvent := <-rch:
					err := watchEvent.Err()
					if err != nil {
						fmt.Errorf("[etcdx] watch %s response error: %s ", key, err.Error())
						break LOOP
					}
					if len(watchEvent.Events) > 0 {
						var (
							createKVs map[string][]byte
							modifyKVs map[string][]byte
							deleteKVs map[string][]byte
						)
						events := &WatchEvent{
							Revision: watchEvent.Header.Revision,
						}
						for _, ev := range watchEvent.Events {
							if ev.IsCreate() {
								if createKVs == nil {
									createKVs = map[string][]byte{}
								}
								createKVs[string(ev.Kv.Key)] = ev.Kv.Value
							} else if ev.IsModify() {
								if modifyKVs == nil {
									modifyKVs = map[string][]byte{}
								}
								modifyKVs[string(ev.Kv.Key)] = ev.Kv.Value
							} else if ev.Type == mvccpb.DELETE {
								if deleteKVs == nil {
									deleteKVs = map[string][]byte{}
								}
								deleteKVs[string(ev.Kv.Key)] = ev.Kv.Value
							} else {
								fmt.Errorf("[etcdx] no found watch type: %s %q", ev.Type, ev.Kv.Key)
							}
						}
						events.KV = map[WatchEventType]map[string][]byte{
							Create: createKVs,
							Modify: modifyKVs,
							Delete: deleteKVs,
						}
						select {
						case retChan <- events:
						default:
						}
					}
					newRevision = watchEvent.Header.Revision + 1

				}
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	return retChan, nil

}

// Close 停止
func (c *Client) Close(ctx context.Context) {
	c.once.Do(func() {
		close(c.closeChan)
		c.wg.Wait()
		c.client.Close()
	})
}

func (c *Client) keepAlive(key string, val string, ttl int64) (clientv3.LeaseID, error) {
	getResp, err := c.client.Get(context.Background(), key)
	if nil != err {
		return 0, err
	}

	if getResp.Count > 0 {
		return 0, nil
	}

	leaseResp, err := c.client.Grant(context.Background(), ttl)
	if err != nil {
		return 0, err
	}

	_, err = c.client.Put(context.Background(), key, val, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return 0, err
	}

	keepAliveChan, err := c.client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return 0, err
	}
	go func() {
		for range keepAliveChan {
		}
	}()
	return leaseResp.ID, nil
}

func (c *Client) revoke(id clientv3.LeaseID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.client.Revoke(ctx, id)

	return err
}
