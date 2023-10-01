package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/bubbles/textinput"
	"github.com/ohzqq/teacozy/util"
)

type State int

const (
	Browsing State = iota
	Input
)

type Model struct {
	*list.Model
	itemDelegate list.DefaultDelegate
	items        Items
	width        int
	height       int
	state        State
	input        textinput.Model
	hasInput     bool
	editable     bool
}

type ListOpt func(*Model)
type ListConfig func(*Model) list.Model

func New(items Items, opts ...ListOpt) *Model {
	w, h := util.TermSize()
	m := &Model{
		width:  w,
		height: h,
		items:  items,
		state:  Browsing,
	}

	m.Model = m.NewListModel(items)

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func EditableList() ListOpt {
	return func(m *Model) {
		m.editable = true
		m.Model.SetLimit(0)
		m.Model.SetFilteringEnabled(false)

		help := func() []key.Binding {
			return []key.Binding{
				m.Model.KeyMap.InsertItem,
				m.Model.KeyMap.RemoveItem,
			}
		}

		m.Model.AdditionalShortHelpKeys = help
		m.Model.AdditionalFullHelpKeys = help

		WithInput("Insert Item: ")(m)
	}
}

func WithInput(prompt string) ListOpt {
	return func(m *Model) {
		m.hasInput = true
		m.input = m.NewTextinputModel()
		m.input.Prompt = prompt
	}
}

// State returns the current filter state.
func (m Model) State() State {
	return m.state
}

// NewTextinputModel returns a textinput.Model with the default styles.
func (m Model) NewTextinputModel() textinput.Model {
	input := textinput.New()
	input.PromptStyle = m.Styles.FilterPrompt
	input.Cursor.Style = m.Styles.FilterCursor
	return input
}

func (m Model) NewListModel(items Items) *list.Model {
	var li []list.Item
	for _, i := range items.ParseFunc() {
		li = append(li, i)
	}

	l := list.New(li, items.NewDelegate(), m.width, m.height)
	return &l
}

// SetBrowsing sets the state to Browsing
func (m *Model) SetBrowsing() {
	m.state = Browsing
}

// IsBrowsing returns whether or not the list state is Browsing.
func (m Model) Browsing() bool {
	return !m.Model.SettingFilter() || m.state == Browsing
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.input.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.resetInput()
			case tea.KeyEnter:
				val := m.input.Value()
				if m.editable {
					cmd = InsertItem(val)
					cmds = append(cmds, cmd)
				}
				m.resetInput()
			}
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.KeyMap.Filter):
				m.state = Input
			}
			if m.Browsing() && m.Selectable() {
				switch msg.Type {
				case tea.KeyEnter:
					if !m.Model.MultiSelectable() {
						m.Model.ToggleItem()
					}
					return m, tea.Quit
				}
			}
		}

		if m.editable {
			cmd = m.handleEditing(msg)
			cmds = append(cmds, cmd)
		}

		li, cmd := m.Model.Update(msg)
		m.Model = &li
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleEditing(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case InsertItemMsg:
		item := NewItem(msg.Value)
		cmd = m.InsertItem(m.Index()+1, item)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.InsertItem):
			if m.hasInput {
				m.SetShowTitle(false)
				m.SetHeight(m.Height() - 1)
				m.state = Input
				cmds = append(cmds, m.input.Focus())
			}
		case key.Matches(msg, m.KeyMap.RemoveItem):
			m.RemoveItem(m.Index())
		}
	}
	return tea.Batch(cmds...)
}

type InsertItemMsg struct {
	Value string
}

func InsertItem(val string) tea.Cmd {
	return func() tea.Msg {
		return InsertItemMsg{
			Value: val,
		}
	}
}

// ResetInput resets the current filtering state.
func (m *Model) ResetInput() {
	m.resetInput()
}

func (m *Model) resetInput() {
	if m.state == Browsing {
		return
	}

	m.state = Browsing
	m.SetShowTitle(true)
	m.input.Reset()
	m.input.Blur()
}

func (m *Model) View() string {
	var views []string

	if m.input.Focused() {
		in := m.input.View()
		views = append(views, in)
	}

	li := m.Model.View()
	views = append(views, li)

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (m *Model) Init() tea.Cmd {
	return nil
}
