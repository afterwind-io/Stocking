package stocking

import (
	"fmt"
	"log"
)

type mLogger struct {
}

func (me *mLogger) Handle(p *HubPackge, next MiddlewareStepFunc) {
	log.Println(fmt.Sprintf("<-- [%v] %s", p.client.id, p.content))

	<-next(nil)

	log.Println(fmt.Sprintf("--> [%v] %s", p.client.id, p.content))
}
