package stocking

import (
	"encoding/json"
)

type mRouter struct {
	routes    map[string]RouterHandler
	otherwise RouterHandler
}

func (me *mRouter) On(route string, handler RouterHandler) {
	me.routes[route] = handler
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

func (me *mRouter) distribute(p RouterPackage) (interface{}, error) {
	var res interface{}
	var e error

	handler, ok := me.routes[p.Route]
	if ok {
		res, e = handler(p)
	} else {
		res, e = me.otherwise(p)
	}

	return res, e
}

func newRouter() *mRouter {
	return &mRouter{
		routes:    make(map[string]RouterHandler),
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
