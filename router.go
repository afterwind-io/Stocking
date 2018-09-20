package stocking

// RouterPackage TODO
type RouterPackage struct {
	Route string
	Body  interface{}
}

// RouterHandler TODO
type RouterHandler = func(p RouterPackage) (interface{}, error)
