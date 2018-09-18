package stocking

type mRouter struct {
	index map[string]RouterHandler
}

func (me *mRouter) Forward(p *HubPackge) error {
	// msg := fmt.Sprint(p.content)
	return nil
}

func (me *mRouter) Backward(p *HubPackge) error {
	return nil
}

func (me *mRouter) On(route string, handler RouterHandler) {
	me.index[route] = handler
}

func (me *mRouter) Otherwise(handler RouterHandler) {

}

func newRouter() *mRouter {
	return &mRouter{
		index: make(map[string]RouterHandler),
	}
}
