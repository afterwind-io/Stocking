package stocking

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents the... um... websocket client
type Client struct {
	connection *websocket.Conn
	id         string
	send       chan []byte
}

func (c *Client) Read(hub chan HubPackge) {
	defer func() {
		c.connection.Close()
	}()

	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		hub <- HubPackge{c, message}
	}
}

func (c *Client) Write() {
	defer func() {
		c.connection.Close()
	}()

	for {
		content := <-c.send

		w, err := c.connection.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		w.Write(content)
	}
}

func newClient(c *websocket.Conn) *Client {
	i := <-id

	client := Client{
		connection: c,
		id:         i,
		send:       make(chan []byte),
	}

	return &client
}
