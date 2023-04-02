package help

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Help struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor    int
	Paginator paginator.Model
	quitting  bool
	Style     style.List
	list      *list.List
	keys      keys.KeyMap
}

var Keys = keys.KeyMap{
	keys.NewBinding("esc").WithHelp("exit screen").Cmd(message.HideHelp()),
	//keys.Up(),
	//keys.Down(),
	keys.Quit(),
	keys.ShowHelp(),
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

func (c Help) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (h Help) Name() string {
	return "help"
}

func (m *Help) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.HideHelpMsg:
		return message.ChangeRoute("default")
	case message.QuitMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		for _, k := range Keys {
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

func (m *Help) Render(w, h int) string {
	return m.list.View()
}

func (tm *Help) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)

	tm.list = list.New(props.Items)

	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Arabic
	tm.Paginator.SetTotalPages((len(tm.Props().Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height

	return nil
}
