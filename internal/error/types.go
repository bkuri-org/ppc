package error

import "fmt"

type SrcError struct {
	Path string
	ID   string
	Msg  string
	Line int
}

func (e SrcError) Error() string {
	return fmt.Sprintf("%s [%s]", e.Msg, e.ID)
}

func (e SrcError) Unwrap() error {
	return nil
}

func New(path, id, msg string) SrcError {
	return SrcError{Path: path, ID: id, Msg: msg}
}

func NewAtLine(path, id, msg string, line int) SrcError {
	return SrcError{Path: path, ID: id, Msg: msg, Line: line}
}
