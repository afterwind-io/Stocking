package stocking

// Middleware TODO
type Middleware interface {
	Handle(p *HubPackge, next MiddlewareStepFunc)
}

// MiddlewareStepFunc TODO
type MiddlewareStepFunc = func(err error) chan chan error
