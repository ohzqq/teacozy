package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type Model struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
	mainRouter reactea.Component[router.Props]
	Fields     []map[string]string
}

type FormComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FormProps]
	view *viewport.Model
}

type FormProps struct {
	Field []map[string]string
}

type FieldComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FieldProps]

	input textinput.Model

	key   string
	value string
}

type FieldProps struct {
	SetValue func(string)
}

func (m *Model) Init(reactea.NoProps) tea.Cmd {
}
