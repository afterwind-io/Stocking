package stocking

// JSONSyntaxError TODO
type JSONSyntaxError struct {
	msg string
}

func (err JSONSyntaxError) Error() string {
	return err.msg
}
