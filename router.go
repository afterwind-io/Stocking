package stocking

// RouterPackage TODO
type RouterPackage struct {
	Client *Client
	Route  string
	Body   interface{}
}

// RouterHandler TODO
type RouterHandler = func(p *RouterPackage) (interface{}, error)

// RouterMessageProtocol TODO
type RouterMessageProtocol struct {
	Code    int         `json:"c"`
	Payload interface{} `json:"p"`
}

// RouterError TODO
type RouterError struct {
	code int
	msg  string
}

func (e RouterError) Error() string {
	return e.msg
}

// NewTextRouterError TODO
func NewTextRouterError(msg string) RouterError {
	return RouterError{0, msg}
}
