package valueobjects

import "errors"

const MaxBodyLen = 1000

type Body string

func NewBody(s string) (Body, error) {
	if len([]rune(s)) > MaxBodyLen {
		return "", errors.New("body is too long")
	}
	return Body(s), nil
}

func (b Body) String() string { return string(b) }
