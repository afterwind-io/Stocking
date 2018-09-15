package stocking

import "github.com/gorilla/websocket"

// Client represents the... um... websocket client
type Client struct {
	Connection *websocket.Conn
	ID         string
}

func createClient(c *websocket.Conn) *Client {
	i := <-id

	client := Client{
		Connection: c,
		ID:         i,
	}

	return &client
}
