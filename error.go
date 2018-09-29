package stocking

type sError struct {
	msg string
}

func (err sError) Error() string {
	return err.msg
}
