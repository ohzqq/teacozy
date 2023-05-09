package teacozy

import (
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
)

type Props struct {
	*pagy.Paginator
	name       string
	Items      Items
	Selected   map[int]struct{}
	InputValue string
	ReadOnly   bool
	SetCurrent func(int)
	SetHelp    func(keys.KeyMap)
}

func NewProps(items Items) Props {
	p := Props{
		Items:    items,
		Selected: make(map[int]struct{}),
	}
	return p
}
