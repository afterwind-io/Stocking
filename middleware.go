package stocking

// Middleware TODO
type Middleware interface {
	Forward(p *HubPackge) error
	Backward(p *HubPackge) error
}
