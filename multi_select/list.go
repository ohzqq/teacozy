package multi

import (
	"github.com/ohzqq/teacozy/list"
)

type List struct {
	*list.Model
}

func New(title string, items list.Items) *List {
	l := list.New(title)
	l.SetItems(items)
	l.List = l.BuildModel()
	return &List{
		List: l,
	}
}
