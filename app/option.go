package app

type Option func(*Page)

func WithSlice[E any](c []E) Option {
	return func(a *Page) {
		a.choices = SliceToChoices(c)
	}
}

func WithMap[K comparable, V any, M ~map[K]V](c []M) Option {
	return func(a *Page) {
		a.choices = MapToChoices(c)
	}
}

func WithRoute(r Route) Option {
	return func(a *Page) {
		a.Routes[r.Name()] = r
	}
}

func NoLimit() Option {
	return func(a *Page) {
		a.noLimit = true
	}
}

func ReadOnly() Option {
	return func(a *Page) {
		a.readOnly = true
	}
}

func WithLimit(l int) Option {
	return func(a *Page) {
		a.limit = l
	}
}

func WithTitle(t string) Option {
	return func(a *Page) {
		a.title = t
	}
}

func DefaultRoute(r string) Option {
	return func(a *Page) {
		a.defaultRoute = r
	}
}

func WithWidth(w int) Option {
	return func(a *Page) {
		a.width = w
	}
}

func WithHeight(h int) Option {
	return func(a *Page) {
		a.height = h
	}
}

func WithSize(w, h int) Option {
	return func(a *Page) {
		a.width = w
		a.height = h
	}
}

func ConfirmChoices() Option {
	return func(a *Page) {
		a.confirmChoices = true
	}
}
