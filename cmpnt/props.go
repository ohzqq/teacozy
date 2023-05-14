package cmpnt

import (
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type PageProps interface {
	SetCurrent(int)
	Current() int
	IsSelected(int) bool
	Items() teacozy.Items
}

type PageComponent interface {
	Mount() reactea.SomeComponent
	KeyMap() keys.KeyMap
}

func NewPageComponent() func(PageProps) PageComponent {
	return func(props PageProps) PageComponent {
		p := New()
		p.Init(props)
		return p
	}
}
