package natsx

import (
	"context"
	"fmt"
	"go/ast"
	"reflect"
	"sync"

	"github.com/nats-io/nats.go"
)

type Server struct {
	wg       sync.WaitGroup
	mutex    sync.Mutex
	connEnc  *nats.EncodedConn
	services map[string]*service             // serverName:service (chess:service)(logic:service)
	sersub   map[string][]*nats.Subscription // subject:所有订阅
}

func NewServer(connEnc *nats.EncodedConn) (*Server, error) {
	if !connEnc.Conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}
	s := &Server{
		connEnc:  connEnc,
		services: make(map[string]*service),
		sersub:   make(map[string][]*nats.Subscription),
	}
	return s, nil
}

func (s *Server) Register(serverName string, svc interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	service := newservice(serverName)
	if v, ok := s.services[serverName]; ok {
		service = v
	} else {
		s.services[serverName] = service
	}
	err := service.parseServices(svc)
	if err != nil {
		return err
	}

	err = s.subscribeMethod()
	if err != nil {
		return err
	}
	return nil
}

// 订阅所有的方法
func (s *Server) subscribeMethod() error {
	for _, service := range s.services {
		for _, obj := range service.objects {
			for subject, v := range obj.methods {
				method := v
				natsSub, subErr := s.connEnc.Subscribe(subject, func(msg *nats.Msg) {
					s.wg.Add(1)
					go func() {
						defer s.wg.Done()
						reply, err := s.handle(context.Background(), obj, msg, method)
						if len(reply) == 0 { // notify (reply == 0)
							return
						}
						if s.connEnc.Conn.IsClosed() {
							fmt.Errorf("conn closed")
							return
						}
						response := Reply{
							Payload: reply,
						}
						if err != nil {
							response.Error = err.Error()
						}
						b, e := s.connEnc.Enc.Encode(msg.Subject, response)
						if e != nil {
							fmt.Errorf("[encode err:%s]", e.Error())
							return
						}
						respMsg := &nats.Msg{
							Subject: msg.Reply,
							Data:    b,
						}
						err = s.connEnc.Conn.PublishMsg(respMsg)
						if err != nil {
							fmt.Errorf("[publishMsg err:%s]", err.Error())
						}
					}()
				})
				if subErr != nil {
					return subErr
				}
				// todo
				s.sersub[subject] = append(s.sersub[subject], natsSub)
			}
		}
	}
	return nil
}

func (s *Server) handle(ctx context.Context, obj *object, msg *nats.Msg, method *method) ([]byte, error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Errorf("[Server] handle ,err:%s", e)
		}
	}()
	req := method.newRequest()
	b := msg.Data
	if len(b) > 0 {
		rpcReq := &Request{}
		if err := s.connEnc.Enc.Decode(msg.Subject, b, rpcReq); err != nil {
			return nil, err
		}
		if len(rpcReq.Header) > 0 {
			// todo
		}
		if len(rpcReq.Payload) > 0 {
			if err := s.connEnc.Enc.Decode(msg.Subject, rpcReq.Payload, req); nil != err {
				return nil, err
			}
		}
	}
	resp, err := method.handle(obj.srv, ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return s.connEnc.Enc.Encode(msg.Subject, resp)
}

type service struct {
	serverName string
	objects    map[interface{}]*object // 对象名:对象
}

func newservice(serverName string) *service {
	return &service{
		serverName: serverName,
		objects:    make(map[interface{}]*object),
	}
}

func (s *service) parseServices(srv interface{}) error {
	if _, ok := s.objects[srv]; ok {
		return fmt.Errorf("srv :%v exist", srv)
	}
	object := newObject(s.serverName, srv)
	s.objects[srv] = object
	return object.parseObject()
}

type object struct {
	serverName string
	srv        interface{}        // 对象名
	methods    map[string]*method // 方法集合 (服务名.对象名.方法名: 方法)
}

func newObject(serverName string, srv interface{}) *object {
	return &object{
		serverName: serverName,
		srv:        srv,
		methods:    make(map[string]*method),
	}
}

func (o *object) parseObject() error {
	val := reflect.ValueOf(o.srv)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("service must be a pointer")
	}
	typ := reflect.Indirect(val).Type()
	// 对象名
	typeName := typ.Name()
	if !ast.IsExported(typeName) {
		return fmt.Errorf("service [%s] typeName:%s must be exported", o.serverName, typeName)
	}

	methods, err := parseMethod(o.srv)
	if err != nil {
		return err
	}
	if len(methods) == 0 {
		return fmt.Errorf("service [%s] type:[%s] has no exported method", o.serverName, typeName)
	}
	for funName, funcVal := range methods {
		subject := o.serverName + "." + typeName + "." + funName
		o.methods[subject] = funcVal
	}
	return nil
}
