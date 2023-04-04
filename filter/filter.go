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
	start       int
	end         int
}

type Props struct {
	*props.Items
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
		m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor-1)

		h := m.Matches[m.Cursor].LineHeight()
		m.Viewport.LineUp(h)

		m.Props().SetCurrent(m.Cursor)

	case message.DownMsg:
		m.Cursor = util.Clamp(0, len(m.Matches)-1, m.Cursor+1)

		offset := m.Viewport.YOffset
		h := m.Matches[m.Cursor].LineHeight()
		if o := h - m.Viewport.Height; o > 0 {
			m.Viewport.LineDown(o)
		} else if m.Cursor <= offset+m.Viewport.Height {
			m.Viewport.LineDown(h)
		}

		m.Props().SetCurrent(m.Cursor)

	case message.ToggleItemMsg:
		if len(m.Matches) == 0 {
			return nil
		}
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

	case message.ShowHelpMsg:
		k := m.KeyMap()
		m.Props().SetHelp(k)
		cmds = append(cmds, message.ChangeRoute("help"))

	case tea.KeyMsg:
		for _, k := range m.KeyMap() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
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

	items := m.Props().RenderItems(m.Matches)

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

func (m *Filter) ReturnSelections() tea.Cmd {
	return message.ReturnSels(m.Props().Limit, m.Props().NumSelected)
}
