package choose

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	quitting bool
	header   string
	list     *list.List
	Style    style.List
}

type Props struct {
	*props.Items
}

func New() *Choose {
	tm := Choose{
		Style: style.ListDefaults(),
	}
	return &tm
}

func (c Choose) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		p := Props{
			Items: props,
		}
		return component, component.Init(p)
	}
}

func (c Choose) Name() string {
	return "choose"
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case message.ShowHelpMsg:
		k := m.KeyMap()
		k = append(k, list.Keys...)
		m.Props().SetHelp(k)
		cmds = append(cmds, message.ChangeRoute("help"))

	case message.StartEditingMsg:
		return message.ChangeRoute("editField")

	case message.ToggleItemMsg:
		m.SetCurrent()
		m.Props().ToggleSelection()
		if m.Props().Limit == 1 {
			return message.ReturnSelections()
		}
		cmds = append(cmds, message.Down())

	case message.StartFilteringMsg:
		return message.ChangeRoute("filter")

	case tea.KeyMsg:
		for _, k := range m.KeyMap() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *Choose) SetCurrent() {
	m.Props().SetCurrent(m.Props().Visible()[m.list.Cursor].Index)
}

func (m *Choose) Render(w, h int) string {
	if m.list.Footer != "" {
		m.Props().SetFooter(m.list.Footer)
	}
	return m.list.View()
}

func (tm *Choose) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)
	tm.list = list.New(props.Items)
	return nil
}
