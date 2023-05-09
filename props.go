package teacozy

import (
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/sahilm/fuzzy"
)

type Props struct {
	*pagy.Paginator
	name       string
	Items      Items
	Selected   map[int]struct{}
	Search     string
	ReadOnly   bool
	SetCurrent func(int)
	SetHelp    func(keys.KeyMap)
	Style      Style
}

func NewProps(items Items) Props {
	p := Props{
		Items:    items,
		Selected: make(map[int]struct{}),
	}
	return p
}

func (c *Props) ExactMatches(search string) fuzzy.Matches {
	if search != "" {
		if m := fuzzy.FindFrom(search, c.Items); len(m) > 0 {
			return m
		}
	}
	return SourceToMatches(c.Items)
}
