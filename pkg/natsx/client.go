package natsx

import (
	"context"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

type Client struct {
	connEnc     *nats.EncodedConn
	serviceName string
}

func NewClient(enc *nats.EncodedConn, serviceName string) (*Client, error) {
	return &Client{
		connEnc:     enc,
		serviceName: serviceName,
	}, nil
}

func (c *Client) Publish(objectName string, methodName string, req interface{}, opt ...CallOption) error {
	return c.call(nil, objectName, methodName, req, nil, opt...)
}

func (c *Client) Request(ctx context.Context, objectName string, methodName string, req interface{}, resp interface{}, opt ...CallOption) error {
	return c.call(ctx, objectName, methodName, req, resp, opt...)
}

func (c *Client) call(ctx context.Context, objectName string, methodName string, req interface{}, resp interface{}, opt ...CallOption) error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintf(os.Stderr, "client call panic:%v\n", e)
		}
	}()
	callOpt := CallOptions{}
	for _, v := range opt {
		v(&callOpt)
	}
	var isPublish bool
	if ctx == nil && resp == nil {
		isPublish = true
	}
	subject := CombineSubject(c.serviceName, callOpt.id, objectName, methodName)
	rpcReq, err := c.newRequest(subject, req, nil)
	if err != nil {
		return err
	}
	if isPublish {
		return c.connEnc.Publish(subject, rpcReq)
	}
	response := &Reply{}
	err = c.connEnc.RequestWithContext(ctx, subject, rpcReq, response)
	if err != nil {
		return err
	}
	if len(response.Error) > 0 {
		return fmt.Errorf(response.Error)
	}
	return c.connEnc.Enc.Decode(subject, response.Payload, resp)
}

func (c *Client) newRequest(subject string, req interface{}, header map[string]string) (*Request, error) {
	payload, err := c.connEnc.Enc.Encode(subject, req)
	if err != nil {
		return nil, err
	}
	return &Request{
		Payload: payload,
		Header:  header,
	}, nil
}
