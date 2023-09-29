package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/teacozy/util"
)

type Model struct {
	list.Model
}

type Items func() []*Item

func New(items Items) *Model {
	var li []list.Item
	for _, i := range items() {
		li = append(li, i)
	}
	w, h := util.TermSize()

	del := list.NewDefaultDelegate()
	//del.ShowDescription = true

	m := list.New(li, del, w, h)

	return &Model{
		Model: m,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.Model.SettingFilter() {
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
