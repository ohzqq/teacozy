package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

// Filter defines a state for the table widget.
type Filter struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List

	list *List

	Input textinput.Model
}

type Props struct {
	*props.Items
}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))
type Option func(*Filter)

func (c Filter) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewTable()
		return component, component.Init(Props{Items: props})
	}
}

func (c Filter) Name() string {
	return "table"
}

// New creates a new model for the table widget.
func NewTable(opts ...Option) *Filter {
	m := Filter{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return &m
}

func (m *Filter) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	m.Input = textinput.New()
	m.Input.Prompt = m.Prompt
	m.Input.PromptStyle = m.Style.Prompt
	m.Input.Placeholder = m.Placeholder
	m.Input.Width = props.Width
	m.Input.Focus()

	m.list = New(props.Items)
	m.list.Focus()

	m.list.SetWidth(props.Width)
	m.list.SetHeight(props.Height)

	return nil
}

func (m Filter) KeyMap() keys.KeyMap {
	var km = keys.KeyMap{
		//keys.ToggleItem(),
		keys.NewBinding("enter").
			WithHelp("return selections").
			Cmd(m.ReturnSelections()),
		keys.NewBinding("/").
			WithHelp("filter list").
			Cmd(StartFiltering()),
		keys.NewBinding("esc").
			WithHelp("stop filtering").
			Cmd(StopFiltering()),
	}
	return km
}

func (m *Filter) ReturnSelections() tea.Cmd {
	return message.ReturnSels(m.Props().Limit, m.Props().NumSelected)
}

func (m *Filter) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case message.QuitMsg:
		return reactea.Destroy

	case StopFilteringMsg:
		m.Input.Reset()
		m.Input.Blur()
		m.list.Focus()

	case StartFilteringMsg:
		m.Input.Focus()

	case message.ToggleItemMsg:
		if len(m.list.VisibleItems()) > 0 {
			m.Props().ToggleSelection(m.list.CurrentItem().Index)
			switch {
			case m.Props().NumSelected == 0 && m.quitting:
				cmds = append(cmds, m.ReturnSelections())
			case m.Props().Limit == 1:
				cmds = append(cmds, m.ReturnSelections())
			case m.Props().NumSelected > 0 || m.Props().Limit > 1:
				cmds = append(cmds, message.LineDown)
			}

		}

	case tea.KeyMsg:
		if m.Input.Focused() {
			m.Input, cmd = m.Input.Update(msg)
			m.Props().Filter(m.Input.Value())
			cmds = append(cmds, cmd)
		}
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

func (m Filter) Render(w, h int) string {
	m.list.SetHeight(m.Props().Height - 1)
	m.list.SetWidth(m.Props().Width)
	m.list.UpdateItems()

	view := m.list.View()
	if m.Input.Focused() {
		view = m.Input.View() + "\n" + view
	}

	return view
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func (m *Filter) ToggleAllItems() tea.Cmd {
	return func() tea.Msg {
		var items []int
		for _, item := range m.Props().Items.Items {
			items = append(items, item.Index)
		}
		m.Props().ToggleSelection(items...)
		return nil
	}
}

func (m Filter) UnfilteredKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Up().WithKeys("k", "up"),
		keys.Down().WithKeys("j", "down"),
		keys.Next().WithKeys("right", "l"),
		keys.Prev().WithKeys("left", "h"),
		keys.NewBinding("/").
			WithHelp("filter list").
			Cmd(StartFiltering()),
		keys.NewBinding("G").
			WithHelp("list bottom").
			Cmd(message.Bottom),
		keys.NewBinding("g").
			WithHelp("list top").
			Cmd(message.Top),
		keys.NewBinding("v").
			WithHelp("toggle all items").
			Cmd(m.ToggleAllItems()),
		keys.NewBinding("enter").
			WithHelp("return selections").
			Cmd(m.ReturnSelections()),
		keys.ToggleItem().WithKeys("tab", " "),
		keys.ShowHelp(),
		keys.Quit().
			WithKeys("ctrl+c", "q", "esc").
			Cmd(m.quit()),
	}
	return km
}

func (m *Filter) quit() tea.Cmd {
	m.quitting = true
	return message.Quit()
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}

type StopFilteringMsg struct{}

func StopFiltering() tea.Cmd {
	return func() tea.Msg {
		return StopFilteringMsg{}
	}
}

type StartFilteringMsg struct{}

func StartFiltering() tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}
