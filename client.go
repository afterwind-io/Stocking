package stocking

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	readTimeout      = 120 * time.Second
	writeTimeout     = 10 * time.Second
	heartbeatTimeout = 60 * time.Second
)

// Client represents the... um... websocket client
type Client struct {
	// conn refers to the underlying websocket connection
	conn *websocket.Conn
	// hub refers to the attached hub
	hub *Hub
	// id is the inner index of the client
	id string
	// channels keeps every channel currently subscribed
	channels map[string]bool
	// lasthb records the timestamp of last ping/inbound-message
	lasthb time.Time
	// send is the channel for sending message
	send chan []byte
	// oops is the channel for breaking the Read/Write goroute loop
	oops chan bool
}

// Close the connection
func (c *Client) Close() {
	p := &HubPackge{
		mtype: tmitClose,
	}
	c.hub.inbound <- p
}

func (c *Client) read() {
	defer func() {
		c.hub.signout <- c
	}()

	for {
		c.conn.SetReadDeadline(time.Now().Add(readTimeout))

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Accept any inbound message as a ping,
		// because the client is safe and sound after all
		c.lasthb = time.Now()

		p, err := newHubPackage(c, message)
		if err != nil {
			log.Println(err.Error())
			break
		}

		c.hub.inbound <- p
	}
}

func (c *Client) write() {
	heartbeat := time.NewTicker(heartbeatTimeout)

	defer func() {
		heartbeat.Stop()

		c.hub.signout <- c
	}()

	for {
		select {
		case <-c.oops:
			return

		case <-heartbeat.C:
			if time.Since(c.lasthb) > heartbeatTimeout {
				log.Printf("Id %v: Connection Timeout", c.id)
				return
			}

		case content := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))

			if err := c.conn.WriteMessage(websocket.TextMessage, content); err != nil {
				log.Printf("Id %v: Write Timeout", c.id)
				return
			}
		}
	}
}

func newClient(c *websocket.Conn, h *Hub) *Client {
	i := <-id

	client := Client{
		id:       i,
		conn:     c,
		hub:      h,
		channels: make(map[string]bool),
		lasthb:   time.Now(),
		send:     make(chan []byte, 1),
		oops:     make(chan bool, 1),
	}

	return &client
}
