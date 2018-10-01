package stocking

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents the... um... websocket client
type Client struct {
	connection *websocket.Conn
	id         string
	channel    map[string]bool
	send       chan []byte
}

func (c *Client) Read(hub chan *HubPackge) {
	defer func() {
		c.connection.Close()
	}()

	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}

		p, err := newHubPackage(c, message)
		if err != nil {
			break
		}

		hub <- p
	}
}

func (c *Client) Write() {
	defer func() {
		c.connection.Close()
	}()

	for {
		content := <-c.send

		if err := c.connection.WriteMessage(
			websocket.TextMessage,
			content,
		); err != nil {
			return
		}
	}
}

func newClient(c *websocket.Conn) *Client {
	i := <-id

	client := Client{
		connection: c,
		id:         i,
		channel:    make(map[string]bool),
		send:       make(chan []byte, 1),
	}

	return &client
}
