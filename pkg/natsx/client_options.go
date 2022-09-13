package natsx

import (
	"fmt"
)

// CallOptions 调用选项
type CallOptions struct {
	id     string            //
	header map[string]string // header
}

// CallOption call option
type CallOption func(options *CallOptions)

// WithCallID call id
func WithCallID(id interface{}) CallOption {
	return func(options *CallOptions) {
		options.id = fmt.Sprint(id)
	}
}

// WithCallHeader header
func WithCallHeader(hd map[string]string) CallOption {
	return func(options *CallOptions) {
		options.header = hd
	}
}
