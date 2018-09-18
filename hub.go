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

			if len(p.content) != 0 {
				p.client.send <- p.content
			}
		}
	}
}

func flow(middlewares []Middleware, p *HubPackge) error {
	if len(middlewares) == 0 {
		return nil
	}

	middleware := middlewares[0]

	if err := middleware.Forward(p); err != nil {
		log.Println("middleware error: ", err)
		return err
	}

	if err := flow(middlewares[1:], p); err != nil {
		return err
	}

	if err := middleware.Backward(p); err != nil {
		log.Println("middleware error: ", err)
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