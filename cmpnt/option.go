package cmpnt

import "github.com/ohzqq/teacozy"

type Option func(*Pager)

func WithSlice[E any](c []E) Option {
	return func(a *Pager) {
		a.Items = teacozy.SliceToChoices(c)
	}
}

func WithMap[K comparable, V any, M ~map[K]V](c []M) Option {
	return func(a *Pager) {
		a.Items = teacozy.MapToChoices(c)
	}
}

func WithWidth(w int) Option {
	return func(a *Pager) {
		a.Height = w
	}
}

func WithHeight(h int) Option {
	return func(a *Pager) {
		a.Height = h
	}
}

func WithSize(w, h int) Option {
	return func(a *Pager) {
		a.Height = w
		a.Height = h
	}
}
