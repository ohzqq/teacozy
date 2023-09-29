package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	list.Model
	items      Items
	li         []list.Item
	width      int
	height     int
	state      State
	selectable bool
}

type ListOpt func(*list.Model)

func New(items Items) *Model {
	var li []list.Item
	for _, i := range items() {
		li = append(li, i)
	}
	w, h := util.TermSize()

	del := list.NewDefaultDelegate()
	m := list.New(li, del, w, h)

	return &Model{
		Model:  m,
		width:  w,
		height: h,
		items:  items,
		li:     li,
		state:  Browsing,
	}
}

func (m Model) NewTextinputModel() textinput.Model {
	input := textinput.New()
	input.PromptStyle = m.Styles.FilterPrompt
	input.Cursor.Style = m.Styles.FilterCursor
	return input
}

func (m Model) Selectable() bool {
	return m.Limit() != 0
}

func (m Model) Input() bool {
	return m.Model.SettingFilter() || m.state == Input
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Filter):
			m.state = Input
		}
		if !m.Input() && m.Selectable() {
			switch msg.Type {
			case tea.KeyEnter:
				if !m.Model.MultiSelect() {
					m.Model.ToggleItem()
				}
				return m, tea.Quit
			}
		}
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.Model.View()
}
