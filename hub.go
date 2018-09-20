package stocking

import (
	"log"
)

// Hub TODO
type Hub struct {
	registry    chan *Client
	inbound     chan HubPackge
	clients     map[*Client]bool
	middlewares []Middleware
}

// HubPackge TODO
type HubPackge struct {
	client  *Client
	content []byte
	mailbox []byte
}

func (hub *Hub) use(middlewares ...Middleware) {
	hub.middlewares = append(hub.middlewares, middlewares...)
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.registry:
			hub.clients[client] = true
			log.Printf("Connected: %v", client.id)

		case p := <-hub.inbound:
			if err := flow(hub.middlewares, &p); err != nil {
				// TODO
				delete(hub.clients, p.client)
				p.client.connection.Close()
			}

			if len(p.mailbox) != 0 {
				p.client.send <- p.mailbox
			}
		}
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

func newHub() *Hub {
	return &Hub{
		registry:    make(chan *Client),
		inbound:     make(chan HubPackge, 100),
		clients:     make(map[*Client]bool),
		middlewares: []Middleware{},
	}
}

func newHubPackage(c *Client, m []byte) HubPackge {
	return HubPackge{
		client:  c,
		content: m,
		mailbox: []byte{},
	}
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
