package cmpnt

import "github.com/ohzqq/teacozy"

type Props interface {
	SetCurrent(int)
	Current() int
	IsSelected(int) bool
	Items() teacozy.Items
}
