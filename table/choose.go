package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

// Choose defines a state for the table widget.
type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List

	list *List

	Viewport viewport.Model
}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))

func (c Choose) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewChoice()
		return component, component.Init(Props{Items: props})
	}
}

func (c Choose) Name() string {
	return "choose"
}

// New creates a new model for the table widget.
func NewChoice() *Choose {
	m := Choose{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}

	return &m
}

func (m *Choose) Init(props Props) tea.Cmd {
	m.UpdateProps(props)

	m.list = New(props.Items)
	m.list.Focus()

	m.Viewport = viewport.New(props.Width, props.Height)

	return nil
}

func (m Choose) KeyMap() keys.KeyMap {
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
			Cmd(message.Bottom()),
		keys.NewBinding("g").
			WithHelp("list top").
			Cmd(message.Top()),
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

func (m *Choose) ReturnSelections() tea.Cmd {
	return message.ReturnSels(m.Props().Limit, m.Props().NumSelected)
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case message.QuitMsg:
		return tea.Quit

	case message.ToggleItemMsg:
		if len(m.list.VisibleItems()) > 0 {
			m.Props().ToggleSelection(m.list.CurrentItem().Index)
			switch {
			case m.Props().NumSelected == 0 && m.quitting:
				cmds = append(cmds, m.ReturnSelections())
			case m.Props().Limit == 1:
				cmds = append(cmds, m.ReturnSelections())
			case m.Props().NumSelected > 0 || m.Props().Limit > 1:
				cmds = append(cmds, message.Down(1))
			}

		}

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

func (m Choose) Render(w, h int) string {
	m.Viewport.Height = m.Props().Height - 1
	m.Viewport.Width = m.Props().Width
	m.list.UpdateItems()

	view := m.list.View()

	return view
}

func (m *Choose) ToggleAllItems() tea.Cmd {
	return func() tea.Msg {
		var items []int
		for _, item := range m.Props().Items.Items {
			items = append(items, item.Index)
		}
		m.Props().ToggleSelection(items...)
		return nil
	}
}

func (m *Choose) quit() tea.Cmd {
	m.quitting = true
	return message.Quit()
}
