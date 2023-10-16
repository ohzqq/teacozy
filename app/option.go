package app

import (
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
)

type Option func(m *Model)

func EditableList(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		items := list.NewItems(parser)
		li := list.Edit(items, opts...)
		m.SetList(li)
	}
}

func ChooseAny(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		items := list.NewItems(parser)
		li := list.ChooseAny(items, opts...)
		m.SetList(li)
	}
}

func ChooseOne(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		items := list.NewItems(parser)
		li := list.ChooseOne(items, opts...)
		m.SetList(li)
	}
}

func WithList(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		items := list.NewItems(parser)
		li := list.New(items, opts...)
		m.SetList(li)
	}
}

func WithPager(render pager.Renderer, text ...string) Option {
	return func(m *Model) {
		p := pager.New(render)
		if len(text) > 0 {
			p.SetText(text[0])
		}
		m.SetPager(p)
	}
}

func WithLayout(layout *Layout) Option {
	return func(m *Model) {
		m.layout = layout
	}
}

// WithDescription sets the list to show an item's description.
func WithDescription() Option {
	return func(m *Model) {
		m.showItemDesc = true
		m.SetPager(pager.New(pager.RenderText))
		m.SetLayout(NewLayout().Vertical().Top().Quarter())
	}
}
