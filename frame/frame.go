package frame

import "github.com/charmbracelet/bubbles/viewport"

type Model struct {
	model  viewport.Model
	Style  Styles
	Width  int
	Height int
}

func New() Model {
	style := DefaultStyle()
	model := viewport.New()
	return Model{}
}
