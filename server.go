package stocking

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Some default startup params
const (
	defaultHost = ":12345"
	defaultRoot = "ws"
)

var (
	upgrader    websocket.Upgrader
	hub         = newHub()
	router      = newRouter()
	middlewares = []Middleware{}
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
	s.Attach()

	log.Fatal(http.ListenAndServe(s.Host, nil))
}

// Attach root handler to an existing http server
func (s *Stocking) Attach() {
	s.boot()

	http.HandleFunc("/"+s.Root, serveClient)
}

// On adds a route handler
func (s *Stocking) On(route string, handler RouterHandler, typeHint interface{}) {
	router.On(route, handler, typeHint)
}

// Otherwise adds a fallback handler when no route hits
func (s *Stocking) Otherwise(handler RouterHandler) {
	router.Otherwise(handler)
}

// Use adds a middleware
func (s *Stocking) Use(ms ...Middleware) {
	middlewares = append(middlewares, ms...)
}

func (s *Stocking) boot() {
	upgrader = s.Upgrader

	// TODO
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	go generateID()

	hub.use(&mLogger{})
	hub.use(middlewares...)
	hub.use(router)
	go hub.run()
}

// NewStocking creates and returns a new stocking, server I mean.
func NewStocking(host, root string) *Stocking {
	if host == "" {
		host = defaultHost
	}

	if root == "" {
		root = defaultRoot
	}

	return &Stocking{
		Host: host,
		Root: root,
	}
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade failed: ", err)
		return
	}

	client := newClient(c, hub)
	hub.signin <- client

	go client.read()
	go client.write()
}
