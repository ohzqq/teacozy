package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	cozykey "github.com/ohzqq/teacozy/key"
)

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.area.Focused() {
			if key.Matches(msg, cozykey.SaveAndExit) {
				cur := m.Model.SelectedItem().(Item)
				val := m.area.Value()
				cur.SetContent(val)
				m.SetItem(m.Model.Index(), cur)
				m.area.Blur()
			}
			m.area, cmd = m.area.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			if m.IsMultiSelect {
			} else {
				switch {
				case key.Matches(msg, cozykey.Enter):
					cmds = append(cmds, m.Action(m))
				}
			}
			switch {
			case key.Matches(msg, m.Keys.ExitScreen):
				cmds = append(cmds, tea.Quit)
			default:
				for label, widget := range m.Widgets {
					if key.Matches(msg, widget.Toggle()) {
						widget.Focus()
						m.ShowWidget()
						cmds = append(cmds, SetFocusedViewCmd(label))
					}
				}
			}
			m.Model, cmd = m.Model.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ToggleItemMsg:
		cur := m.Model.SelectedItem().(Item)
		cur.ToggleSelected()
		//if m.IsMultiSelect {
		//  cur.IsSelected = !cur.IsSelected
		//}
		m.SetItem(m.Model.Index(), cur)
	case SetFocusedViewMsg:
		m.FocusedView = string(msg)
		if m.FocusedView == "list" && m.CurrentWidget() != nil {
			m.HideWidget()
		}
	case EditItemMsg:
		cur := m.Model.SelectedItem().(Item)
		m.area = cur.Edit()
		m.area.Focus()
	case UpdateVisibleItemsMsg:
		switch string(msg) {
		case "selected":
		}
	}

	switch focus := m.FocusedView; focus {
	case "list":
		switch msg := msg.(type) {
		//case UpdateDisplayedItemsMsg:
		//items := m.DisplayItems(string(msg))
		//m.Model.SetHeight(m.GetHeight(items))
		//cmds = append(cmds, m.Model.SetItems(items))
		case UpdateWidgetContentMsg:
			m.CurrentWidget().SetContent(string(msg))
			m.HideWidget()
		}

		//m.Model, cmd = m.Model.Update(msg)
		//cmds = append(cmds, cmd)
	default:
		if m.CurrentWidget() != nil {
			cmds = append(cmds, m.CurrentWidget().Update(&m, msg))
		}
	}

	return m, tea.Batch(cmds...)
}

func MultiSelectAction(m List) func(List) tea.Cmd {
	fn := func(m List) tea.Cmd {
		for _, item := range m.Items.All {
			i := item.(Item)
			m.Items.Selected[i.Idx] = item
		}
		if m.Items.HasSelections() {
			return tea.Quit
		}
		return nil
	}
	return fn
}

func SingleSelectAction(m List) func(List) tea.Cmd {
	fn := func(m List) tea.Cmd {
		cur := m.Model.SelectedItem()
		curItem := cur.(Item)
		curItem.IsSelected = true
		m.Items.Selected[curItem.Idx] = cur
		if m.Items.HasSelections() {
			return tea.Quit
		}
		return nil
	}
	return fn
}

type UpdateWidgetContentMsg string

func UpdateWidgetContentCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return UpdateWidgetContentMsg(s)
	}
}
func (m List) CurrentWidget() Widget {
	for _, w := range m.Widgets {
		if w.Focused() {
			return w
		}
	}
	return nil
}

func (m *List) HideWidget() {
	m.focusWidget = false
	if m.CurrentWidget() != nil {
		m.CurrentWidget().Blur()
	}
}

func (m *List) ShowWidget() {
	m.focusWidget = true
	if m.CurrentWidget() != nil {
		m.CurrentWidget().Focus()
	}
}
