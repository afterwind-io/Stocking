package stocking

import (
	"encoding/json"
)

type mRouter struct {
	routes    map[string]RouterHandler
	otherwise RouterHandler
}

func (me *mRouter) Handle(p *HubPackge, next MiddlewareStepFunc) {
	if !json.Valid(p.content) {
		<-next(JSONSyntaxError{"Invalid Json Systax."})
		return
	}

	var msg TextMessageProtocol
	if err := json.Unmarshal(p.content, &msg); err != nil {
		<-next(err)
		return
	}

	pkg := RouterPackage{
		route: msg.r,
		body:  msg.b,
	}
	handler, ok := me.routes[msg.r]
	if ok {
		handler(pkg)
	} else {
		me.otherwise(pkg)
	}

	<-next(nil)

	// TODO
}

func (me *mRouter) On(route string, handler RouterHandler) {
	me.routes[route] = handler
}

func (me *mRouter) Otherwise(handler RouterHandler) {
	me.otherwise = handler
}

func newRouter() *mRouter {
	return &mRouter{
		routes:    make(map[string]RouterHandler),
		otherwise: func(p RouterPackage) {},
	}
}
