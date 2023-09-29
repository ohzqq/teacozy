package list

import (
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
	m.SetLimit(0)
	m.SetFilteringEnabled(false)

	return &Model{
		Model:      m,
		width:      w,
		height:     h,
		items:      items,
		li:         li,
		state:      Browsing,
		selectable: m.Limit() != 0,
	}
}

func (m Model) NewTextinputModel() textinput.Model {
	input := textinput.New()
	input.PromptStyle = m.Styles.FilterPrompt
	input.Cursor.Style = m.Styles.FilterCursor
	return input
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.Model.SettingFilter() && m.selectable {
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

func (m *Model) UpdateKeys(state State) {
	switch state {
	case Input:
		m.KeyMap.CursorUp.SetEnabled(false)
		m.KeyMap.CursorDown.SetEnabled(false)
		m.KeyMap.NextPage.SetEnabled(false)
		m.KeyMap.PrevPage.SetEnabled(false)
		m.KeyMap.GoToStart.SetEnabled(false)
		m.KeyMap.GoToEnd.SetEnabled(false)
		m.KeyMap.Filter.SetEnabled(false)
		m.KeyMap.ClearFilter.SetEnabled(false)
		m.KeyMap.CancelWhileFiltering.SetEnabled(true)
		m.KeyMap.AcceptWhileFiltering.SetEnabled(m.FilterValue() != "")
		m.KeyMap.Quit.SetEnabled(false)
		m.KeyMap.ShowFullHelp.SetEnabled(false)
		m.KeyMap.CloseFullHelp.SetEnabled(false)

	default:
		hasItems := len(m.Items()) != 0
		m.KeyMap.CursorUp.SetEnabled(hasItems)
		m.KeyMap.CursorDown.SetEnabled(hasItems)

		hasPages := m.Paginator.TotalPages > 1
		m.KeyMap.NextPage.SetEnabled(hasPages)
		m.KeyMap.PrevPage.SetEnabled(hasPages)
		m.KeyMap.GoToStart.SetEnabled(hasItems)
		m.KeyMap.GoToEnd.SetEnabled(hasItems)
		m.KeyMap.Filter.SetEnabled(m.FilteringEnabled() && hasItems)
		m.KeyMap.ClearFilter.SetEnabled(m.IsFiltered())
		m.KeyMap.CancelWhileFiltering.SetEnabled(false)
		m.KeyMap.AcceptWhileFiltering.SetEnabled(false)
		m.KeyMap.ToggleItem.SetEnabled(true)
		//m.KeyMap.Quit.SetEnabled(!m.disableQuitKeybindings)

		if m.Help.ShowAll {
			m.KeyMap.ShowFullHelp.SetEnabled(true)
			m.KeyMap.CloseFullHelp.SetEnabled(true)
		} else {
			//minHelp := countEnabledBindings(m.FullHelp()) > 1
			//m.KeyMap.ShowFullHelp.SetEnabled(minHelp)
			//m.KeyMap.CloseFullHelp.SetEnabled(minHelp)
		}
	}
}
