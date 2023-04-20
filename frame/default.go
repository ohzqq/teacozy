package frame

import (
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/item"
)

type Default struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[item.Props]
}

func NewDefault() *Default {
	return &Default{}
}

func (c *Default) Render(w, h int) string {
	view := item.Renderer(c.Props(), w, h)
	return view
}
