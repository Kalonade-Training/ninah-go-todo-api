package vo

import "fmt"

type Title struct{ v string }

func NewTitle(s string) (Title, error) {
	if n := len([]rune(s)); n < 1 || n > 200 {
		return Title{}, fmt.Errorf("Title must be 1 to 200 characters")
	}
	return Title{v: s}, nil
}
func (t Title) String() string { return t.v }
