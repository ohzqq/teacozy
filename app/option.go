package app

import (
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
)

// Option configures a Model.
type Option func(m *Model)

// EditableList is a convenience method for a list where items can be inserted
// or removed.
func EditableList(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		o := []list.Option{list.Editable(true)}
		o = append(o, opts...)
		m.SetList(parser, opts...)
	}
}

// ChooseAny is a convenience method for a multi-selectable list.
func ChooseAny(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		opts = append(opts, list.WithLimit(list.SelectAll))
		m.SetList(parser, opts...)
	}
}

// ChooseOne is a convenience method for a single option list.
func ChooseOne(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		opts = append(opts, list.WithLimit(list.SelectOne))
		m.SetList(parser, opts...)
	}
}

// WithList configures a list.
func WithList(parser list.ParseItems, opts ...list.Option) Option {
	return func(m *Model) {
		m.SetList(parser, opts...)
	}
}

// WithPager configures a pager.
func WithPager(render pager.Renderer, text ...string) Option {
	return func(m *Model) {
		m.SetPager(render, text...)
	}
}

// WithLayout sets the Model's layout.
func WithLayout(layout *Layout) Option {
	return func(m *Model) {
		m.layout = layout
	}
}

// WithDescription sets the list to show an item's description.
func WithDescription() Option {
	return func(m *Model) {
		m.showItemDesc = true
		m.SetPager(pager.RenderText)
		m.SetLayout(NewLayout().Vertical().Top().Quarter())
	}
}
