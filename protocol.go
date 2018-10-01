package stocking

import "encoding/json"

// TextMessageInboundType TODO
const (
	tmitConnect   = "0"
	tmitError     = "1"
	tmitClose     = "2"
	tmitPingPong  = "3"
	tmitMessage   = "4"
	tmitBroadcast = "5"
	tmitJoin      = "6"
)

// TextMessageInboundSegment TODO
const (
	tmisType = iota
	tmisCCode
	tmisContent
)

// TextMessageProtocol TODO
type TextMessageProtocol struct {
	Event   string          `json:"e"`
	Payload json.RawMessage `json:"p"`
}
