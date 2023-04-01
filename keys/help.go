package keys

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Help struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor   int
	Viewport *viewport.Model
	quitting bool
	Style    style.List
	lineInfo string
}

type Props struct {
	*props.Items
}

func NewProps(items *props.Items) Props {
	return Props{
		Items: items,
	}
}

func New() *Help {
	return &Help{
		Style: style.ListDefaults(),
	}
}

func (h Help) Name() string {
	return "help"
}

func (h *Help) Init(props Props) tea.Cmd {
	h.UpdateProps(props)
	return nil
}
