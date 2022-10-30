package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	cozykey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	Model            list.Model
	Items            Items
	Selections       []Item
	Keys             cozykey.KeyMap
	Title            string
	ShowSelectedOnly bool
	FocusedView      string
	IsMultiSelect    bool
	width            int
	height           int
	ShowMenu         bool
	Action           ListAction
}

func New(title string, items Items, multi bool) List {
	m := List{Items: items, Keys: cozykey.DefaultKeys()}
	m.Model = list.New(items.All, NewItemDelegate(multi), m.Width(), m.Height())
	m.Model.Title = title
	m.Model.Styles = ListStyles()
	m.Model.SetShowStatusBar(false)
	m.Model.SetShowHelp(false)
	return m
}

func (l List) Width() int {
	return util.TermWidth()
}

func (l List) Height() int {
	return util.TermHeight()
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.IsMultiSelect {
		} else {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cmds = append(cmds, m.Action(m))
			}
		}
		switch {
		case key.Matches(msg, m.Keys.ExitScreen):
			cmds = append(cmds, tea.Quit)
		}
	case ToggleItemMsg:
		cur := m.Model.SelectedItem().(Item)
		if m.IsMultiSelect {
			cur.IsSelected = !cur.IsSelected
		}
		m.SetItem(m.Model.Index(), cur)
	case UpdateVisibleItemsMsg:
		switch string(msg) {
		case "selected":
		}
	}
	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *List) SetItem(modelIndex int, item Item) {
	m.Model.SetItem(modelIndex, item)
	m.Items.All[item.Idx] = item
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	return m.Model.View()
}

//func (l List) GetHeight(items []list.Item) int {
//  max := util.TermHeight()
//  total := len(items)
//  cur := l.Model.Height()

//  switch {
//  case l.isFullScreen:
//    return max
//  case cur > max:
//    return max
//  case total < max:
//    return total + 6
//  default:
//    return max
//  }
//}
