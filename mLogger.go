package stocking

import (
	"fmt"
	"log"
)

type mLogger struct {
}

func (me *mLogger) Forward(p *HubPackge) error {
	log.Println(fmt.Sprintf("<-- [%v] %s", p.client.id, p.content))
	return nil
}

func (me *mLogger) Backward(p *HubPackge) error {
	log.Println(fmt.Sprintf("--> [%v] %s", p.client.id, p.content))
	return nil
}

