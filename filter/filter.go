package filter

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
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
	lineInfo    string
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

func New() *Filter {
	tm := Filter{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func (c Filter) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (c Filter) Name() string {
	return "filter"
}

func (m *Filter) KeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.Up(),
		keys.Down(),
		keys.ToggleItem(),
		keys.Quit(),
		keys.ShowHelp(),
		keys.NewBinding("enter").
			WithHelp("return selections").
			Cmd(m.ReturnSelections()),
		keys.NewBinding("esc").
			WithHelp("stop filtering").
			Cmd(message.StopFiltering()),
	}
}

func (m *Filter) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.UpMsg:
		//h := m.Props().Visible()[m.Cursor].LineHeight()
		offset := m.Viewport.YOffset
		m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor-1)
		if m.Cursor < offset {
			m.Viewport.SetYOffset(m.Cursor)
		}
		//m.lineInfo = fmt.Sprintf("(cursor %d) < (offset %d)\n", m.Cursor, offset)
		m.Props().SetCurrent(m.Cursor)

	case message.DownMsg:
		h := m.Props().Visible()[m.Cursor].LineHeight()
		offset := m.Viewport.YOffset - h
		m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor+1)
		if m.Cursor+h >= offset+m.Viewport.Height {
			m.Viewport.LineDown(h)
		} else if m.Cursor == len(m.Matches)-1 {
			m.Viewport.GotoBottom()
		}
		//m.lineInfo = fmt.Sprintf("down: %d, (cursor %d) >= %d (offset %d + height %d) \n", h, m.Cursor, offset+m.Viewport.Height, offset, m.Viewport.Height)
		m.Props().SetCurrent(m.Cursor)

	case message.ToggleItemMsg:
		idx := m.Matches[m.Cursor].Index

		if m.Props().NumSelected == 0 && m.quitting {
			cmds = append(cmds, message.ReturnSelections())
		}

		m.Props().ToggleSelection(idx)

		if m.Props().Limit == 1 {
			return message.ReturnSelections()
		}

		cmds = append(cmds, message.Down())

	case message.StopFilteringMsg:
		if m.Props().Limit == 1 {
			cmds = append(cmds, message.ToggleItem())
		}
		m.Input.Reset()
		m.Input.Blur()
		return message.ChangeRoute("default")

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, filterKey.StopFiltering):
			cmds = append(cmds, message.StopFiltering())
		case key.Matches(msg, filterKey.Up):
			cmds = append(cmds, message.Up())
		case key.Matches(msg, filterKey.Down):
			cmds = append(cmds, message.Down())
		case key.Matches(msg, filterKey.ToggleItem):
			cmds = append(cmds, message.ToggleItem())
		case key.Matches(msg, filterKey.Quit):
			m.quitting = true
			cmds = append(cmds, message.ReturnSelections())
		case key.Matches(msg, filterKey.ReturnSelections):
			if m.Props().Limit == 1 {
				return message.ToggleItem()
			}
			if m.Props().NumSelected == 0 {
				m.quitting = true
				return message.ToggleItem()
			}
			cmds = append(cmds, message.ReturnSelections())
		}
		m.Input, cmd = m.Input.Update(msg)
		m.Matches = m.Props().Visible(m.Input.Value())
		cmds = append(cmds, cmd)
	}

	m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor)
	m.Props().SetCurrent(m.Cursor)
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

func (m *Filter) ReturnSelections() tea.Cmd {
	return message.ReturnSels(m.Props().Limit, m.Props().NumSelected)
}
