package redisclix

import (
	"context"
	"fmt"
	"hash/crc32"
	"net"
	"time"
)

// 选项
type options struct {
	concurrency    int
	hash           hashFun
	maxCmdQueue    int // 每个thread最大命令队列
	maxReturnQueue int // 最大返回队列
	pushFunc       PushFunc
	errorHandler   ErrorHandler
	onStop         []func()
	overtimeDur    time.Duration // 超时时间
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func defaultOption() options {
	return options{
		maxCmdQueue:    1024,
		maxReturnQueue: 1024,
		concurrency:    8,
		overtimeDur:    100 * time.Millisecond,
		hash: func(key string, count int) int {
			if count == 1 {
				return 0
			}
			return int(crc32.ChecksumIEEE([]byte(key)) % uint32(count))
		},
		errorHandler: func(err error) {
			if _, ok := err.(net.Error); ok {
				fmt.Errorf("[AsyncClient] net err[%v]", err)
			}
		},
	}
}

type (
	hashFun      func(string, int) int
	ErrorHandler func(error) // 错误处理
	PushFunc     func(ctx context.Context, f func()) error
)

// WithMaxCmdQueue 设置最大命令队列
func WithMaxCmdQueue(size int) Option {
	return optionFunc(func(opt *options) {
		opt.maxCmdQueue = size
	})
}

// WithMaxReturnQueue 设置最大返回队列
func WithMaxReturnQueue(size int) Option {
	return optionFunc(func(opt *options) {
		opt.maxReturnQueue = size
	})
}

// WithPushFunc 设置func
func WithPushFunc(f func(ctx context.Context, f func()) error) Option {
	return optionFunc(func(opt *options) {
		opt.pushFunc = f
	})
}

func WithStopHandler(f func()) Option {
	return optionFunc(func(opt *options) {
		opt.onStop = append(opt.onStop, f)
	})
}

// WithConcurrency 设置最大并发数
func WithConcurrency(n int) Option {
	return optionFunc(func(opt *options) {
		opt.concurrency = n
	})
}

// WithHashFunc 设置最大返回队列
func WithHashFunc(hash hashFun) Option {
	return optionFunc(func(opt *options) {
		opt.hash = hash
	})
}
