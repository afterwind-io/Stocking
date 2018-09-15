package stocking

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Stocking is the instance of the websocket server
// with basic configs
type Stocking struct {
	// the address for http.ListenAndServe
	Host string
	// the root pattern for connection
	Root string
	// the underlying Upgrader for gorilla/websocket
	Upgrader websocket.Upgrader
}

// Start the server
func (s *Stocking) Start() {
	upgrader = s.Upgrader

	go generateID()

	http.HandleFunc("/"+s.Root, handler)
	log.Fatal(http.ListenAndServe(s.Host, nil))
}

var upgrader websocket.Upgrader
var clients = make(map[*Client]bool)
var rooms = make(map[string]*Client)

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade failed: ", err)
	}

	client := createClient(c)
	clients[client] = true

	go serveClient(client)
}

func serveClient(c *Client) {
	defer c.Connection.Close()

	for {
		_, content, err := c.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println(fmt.Sprintf("[%v] %s", c.ID, content))
	}
}
