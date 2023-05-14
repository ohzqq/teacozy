package cmpnt

import (
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type Props interface {
	SetCurrent(int)
	Current() int
	IsSelected(int) bool
	Items() teacozy.Items
}

type Page interface {
	Component() reactea.SomeComponent
	KeyMap() keys.KeyMap
}
