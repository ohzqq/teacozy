package frame

import (
	"github.com/ohzqq/teacozy/item"
)

type Option func(*App)

func WithSlice[E any](c []E) Option {
	return func(a *App) {
		a.choices = item.SliceToChoices(c)
	}
}

func WithMap[K comparable, V any, M ~map[K]V](c []M) Option {
	return func(a *App) {
		a.choices = item.MapToChoices(c)
	}
}

func WithRoute(r Route) Option {
	return func(a *App) {
		//a.NewRoute(r)
		a.Routes[r.Name()] = r
	}
}

func NoLimit() Option {
	return func(a *App) {
		a.noLimit = true
	}
}

func ReadOnly() Option {
	return func(a *App) {
		a.ReadOnly = true
	}
}

func WithLimit(l int) Option {
	return func(a *App) {
		a.limit = l
	}
}

func WithTitle(t string) Option {
	return func(a *App) {
		a.title = t
	}
}

func DefaultRoute(r string) Option {
	return func(a *App) {
		a.defaultRoute = r
	}
}

func WithWidth(w int) Option {
	return func(a *App) {
		a.width = w
	}
}

func WithHeight(h int) Option {
	return func(a *App) {
		a.height = h
	}
}

func WithSize(w, h int) Option {
	return func(a *App) {
		a.width = w
		a.height = h
	}
}

func ConfirmChoices() Option {
	return func(a *App) {
		a.confirmChoices = true
	}
}
