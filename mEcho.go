package stocking

type mEcho struct {
}

func (me *mEcho) Forward(p *HubPackge) error {
	return nil
}

func (me *mEcho) Backward(p *HubPackge) error {
	hello := []byte("Hello ")
	p.content = append(hello, p.content...)

	return nil
}
