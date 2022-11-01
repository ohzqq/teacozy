package multi

import (
	"github.com/ohzqq/urbooks-core/bubbles/list"
)

type List struct {
	*list.List
}

func New(title string, items list.Items) *List {
	l := list.New(title)
	l.SetItems(items)
	l.Model = l.BuildModel()
	return &List{
		List: l,
	}
}
