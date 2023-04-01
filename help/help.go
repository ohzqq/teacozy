package help

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Help struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor    int
	Paginator paginator.Model
	quitting  bool
	Style     style.List
	keys      keys.KeyMap
}

var Keys = keys.KeyMap{
	keys.NewBinding("esc").WithHelp("exit screen").Cmd(message.HideHelpCmd()),
	keys.Up(),
	keys.Down(),
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
	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))
	switch msg := msg.(type) {
	case message.HideHelpMsg:
		return message.ChangeRouteCmd("default")
	case message.QuitMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)

	case message.NextMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor+m.Props().Height)
		m.Props().Items.SetCurrent(m.Cursor)
		m.Paginator.NextPage()

	case message.PrevMsg:
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor-m.Props().Height)
		m.Props().Items.SetCurrent(m.Cursor)
		m.Paginator.PrevPage()

	case message.TopMsg:
		m.Cursor = 0
		m.Paginator.Page = 0
		m.Props().SetCurrent(m.Cursor)

	case message.BottomMsg:
		m.Cursor = len(m.Props().Visible()) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
		m.Props().SetCurrent(m.Cursor)

	case message.UpMsg:
		m.Cursor--
		if m.Cursor < 0 {
			m.Cursor = len(m.Props().Visible()) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		}
		if m.Cursor < start {
			m.Paginator.PrevPage()
		}
		m.Props().SetCurrent(m.Cursor)

	case message.DownMsg:
		m.Cursor++
		if m.Cursor >= len(m.Props().Visible()) {
			m.Cursor = 0
			m.Paginator.Page = 0
		}
		if m.Cursor >= end {
			m.Paginator.NextPage()
		}
		m.Props().SetCurrent(m.Cursor)

	case tea.KeyMsg:
		for _, k := range Keys {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *Help) Render(w, h int) string {
	var s strings.Builder
	start, end := m.Paginator.GetSliceBounds(len(m.Props().Visible()))

	items := m.Props().RenderItems(
		m.Cursor%m.Props().Height,
		m.Props().Visible()[start:end],
	)
	s.WriteString(items)

	var view string
	view = s.String()
	if m.Paginator.TotalPages <= 1 {
		//view = s.String()
	} else if m.Paginator.TotalPages > 1 {
		p := style.Footer.Render(m.Paginator.View())
		//view = lipgloss.JoinVertical(lipgloss.Left, view, p)
		m.Props().Footer(p)
	}

	return view
}

func (tm *Help) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)

	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Arabic
	tm.Paginator.SetTotalPages((len(tm.Props().Visible()) + props.Height - 1) / props.Height)
	tm.Paginator.PerPage = props.Height

	return nil
}
