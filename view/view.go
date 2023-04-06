package view

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type View struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	quitting bool
	Style    style.List
	list     *list.List
}

type Props struct {
	*props.Items
}

func NewProps(items *props.Items) Props {
	return Props{
		Items: items,
	}
}

func New() *View {
	return &View{
		Style: style.ListDefaults(),
	}
}

func (m View) KeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.NewBinding("esc", "q").WithHelp("exit screen").Cmd(message.HideHelp()),
		keys.Quit(),
		keys.ShowHelp(),
	}
	return km
}

func (m View) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (m View) Name() string {
	return "view"
}

func (m *View) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.HideHelpMsg:
		return message.ChangeRoute("prev")
	case message.QuitMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		for _, k := range m.KeyMap() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *View) Render(w, h int) string {
	if m.list.Footer != "" {
		m.Props().SetFooter(m.list.Footer)
	}
	return m.list.View()
}

func (m *View) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	m.list = list.New(props.Items)
	return nil
}

type ShowViewMsg struct{}

func ShowView() tea.Cmd {
	return func() tea.Msg {
		return ShowViewMsg{}
	}
}

type HideViewMsg struct{}

func HideView() tea.Cmd {
	return func() tea.Msg {
		return HideViewMsg{}
	}
}
