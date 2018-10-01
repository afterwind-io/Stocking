package stocking

// RouterPackage TODO
type RouterPackage struct {
	Route string
	Body  interface{}
}

// RouterHandler TODO
type RouterHandler = func(p RouterPackage) (interface{}, error)

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
