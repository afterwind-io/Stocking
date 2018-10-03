package stocking

import (
	"log"
)

// Hub TODO
type Hub struct {
	signin      chan *Client
	signout     chan *Client
	inbound     chan *HubPackge
	clients     map[string]*Client
	channels    map[string]map[string]bool
	middlewares []Middleware
}

func (h *Hub) use(middlewares ...Middleware) {
	h.middlewares = append(h.middlewares, middlewares...)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.signin:
			h.onSignin(client)

		case client := <-h.signout:
			h.onSignout(client)

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
		hub.onJoin(p.ccode, p.content, p.client)
	}
}

func (h *Hub) onSignin(c *Client) {
	h.clients[c.id] = c

	// Join a channel named by client id,
	// so we can emit messages to the specified client
	// by broadcasting to this channel
	hub.onJoin("1", c.id, c)

	log.Printf("Connected: %v", c.id)
	c.send <- []byte(tmitConnect)
}

func (h *Hub) onSignout(c *Client) {
	client, ok := h.clients[c.id]
	if !ok {
		return
	}

	for channel := range client.channels {
		h.onJoin("0", channel, c)
	}
	h.onJoin("0", c.id, c)

	delete(h.clients, c.id)

	c.conn.Close()

	// Stop the Read/Write goroute
	c.oops <- true
}

func (h *Hub) onConnect(p *HubPackge) {
}

func (h *Hub) onClose(p *HubPackge) {
	h.onSignout(p.client)
}

func (h *Hub) onPingPong(p *HubPackge) {
	p.client.send <- []byte(tmitPingPong)
}

func (h *Hub) onMessage(p *HubPackge) {
	if err := flow(h.middlewares, p); err != nil {
		h.onSignout(p.client)
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

	for id := range group {
		client, ok := h.clients[id]

		if !ok || (p.client != nil && id == p.client.id) {
			continue
		}

		client.send <- p.raw
	}
}

func (h *Hub) onJoin(ccode string, channel string, c *Client) {
	if ccode == "1" {
		if _, ok := h.channels[channel]; !ok {
			h.channels[channel] = make(map[string]bool)
		}

		h.channels[channel][c.id] = true
		c.channels[channel] = true
	} else if ccode == "0" {
		if _, ok := c.channels[channel]; ok {
			delete(c.channels, channel)
			delete(h.channels[channel], c.id)
		}
	}
}

func newHub() *Hub {
	return &Hub{
		signin:      make(chan *Client),
		signout:     make(chan *Client),
		inbound:     make(chan *HubPackge, 100),
		clients:     make(map[string]*Client),
		channels:    make(map[string]map[string]bool),
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
