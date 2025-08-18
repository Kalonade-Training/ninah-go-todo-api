package vo

import "fmt"

type Description struct{ v string }

func NewDescription(s string) (Description, error) {
	if len([]rune(s)) > 1000 {
		return Description{}, fmt.Errorf("Description is too long")
	}
	return Description{v: s}, nil
}
func (d Description) String() string { return d.v }
