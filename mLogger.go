package stocking

import (
	"fmt"
	"log"
)

type mLogger struct {
}

func (me *mLogger) Handle(p *HubPackge, next MiddlewareStepFunc) {
	log.Println(fmt.Sprintf("<-- [%v] %v, %v, %v", p.client.id, p.mtype, p.ack, p.content))

	done := <-next(nil)

	log.Println(fmt.Sprintf("--> [%v] %v, %v, %v", p.client.id, p.mtype, p.ack, p.content))

	done <- nil
}
