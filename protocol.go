package stocking

import "encoding/json"

// TextMessageInboundType TODO
type TextMessageInboundType = string

const (
	tmitConnect   = "0"
	tmitReconnect = "1"
	tmitError     = "2"
	tmitClose     = "3"
	tmitPing      = "4"
	tmitPong      = "5"
	tmitMessage   = "6"
	tmitAck       = "7"
	tmitBroadcast = "8"
)

// TextMessageInboundSegment TODO
const (
	tmisType = iota
	tmisAck
	tmisContent
)

// TextMessageInboundProtocol TODO
type TextMessageInboundProtocol struct {
	// Route name
	Route string `json:"r"`
	// message Body
	Payload json.RawMessage `json:"p"`
}

// TextMessageOutboundProtocol TODO
type TextMessageOutboundProtocol struct {
	Code    int         `json:"c"`
	Payload interface{} `json:"p"`
}
