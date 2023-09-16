package list

type Option func(*Model)

func WithHeader(text string) Option {
	return func(m *Model) {
		m.header = text
	}
}

func WithStyle(s Style) Option {
	return func(m *Model) {
		m.Style = s
	}
}

func WithLimit(l int) Option {
	return func(m *Model) {
		m.limit = l
	}
}

func NoLimit() Option {
	return func(m *Model) {
		m.limit = len(m.Choices)
	}
}

func WithHeight(h int) Option {
	return func(m *Model) {
		m.height = h
	}
}

func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

func WithSize(w, h int) Option {
	return func(m *Model) {
		m.width = w
		m.height = h
	}
}
