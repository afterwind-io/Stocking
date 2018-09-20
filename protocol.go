package stocking

// TextMessageInboundProtocol TODO
type TextMessageInboundProtocol struct {
	// Route name
	Route string `json:"route"`
	// message Body
	Body interface{} `json:"body"`
}

// TextMessageOutboundProtocol TODO
type TextMessageOutboundProtocol struct {
	Error string      `json:"error"`
	Body  interface{} `json:"body"`
}
