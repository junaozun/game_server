package redisclix

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// AsyncClient redis 异步调用
type AsyncClient struct {
	cli      IClient
	clients  []*cmdClient
	out      chan Callbacker
	stopChan chan struct{}
	opt      options
	wg       sync.WaitGroup
	once     sync.Once
}

func NewAsyncWithConfig(cfg Config, opts ...Option) (*AsyncClient, error) {
	cli, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, WithStopHandler(func() {
		if err := cli.Close(); err != nil {
			fmt.Errorf("[NewAsyncWithConfig] err [%v]", err)
		}
	}))
	return NewAsync(cli, opts...), nil
}

func NewAsync(cli IClient, opts ...Option) *AsyncClient {
	var options = defaultOption()
	for _, o := range opts {
		o.apply(&options)
	}
	c := &AsyncClient{
		out:      make(chan Callbacker, options.maxReturnQueue),
		stopChan: make(chan struct{}),
		opt:      options,
		cli:      cli,
	}

	for i := 0; i < options.concurrency; i++ {
		client := c.newCmdClient(options.maxCmdQueue)
		c.clients = append(c.clients, client)
	}

	// 启动所有cmdClient
	c.wg.Add(1)
	go c.runAllClient()

	if options.pushFunc != nil {
		c.wg.Add(1)
		go c.pushCallback(options.pushFunc)
	}
	return c
}

type cmdClient struct {
	cmds chan Command // redis client 的命令队列
}

// Command redis 命令(无需回调)
type Command interface {
	ExecCmd(cli IClient) error // 执行命令
	Key() string               // key
}

// Callbacker 执行回掉
type Callbacker interface {
	Callback()
}

// CallbackCmder 带回调的redis命令接口
type CallbackCmder interface {
	Command
	Callbacker // 执行回掉
}

func (c *AsyncClient) newCmdClient(maxCmdQueue int) *cmdClient {
	return &cmdClient{
		cmds: make(chan Command, maxCmdQueue),
	}
}

func (c *AsyncClient) runAllClient() {
	defer c.wg.Done()
	defer close(c.out)

	var wg sync.WaitGroup
	for _, w := range c.clients {
		wg.Add(1)
		go func(cw *cmdClient) {
			defer wg.Done()
			c.serve(cw)
		}(w)
	}

	wg.Wait()
}

func (c *AsyncClient) serve(w *cmdClient) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Errorf("panic: %v: %v", err, buf)
		}
	}()

LOOP:
	for {
		select {
		case cmd := <-w.cmds:
			begin := time.Now()
			if err := cmd.ExecCmd(c.cli); err != nil {
				fmt.Errorf("AsyncClient: server cmd[%v] %s", cmd.Key(), err)
				if c.opt.errorHandler != nil {
					c.opt.errorHandler(err)
				}
			} else {
				executeTime := time.Since(begin)
				if executeTime > c.opt.overtimeDur {
					fmt.Errorf("[AsyncClient] cmd[%v] execute slow, time[%vms]", cmd.Key(), executeTime.Milliseconds())
				}
			}
			if cbCmd, ok := cmd.(CallbackCmder); ok {
				c.out <- cbCmd
			}
		case <-c.stopChan:
			break LOOP
		}
	}

	// 主动关闭后，还需要将队列中的命令全部执行完
	for {
		select {
		case cmd := <-w.cmds:
			if err := cmd.ExecCmd(c.cli); err != nil {
				fmt.Errorf("redisasync: exit server cmd[%v] %s", cmd.Key(), err)
			}
		default:
			return
		}
	}
}

// Sync 同步调用
func (c *AsyncClient) Sync() IClient {
	return c.cli
}

func (c *AsyncClient) Out() <-chan Callbacker {
	return c.out
}

func (c *AsyncClient) pushCallback(pushFunc PushFunc) {
	defer c.wg.Done()
	for cb := range c.Out() {
		if err := pushFunc(context.Background(), cb.Callback); err != nil {
			fmt.Errorf("[AsyncClient] pushCallback %v", err)
		}
	}
}

// Stop 关闭
func (c *AsyncClient) Stop(ctx context.Context) (err error) {
	c.once.Do(func() {
		close(c.stopChan)
		over := make(chan struct{})
		go func() {
			c.wg.Wait()
			close(over)
		}()
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-over:
			break
		}
	})
	return
}

func (c *AsyncClient) addCmd(cmd Command) {
	index := c.opt.hash(cmd.Key(), len(c.clients))
	c.clients[index].cmds <- cmd
}

func (c *AsyncClient) Set(key string, value interface{}) error {
	req, err := Encode(value)
	if err != nil {
		return err
	}
	c.addCmd(&setCmd{
		key:   key,
		value: req,
	})
	return nil
}

func (c *AsyncClient) Do(cb func(string, error), args ...interface{}) error {
	var encodeArgs = make([]interface{}, 0, len(args))
	for _, arg := range args {
		v, err := Encode(arg)
		if err != nil {
			return err
		}
		encodeArgs = append(encodeArgs, v)
	}
	c.addCmd(&doCmd{
		cmd:      encodeArgs,
		callback: cb,
	})
	return nil
}

func (c *AsyncClient) Get(key string, cb func(string, error)) {
	c.addCmd(&getCmd{
		key:      key,
		callback: cb,
	})
}

func (c *AsyncClient) Del(key string) {
	c.addCmd(&delCmd{
		key: key,
	})
}

func (c *AsyncClient) HMSet(key string, value ArrayReq) {
	c.addCmd(&hmsetCmd{
		key:   key,
		value: value,
	})
}

func (c *AsyncClient) HMGet(key string, field []string, cb func([]string, error)) {
	c.addCmd(&hmgetCmd{
		key:      key,
		field:    field,
		callback: cb,
	})
}
