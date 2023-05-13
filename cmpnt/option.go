package cmpnt

import "github.com/ohzqq/teacozy"

type Option func(*Pager)

func WithSlice[E any](c []E) Option {
	return func(a *Pager) {
		a.Choices = teacozy.SliceToChoices(c)
	}
}

func WithMap[K comparable, V any, M ~map[K]V](c []M) Option {
	return func(a *Pager) {
		a.Choices = teacozy.MapToChoices(c)
	}
}

func NoLimit() Option {
	return func(a *Pager) {
		a.NoLimit = true
	}
}

func ReadOnly() Option {
	return func(a *Pager) {
		a.readOnly = true
	}
}

func WithLimit(l int) Option {
	return func(a *Pager) {
		a.Limit = l
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

func ConfirmChoices() Option {
	return func(a *Pager) {
		a.ConfirmChoices = true
	}
}
