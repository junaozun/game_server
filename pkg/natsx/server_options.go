package natsx

import (
	"fmt"
)

// ServiceOption service option
type ServiceOption func(options *serviceOptions)

// serviceOptions service 选项
type serviceOptions struct {
	id             string            // id
	errorHandler   func(interface{}) // error handler
	recoverHandler func(interface{}) // recover handler
}

// WithServiceID id
func WithServiceID(id interface{}) ServiceOption {
	return func(options *serviceOptions) {
		options.id = fmt.Sprintf("%v", id)
	}
}
