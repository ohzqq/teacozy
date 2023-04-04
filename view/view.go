package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/style"
)

type Filter struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Viewport *viewport.Model
	quitting bool
	Style    style.List
}

type Props struct {
	Header string
	Body   string
	Footer string
}

func Renderer(props Props, w, h int) string {
	vp := viewport.New(w, h)
	vp.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Left,
			props.Header,
			props.Body,
			props.Footer,
		),
	)
	return vp.View()
}
