package natsx

import (
	"context"
	"fmt"
	"go/ast"
	"os"
	"reflect"
	"sync"

	"github.com/nats-io/nats.go"
)

type Server struct {
	wg          sync.WaitGroup
	mutex       sync.Mutex
	connEnc     *nats.EncodedConn
	services    map[string]*service             // serverName(chess、logic、gvg):service
	opt         serviceOptions                  // options
	servicesSub map[string][]*nats.Subscription // serviceName:所有订阅
}

func NewServer(connEnc *nats.EncodedConn) (*Server, error) {
	if !connEnc.Conn.IsConnected() {
		return nil, fmt.Errorf("enc is not connected")
	}
	opt := serviceOptions{
		errorHandler: func(i interface{}) {
			fmt.Fprintf(os.Stderr, "error:%v\n", i)
		},
		recoverHandler: func(i interface{}) {
			fmt.Fprintf(os.Stderr, "server panic:%v\n", i)
		},
	}
	s := &Server{
		opt:         opt,
		connEnc:     connEnc,
		services:    make(map[string]*service),
		servicesSub: make(map[string][]*nats.Subscription),
	}
	return s, nil
}

func (s *Server) Register(serverName string, svc interface{}, opts ...ServiceOption) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range opts {
		v(&s.opt)
	}
	serverName = CombineSubject(serverName, s.opt.id)
	service := newservice(serverName)
	if v, ok := s.services[serverName]; ok {
		service = v
	} else {
		s.services[serverName] = service
	}
	err := service.parseServices(svc)
	if err != nil {
		panic(err)
	}

	err = s.subscribeMethod()
	if err != nil {
		panic(err)
	}
}

// 订阅所有的方法
func (s *Server) subscribeMethod() error {
	for _, service := range s.services {
		for _, obj := range service.objects {
			for subject, v := range obj.methods {
				method := v
				cb := func(msg *nats.Msg) {
					s.wg.Add(1)
					go func() {
						defer s.wg.Done()
						err := s.DealSubFunction(context.Background(), obj, msg, method)
						if err != nil {
							s.opt.errorHandler(err.Error())
						}
					}()
				}
				natsSub, subErr := s.connEnc.Subscribe(subject, cb)
				if subErr != nil {
					return subErr
				}
				s.servicesSub[service.serverName] = append(s.servicesSub[service.serverName], natsSub)
			}
		}
	}
	return nil
}

func (s *Server) DealSubFunction(ctx context.Context, obj *object, msg *nats.Msg, method *method) error {
	if s.opt.recoverHandler != nil {
		defer func() {
			if e := recover(); e != nil {
				s.opt.recoverHandler(e)
			}
		}()
	}
	reply, err := s.handle(ctx, obj, msg, method)
	if len(reply) == 0 { // notify (reply == 0)
		return nil
	}
	if s.connEnc.Conn.IsClosed() {
		return fmt.Errorf("conn closed")
	}
	response := &Reply{
		Payload: reply,
	}
	if err != nil {
		response.Error = err.Error()
	}
	b, e := s.connEnc.Enc.Encode(msg.Subject, response)
	if e != nil {
		return fmt.Errorf("[encode err:%s]", e.Error())
	}
	respMsg := &nats.Msg{
		Subject: msg.Reply,
		Data:    b,
	}
	return s.connEnc.Conn.PublishMsg(respMsg)
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

func (s *Server) Close(ctx context.Context) error {
	s.CancelAllServerSub()
	if err := s.connEnc.Flush(); err != nil {
		return err
	}
	return nil
}

// CancelAllServerSub  取消所有服务订阅
func (s *Server) CancelAllServerSub() {
	s.mutex.Lock()
	srvNames := make([]string, 0, len(s.servicesSub))
	for k := range s.servicesSub {
		srvNames = append(srvNames, k)
	}
	s.mutex.Unlock()
	for _, v := range srvNames {
		s.cancelServerSubByName(v)
	}
}

// cancelServerSubByName 根据服务名取消该服务下的所有订阅
func (s *Server) cancelServerSubByName(srvName string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	subs, ok := s.servicesSub[srvName]
	if !ok {
		return false
	}
	for _, sub := range subs {
		sub.Unsubscribe()
		delete(s.servicesSub, srvName)
	}
	return true
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
	methods    map[string]*method // 方法集合 (服务名.服务子id.对象名.方法名: 方法)
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
		subject := CombineSubject(o.serverName, typeName, funName)
		o.methods[subject] = funcVal
	}
	return nil
}
