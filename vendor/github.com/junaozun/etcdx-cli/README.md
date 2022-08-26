# etcdx-cli
``对etcd client的包装使用``

# 目前拥有的方法

````
        // Put 放入键值对
        Put(ctx context.Context, key, value string) error

	// PutJSON 放入json对象
	PutJSON(ctx context.Context, key string, data interface{}) error

	// PutWithTTL 放入带存活时间的键值
	PutWithTTL(ctx context.Context, key, value string, ttl int64) error

	// PutJSONWithTTL 放入带存活时间的json对象
	PutJSONWithTTL(ctx context.Context, key string, data interface{}, ttl int64) error

	// Get 获取值
	Get(ctx context.Context, key string) ([]byte, error)

	// GetString 获取字符串值
	GetString(ctx context.Context, key string) (string, error)

	// GetJSON 获取json对象
	GetJSON(ctx context.Context, key string, data interface{}) error

	// GetAllKey 获取前缀key匹配key所有键
	GetAllKey(ctx context.Context, key string) ([]string, error)

	// GetAllValue 获取前缀key匹配key所有值
	GetAllValue(ctx context.Context, key string) ([][]byte, error)

	// Delete 删除键值
	Delete(ctx context.Context, key string) error

	// KV 获取前缀key匹配的所有键值对
	KV(ctx context.Context, key string) (map[string][]byte, error)

	// PutIfNotExist 只有当不存在时才放入成功
	PutIfNotExist(ctx context.Context, create KV, kv ...KV) (bool, error)

	// UpdateTx 更新事务
	UpdateTx(ctx context.Context, key, val string) error

	// KeepAlive 保活
	KeepAlive(ctx context.Context, key string, val string, ttl int64) (err error)

	// Watch 监听
	Watch(ctx context.Context, key string, prefix bool) (<-chan *WatchEvent, error)

	// Close 关闭
	Close(ctx context.Context)``
