package natsx

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
)

// NewNatsConnEnc 创建enc
func NewNatsConnEnc(url string, encType string, option ...nats.Option) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(url, option...)
	if err != nil {
		return nil, err
	}
	enc, err1 := nats.NewEncodedConn(nc, encType)
	if nil != err1 {
		return nil, err1
	}
	return enc, nil
}

// NewNatsPBEnc 创建 protobuf enc
func NewNatsPBEnc(url string, option ...nats.Option) (*nats.EncodedConn, error) {
	return NewNatsConnEnc(url, protobuf.PROTOBUF_ENCODER, option...)
}

// NewNatsJSONEnc 创建 json enc
func NewNatsJSONEnc(url string, option ...nats.Option) (*nats.EncodedConn, error) {
	return NewNatsConnEnc(url, nats.JSON_ENCODER, option...)
}

// NewNatsGobEnc 创建 json enc
func NewNatsGobEnc(url string, option ...nats.Option) (*nats.EncodedConn, error) {
	return NewNatsConnEnc(url, nats.GOB_ENCODER, option...)
}
