package cmpnt

import "github.com/ohzqq/teacozy"

type Option func(*Page)

func WithSlice[E any](c []E) Option {
	return func(a *Page) {
		a.Choices = teacozy.SliceToChoices(c)
	}
}

func WithMap[K comparable, V any, M ~map[K]V](c []M) Option {
	return func(a *Page) {
		a.Choices = teacozy.MapToChoices(c)
	}
}

func NoLimit() Option {
	return func(a *Page) {
		a.NoLimit = true
	}
}

func ReadOnly() Option {
	return func(a *Page) {
		a.readOnly = true
	}
}

func WithLimit(l int) Option {
	return func(a *Page) {
		a.Limit = l
	}
}

func WithTitle(t string) Option {
	return func(a *Page) {
		a.Title = t
	}
}

func WithWidth(w int) Option {
	return func(a *Page) {
		a.Height = w
	}
}

func WithHeight(h int) Option {
	return func(a *Page) {
		a.Height = h
	}
}

func WithSize(w, h int) Option {
	return func(a *Page) {
		a.Height = w
		a.Height = h
	}
}

func ConfirmChoices() Option {
	return func(a *Page) {
		a.ConfirmChoices = true
	}
}
