package stocking

import (
	"log"
)

// Hub TODO
type Hub struct {
	registry    chan *Client
	inbound     chan *HubPackge
	clients     map[string]*Client
	channels    map[string][]string
	middlewares []Middleware
}

func (h *Hub) use(middlewares ...Middleware) {
	h.middlewares = append(h.middlewares, middlewares...)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.registry:
			h.onRegister(client)

		case p := <-h.inbound:
			h.switchType(p)
		}
	}
}

func (h *Hub) broadcast(target string, content []byte) {
	hub.inbound <- &HubPackge{
		client:  nil,
		raw:     content,
		mtype:   tmitBroadcast,
		ccode:   target,
		content: string(content),
	}
}

func (h *Hub) switchType(p *HubPackge) {
	switch p.mtype {
	case tmitConnect:
		hub.onConnect(p)

	case tmitClose:
		hub.onClose(p)

	case tmitPingPong:
		hub.onPingPong(p)

	case tmitMessage:
		hub.onMessage(p)

	case tmitBroadcast:
		hub.onBroadcast(p)

	case tmitJoin:
		hub.onJoin(p)
	}
}

func (h *Hub) onRegister(c *Client) {
	h.clients[c.id] = c

	p := &HubPackge{
		client:  c,
		ccode:   "1",
		content: c.id,
	}
	hub.onJoin(p)

	log.Printf("Connected: %v", c.id)
	c.send <- []byte(tmitConnect)
}

func (h *Hub) onConnect(p *HubPackge) {
}

func (h *Hub) onClose(p *HubPackge) {
	// TODO
	p.client.connection.Close()
}

func (h *Hub) onPingPong(p *HubPackge) {
	p.client.send <- []byte(tmitPingPong)
}

func (h *Hub) onMessage(p *HubPackge) {
	if err := flow(h.middlewares, p); err != nil {
		// TODO
		delete(h.clients, p.client.id)
		p.client.connection.Close()
	} else if p.hasAck() {
		p.client.send <- p.encode()
	}
}

func (h *Hub) onBroadcast(p *HubPackge) {
	group, ok := h.channels[p.ccode]
	if !ok {
		p.error("0", `Channel "`+p.ccode+`" not found`)
		p.client.send <- p.encode()
		return
	}

	for _, id := range group {
		client, ok := h.clients[id]

		if !ok || (p.client != nil && id == p.client.id) {
			continue
		}

		client.send <- p.raw
	}
}

func (h *Hub) onJoin(p *HubPackge) {
	if p.ccode == "1" {
		group, ok := h.channels[p.content]
		if ok {
			h.channels[p.content] = append(group, p.client.id)
		} else {
			h.channels[p.content] = []string{p.client.id}
		}

		p.client.channel[p.content] = true
	} else if p.ccode == "0" {
		if _, ok := p.client.channel[p.content]; !ok {
			p.error("0", "Nope")
		} else {
			delete(p.client.channel, p.content)
			delete(h.channels, p.content)
		}
	}
}

func newHub() *Hub {
	return &Hub{
		registry:    make(chan *Client),
		inbound:     make(chan *HubPackge, 100),
		clients:     make(map[string]*Client),
		channels:    make(map[string][]string),
		middlewares: []Middleware{},
	}
}

func flow(middlewares []Middleware, p *HubPackge) error {
	if len(middlewares) == 0 {
		return nil
	}

	middleware := middlewares[0]
	forward, backward, next, done := step()

	go middleware.Handle(p, next)

	if err := <-forward; err != nil {
		log.Println("middleware error: ", err)
		return err
	}

	if err := flow(middlewares[1:], p); err != nil {
		return err
	}

	backward <- done

	if err := <-done; err != nil {
		return err
	}

	return nil
}

func step() (chan error, chan chan error, MiddlewareStepFunc, chan error) {
	forward := make(chan error)
	backward := make(chan chan error)
	next := func(err error) chan chan error {
		forward <- err
		return backward
	}
	done := make(chan error)

	return forward, backward, next, done
}
