package redisclix

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// goRedis 接入 github.com/go-redis/redis/v8
type goRedis struct {
	rdb *redis.Client
}

func newGoRedis(cfg Config) *goRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Server,
		Password:     cfg.Auth,
		DB:           cfg.Index,
		MaxRetries:   2,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  -1 * time.Second, // 不超时
		WriteTimeout: -1 * time.Second, // 不超时
	})
	return &goRedis{
		rdb: rdb,
	}
}

func (c *goRedis) Ping() error {
	cmd := c.rdb.Ping(context.Background())
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (c *goRedis) Close() error {
	return c.rdb.Close()
}

// Set 设置值
// @param value - 如果为proto.Message, 则自动序列化
func (c *goRedis) Set(key string, value interface{}, expireTime int64) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	cmd := c.rdb.Set(context.Background(), key, data, time.Duration(expireTime)*time.Second)
	return c.parseError(cmd.Err())
}

func (c *goRedis) Get(key string) (string, error) {
	return c.rdb.Get(context.Background(), key).Result()
}

func (c *goRedis) HisExist(key string, field string) (bool, error) {
	return c.rdb.HExists(context.Background(), key, field).Result()
}

func (c *goRedis) IsExist(key string) (int64, error) {
	return c.rdb.Exists(context.Background(), key).Result()
}

func (c *goRedis) Keys(prefix string) ([]string, error) {
	return c.rdb.Keys(context.Background(), prefix).Result()
}

func (c *goRedis) Del(keys ...string) error {
	cmd := c.rdb.Del(context.Background(), keys...)
	return c.parseError(cmd.Err())
}

/*MGet
    SET key1 "hello"
	SET key2 "world"
	MGET key1 key2 key3
		1) "Hello"
		2) "World"
		3) (nil)
*/
func (c *goRedis) MGet(keys ...string) ([]string, error) {
	cmd := c.rdb.MGet(context.Background(), keys...)
	return sliceCmdConvert(cmd)
}

/*MSet
SET key1 "hello"
SET key2 "world"
MSET key1 "foo" key2 "bar"
	1) "OK"
*/
func (c *goRedis) MSet(fields ArrayReq) error {
	cmd := c.rdb.MSet(context.Background(), fields...)
	return c.parseError(cmd.Err())
}

/*HGet
HSET key1 field1 "foo"
HGET key1 field1
	1) "foo"
*/
func (c *goRedis) HGet(key string, field string) (string, error) {
	return c.rdb.HGet(context.Background(), key, field).Result()
}

/*HSet
HSET key1 field1 "foo"
HGET key1 field1
	1) "OK"
*/
func (c *goRedis) HSet(key string, field string, value interface{}) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	cmd := c.rdb.HSet(context.Background(), key, field, data)
	return c.parseError(cmd.Err())
}

/*HMSet
HMSET key1 field1 "hello" field2 "world"
	1) "OK"
*/
func (c *goRedis) HMSet(key string, value ArrayReq) error {
	cmd := c.rdb.HMSet(context.Background(), key, value...)
	return c.parseError(cmd.Err())
}

/*HMGet
HMSET key1 field1 "hello" field2 "world"
HMGET key1 field1 field2 field3
	1) "Hello"
	2) "World"
	3) (nil)
*/
func (c *goRedis) HMGet(key string, fields ...string) ([]string, error) {
	cmd := c.rdb.HMGet(context.Background(), key, fields...)
	return sliceCmdConvert(cmd)
}

/*HGetAll
HMSET key1 field1 "hello" field2 "world"
HGETALL key1
	1) "field1"
	2) "hello"
	3) "field2"
	4) "world"
*/
func (c *goRedis) HGetAll(key string) (map[string]string, error) {
	cmd := c.rdb.HGetAll(context.Background(), key)
	return cmd.Result()
}

/*HDel
HMSET key1 field1 "hello" field2 "world"
HDEL key1 field1 field2
	1) 2
*/
func (c *goRedis) HDel(key string, fields ...string) error {
	cmd := c.rdb.HDel(context.Background(), key, fields...)
	return c.parseError(cmd.Err())
}

/*HVals
HMSET key1 field1 "hello" field2 "world"
HVALS key1
	1) "hello"
	2) "world"
*/
func (c *goRedis) HVals(ctx context.Context, key string) ([]string, error) {
	cmd := c.rdb.HVals(ctx, key)
	return cmd.Val(), c.parseError(cmd.Err())
}

/*ZAdd
ZADD key1 1 "one"
ZADD key1 2 "two"
ZADD key1 3 "three"
ZRANGE key1 0 -1 WITHSCORES
	1) "one"
	2) "1"
	3) "two"
	4) "2"
	5) "three"
	6) "3"
*/
func (c *goRedis) ZAdd(key string, score ZAddReq) (int64, error) {
	if score.Count() == 0 {
		return 0, fmt.Errorf("score is empty")
	}
	var memebers = make([]*redis.Z, 0, len(score))
	for score, member := range score {
		memebers = append(memebers, &redis.Z{
			Score:  score,
			Member: member,
		})
	}
	cmd := c.rdb.ZAdd(context.Background(), key, memebers...)
	return cmd.Val(), c.parseError(cmd.Err())
}

func (c *goRedis) Do(args ...interface{}) (string, error) {
	ret := c.rdb.Do(context.Background(), args...)
	if ret.Err() != nil {
		return "", ret.Err()
	}
	return ret.String(), nil
}

/*ZRange
ZADD key1 1 "one"
ZADD key1 2 "two"
ZADD key1 3 "three"
ZRANGE key1 0 -1 WITHSCORES
	1) "one"
	2) "1"
	3) "two"
	4) "2"
	5) "three"
	6) "3"
*/
func (c *goRedis) ZRange(key string, start, stop int64) ([]string, error) {
	cmd := c.rdb.ZRange(context.Background(), key, start, stop)
	return cmd.Val(), c.parseError(cmd.Err())
}

func (c *goRedis) RPush(key string, value ...interface{}) error {
	return c.rdb.RPush(context.Background(), key, value).Err()
}

func (c *goRedis) LRange(key string, start, end int64) ([]string, error) {
	return c.rdb.LRange(context.Background(), key, start, end).Result()
}

func (c *goRedis) ListAll(key string) ([]string, error) {
	return c.LRange(key, 0, -1)
}

func sliceCmdConvert(cmd *redis.SliceCmd) ([]string, error) {
	var vals = make([]string, 0, len(cmd.Val()))
	for _, v := range cmd.Val() {
		switch val := v.(type) {
		case string:
			vals = append(vals, val)
		case nil:
			vals = append(vals, "")
		default:
			return nil, fmt.Errorf("unsupported type %T", val)
		}
	}
	return vals, nil
}

func (c *goRedis) parseError(err error) error {
	if err == redis.Nil {
		return nil
	}
	return err
}
