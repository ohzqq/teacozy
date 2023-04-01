package help

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Help struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor   int
	Viewport *viewport.Model
	quitting bool
	Style    style.List
	lineInfo string
	keys     KeyMap
}

var keys = KeyMap{
	NewBinding("esc").WithHelp("exit screen").Cmd(message.HideHelpCmd()),
	Up(),
	Down(),
	Quit(),
	ShowHelp(),
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
	case message.QuitMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)

	case message.UpMsg:
		offset := m.Viewport.YOffset
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor-1)
		if m.Cursor < offset {
			m.Viewport.SetYOffset(m.Cursor)
		}
		m.Props().SetCurrent(m.Cursor)

	case message.DownMsg:
		h := m.Props().Visible()[m.Cursor].LineHeight()
		offset := m.Viewport.YOffset - h
		m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor+1)
		if m.Cursor+h >= offset+m.Viewport.Height {
			m.Viewport.LineDown(h)
		} else if m.Cursor == len(m.Props().Visible())-1 {
			m.Viewport.GotoBottom()
		}
		m.Props().SetCurrent(m.Cursor)

	case tea.KeyMsg:
		for _, k := range keys {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
		cmds = append(cmds, cmd)
	}

	m.Cursor = util.Clamp(0, len(m.Props().Visible())-1, m.Cursor)
	m.Props().SetCurrent(m.Cursor)
	return tea.Batch(cmds...)
}

func (m *Help) Render(w, h int) string {
	m.Viewport.Height = m.Props().Height
	m.Viewport.Width = m.Props().Width

	var s strings.Builder
	items := m.Props().RenderItems(m.Cursor, m.Props().Visible())
	s.WriteString(items)

	m.Viewport.SetContent(s.String())

	view := m.Viewport.View()
	return view
}

func (tm *Help) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)

	v := viewport.New(0, 0)
	tm.Viewport = &v

	return nil
}
