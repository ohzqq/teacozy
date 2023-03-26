package form

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"golang.org/x/exp/maps"
)

type Model struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
	mainRouter reactea.Component[router.Props]
	Fields     []map[string]string
}

type Form struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FormProps]

	mainRouter reactea.Component[router.Props]
	view       *viewport.Model
}

type FormProps struct {
	Fields   []*FieldProps
	SetValue func(int, string)
}

func New(fields []map[string]string) *Model {
	return &Model{
		mainRouter: router.New(),
		Fields:     fields,
	}
}

func (m *Model) Init(reactea.NoProps) tea.Cmd {
	routes := make(map[string]router.RouteInitializer)

	var fields []*FieldProps
	for idx, field := range m.Fields {
		r := strconv.Itoa(idx)
		for key, val := range field {
			props := FieldProps{
				name:     r,
				idx:      idx,
				key:      key,
				SetValue: m.setFieldValue,
				input:    textinput.New(),
			}
			props.input.SetValue(val)
			fields = append(fields, &props)
		}
	}

	routes["default"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := NewForm()
		props := FormProps{
			Fields:   fields,
			SetValue: m.setFieldValue,
		}
		return comp, comp.Init(props)
	}

	return m.mainRouter.Init(routes)
}

func (c *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}
	return c.mainRouter.Update(msg)
}

func (c *Model) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}

func (c *Model) setFieldValue(idx int, val string) {
	keys := maps.Keys(c.Fields[idx])
	c.Fields[idx][keys[0]] = val
}
