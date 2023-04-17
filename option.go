package teacozy

import "github.com/ohzqq/teacozy/item"

type Option func(*App)

func WithSlice[E any](c []E) Option {
	return func(a *App) {
		a.Choices = item.SliceToChoices(c)
	}
}

func WithMap(c []map[string]string) Option {
	return func(a *App) {
		a.Choices = item.MapToChoices(c)
	}
}

func NoLimit() Option {
	return func(a *App) {
		a.limit = -1
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
func ConfirmChoices() Option {
	return func(a *App) {
		a.confirmChoices = true
	}
}

func Editable() Option {
	return func(a *App) {
		a.editable = true
	}
}

func WithFilter() Option {
	return func(a *App) {
		a.filterable = true
	}
}
