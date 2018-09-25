package stocking

import (
	"encoding/json"
	"reflect"
	"sync"
)

type routeMeta struct {
	handler  RouterHandler
	typeHint interface{}
	mux      sync.Mutex
}

type mRouter struct {
	routes    map[string]routeMeta
	otherwise RouterHandler
}

func (me *mRouter) On(route string, handler RouterHandler, typeHint interface{}) {
	me.routes[route] = routeMeta{
		handler:  handler,
		typeHint: typeHint,
		mux:      sync.Mutex{},
	}
}

func (me *mRouter) Otherwise(handler RouterHandler) {
	me.otherwise = handler
}

func (me *mRouter) Handle(p *HubPackge, next MiddlewareStepFunc) {
	pkg, err := unmarshal(p.content)
	if err != nil {
		<-next(err)
		return
	}

	err = me.unserialize(&pkg)
	if err != nil {
		<-next(err)
		return
	}

	res, e := me.distribute(pkg)

	done := <-next(nil)

	mail, err := marshal(res, e)
	if err != nil {
		done <- err
		return
	}

	p.mailbox = mail

	done <- nil
}

func (me *mRouter) unserialize(p *RouterPackage) error {
	var body interface{}

	meta, ok := me.routes[p.Route]
	if ok && meta.typeHint != nil {
		body = clone(meta.typeHint)
	} else {
		b := make(map[string]interface{})
		body = &b
	}

	err := json.Unmarshal(p.Body.(json.RawMessage), body)
	if err != nil {
		return err
	}

	p.Body = body

	return nil
}

func (me *mRouter) distribute(p RouterPackage) (interface{}, error) {
	var res interface{}
	var e error

	meta, ok := me.routes[p.Route]
	if ok {
		res, e = meta.handler(p)
	} else {
		res, e = me.otherwise(p)
	}

	return res, e
}

func newRouter() *mRouter {
	return &mRouter{
		routes:    make(map[string]routeMeta),
		otherwise: blackhole,
	}
}

func unmarshal(raw []byte) (RouterPackage, error) {
	if !json.Valid(raw) {
		return RouterPackage{}, JSONSyntaxError{"Invalid Json Systax."}
	}

	var msg TextMessageInboundProtocol
	if err := json.Unmarshal(raw, &msg); err != nil {
		return RouterPackage{}, err
	}

	return RouterPackage{
		Route: msg.Route,
		Body:  msg.Body,
	}, nil
}

func marshal(body interface{}, e error) ([]byte, error) {
	errorMsg := ""
	if e != nil {
		errorMsg = e.Error()
	}

	msg := TextMessageOutboundProtocol{
		Error: errorMsg,
		Body:  body,
	}

	res, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func blackhole(p RouterPackage) (interface{}, error) {
	return struct{}{}, nil
}

// clone 接收一个预期类型为T的参数，返回一个新的*T值
func clone(source interface{}) interface{} {
	return reflect.New(reflect.ValueOf(source).Type()).Interface()
}
