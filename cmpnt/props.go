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

type Page interface {
	Mount() reactea.SomeComponent
	KeyMap() keys.KeyMap
}

func NewPage() func(PageProps) Page {
	return func(props PageProps) Page {
		p := New()
		p.Init(props)
		return p
	}
}
