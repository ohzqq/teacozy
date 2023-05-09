package teacozy

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type Page interface {
	Header() reactea.SomeComponent
	Main() reactea.SomeComponent
	Footer() reactea.SomeComponent
}

type BasicPage struct {
	reactea.BasicComponent
	header string
	main   string
	footer string
}

func NewPage() *BasicPage {
	return &BasicPage{
		header: "header",
		main:   "main",
		footer: "footer",
	}
}

func (c *BasicPage) Header() reactea.BasicComponent {
	return reactea.Componentify[string](c.Render)
}

func (c *BasicPage) Main() reactea.BasicComponent {
	return reactea.Componentify[string](c.Render)
}

func (c *BasicPage) Footer() reactea.BasicComponent {
	return reactea.Componentify[string](c.Render)
}

func (c *BasicPage) Render(w, h int) string {
	return lipgloss.JoinVertical(lipgloss.Left, c.header, c.main, c.footer)
}
