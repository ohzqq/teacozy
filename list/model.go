package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/teacozy/util"
)

type Model struct {
	list.Model
}

type Item struct {
	Value string
}

func New(in []string) *Model {
	var items []list.Item
	for _, i := range in {
		items = append(items, Item{Value: i})
	}

	w, h := util.TermSize()

	m := list.New(items, list.NewDefaultDelegate(), w, h)
	m.SetNoLimit()

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

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.Model.View()
}

func (i Item) FilterValue() string { return i.Value }
func (i Item) Title() string       { return i.Value }
func (i Item) Description() string { return i.Value }
