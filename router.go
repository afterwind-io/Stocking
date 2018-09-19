package stocking

// RouterPackage TODO
type RouterPackage struct {
	route string
	body  string
}

// RouterHandler TODO
type RouterHandler = func(p RouterPackage)

