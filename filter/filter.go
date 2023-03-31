package filter

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Filter struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor      int
	Matches     []props.Item
	Input       textinput.Model
	Viewport    *viewport.Model
	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List
}

type Props struct {
	*props.Items
}

type KeyMap struct {
	Up               key.Binding
	Down             key.Binding
	ToggleItem       key.Binding
	Quit             key.Binding
	ReturnSelections key.Binding
	StopFiltering    key.Binding
}

func NewFilter() *Filter {
	tm := Filter{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func RouteInitializer(props Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewFilter()
		return component, component.Init(props)
	}
}

func (c Filter) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewFilter()
		return component, component.Init(Props{Items: props})
	}
}

func (c Filter) Name() string {
	return "filter"
}

func (m *Filter) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case message.ReturnSelectionsMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case message.UpMsg:
		m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor-1)
		if m.Cursor < m.Viewport.YOffset {
			m.Viewport.SetYOffset(m.Cursor)
		}
	case message.DownMsg:
		h := lipgloss.Height(m.Props().Visible()[m.Cursor].Str)
		m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor+1)
		if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height-h {
			m.Viewport.LineDown(h)
		} else if m.Cursor == len(m.Matches)-1 {
			m.Viewport.GotoBottom()
		}
	case message.ToggleItemMsg:
		if m.Props().Limit == 1 {
			return nil
		}
		if m.Cursor >= 0 {
			idx := m.Props().Visible()[m.Cursor].Index
			m.Props().ToggleSelection(idx)
		}
		cmds = append(cmds, message.DownCmd())
	case message.StopFilteringMsg:
		if m.Props().Limit == 1 {
			cmds = append(cmds, message.ToggleItemCmd())
		}

		m.Input.Reset()
		m.Input.Blur()
		return message.ChangeRouteCmd("default")
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, filterKey.StopFiltering):
			cmds = append(cmds, message.StopFilteringCmd())
		case key.Matches(msg, filterKey.Up):
			cmds = append(cmds, message.UpCmd())
		case key.Matches(msg, filterKey.Down):
			cmds = append(cmds, message.DownCmd())
		case key.Matches(msg, filterKey.ToggleItem):
			cmds = append(cmds, message.ToggleItemCmd())
		case key.Matches(msg, filterKey.Quit):
			m.quitting = true
			cmds = append(cmds, message.ReturnSelectionsCmd())
		case key.Matches(msg, filterKey.ReturnSelections):
			cmds = append(cmds, message.ReturnSelectionsCmd())
		}
		m.Input, cmd = m.Input.Update(msg)
		m.Matches = m.Props().Visible(m.Input.Value())
		cmds = append(cmds, cmd)
	}

	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor)
	return tea.Batch(cmds...)
}

func (m *Filter) Render(w, h int) string {
	m.Viewport.Height = m.Props().Height
	m.Viewport.Width = m.Props().Width

	var s strings.Builder
	items := m.Props().RenderItems(m.Cursor, m.Matches)
	s.WriteString(items)

	m.Viewport.SetContent(s.String())

	view := m.Input.View() + "\n" + m.Viewport.View()
	return view
}

func (tm *Filter) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)
	tm.Matches = tm.Props().Visible()

	tm.Input = textinput.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.PromptStyle = tm.Style.Prompt
	tm.Input.Placeholder = tm.Placeholder
	tm.Input.Width = tm.Props().Width

	v := viewport.New(0, 0)
	tm.Viewport = &v
	tm.Input.Focus()

	return nil
}

var filterKey = KeyMap{
	ToggleItem: key.NewBinding(
		key.WithKeys(" ", "tab"),
		key.WithHelp("space", "select item"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "move cursor down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move cursor up"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	StopFiltering: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "stop filtering"),
	),
	ReturnSelections: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "return selections"),
	),
}

func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
