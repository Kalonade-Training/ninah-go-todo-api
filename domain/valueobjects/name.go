package valueobjects

import "errors"

const MaxNameLen = 50

type Name string

func NewName(s string) (Name, error) {
	if len([]rune(s)) == 0 {
		return "", errors.New("name required")
	}
	if len([]rune(s)) > MaxNameLen {
		return "", errors.New("name is over too long")
	}
	return Name(s), nil
}

func (n Name) String() string { return string(n) }
