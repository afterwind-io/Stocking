package stocking

// Some default startup params
const (
	DefaultHost = ":12345"
	DefaultRoot = "ws"
)

// NewStocking creates and returns a new stocking, server I mean.
func NewStocking(host, root string) *Stocking {
	if host == "" {
		host = DefaultHost
	}

	if root == "" {
		root = DefaultRoot
	}

	return &Stocking{
		Host: host,
		Root: root,
	}
}
