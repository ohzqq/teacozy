package cmpnt

import "github.com/ohzqq/teacozy"

type Props struct {
	SetCurrent func(int)
	Current    func() int
	Items
}

func NewProps(items teacozy.Items) Props {
	return Props{
		Items: NewItems(items),
	}
}
