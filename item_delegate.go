package teacozy

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"github.com/ohzqq/teacozy/style"
)

type itemDelegate struct {
	multiSelect bool
	showKeys    bool
	styles      style.ItemStyle
}

func NewItemDelegate() itemDelegate {
	return itemDelegate{
		styles: style.ItemStyles(),
	}
}

func (d *itemDelegate) ShowKeys() {
	d.showKeys = true
}

func (d *itemDelegate) MultiSelect() {
	d.multiSelect = true
}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var (
		curItem *Item
		cmds    []tea.Cmd
	)

	switch i := m.SelectedItem().(type) {
	case *Item:
		curItem = i
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Info):
			if info := curItem.Info; info.String() != "" {
				cmds = append(cmds, ShowInfoCmd(curItem))
			}
		case key.Matches(msg, Keys.EditField):
			cmds = append(cmds, EditContentCmd(curItem))
		case key.Matches(msg, Keys.ToggleItem):
			m.CursorDown()
			if curItem.HasList() {
				return ToggleListCmd(curItem)
			}
			if d.multiSelect {
				return ToggleSelectedCmd(curItem)
			}
		}
	}
	return tea.Batch(cmds...)
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		iStyle  = &d.styles
		content string
		curItem *Item
	)

	switch i := listItem.(type) {
	case *Item:
		curItem = i
		content = i.FilterValue()
		//  if d.showKeys() {
		//    content =
		//}
		//content = c
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - iStyle.CurrentItem.GetPaddingLeft() - iStyle.CurrentItem.GetPaddingRight())
		content = padding.String(truncate.StringWithTail(content, textwidth, style.Ellipsis), textwidth)
	}

	var (
		isCurrent  = index == m.Index()
		isSelected = curItem.IsSelected
		isSub      = curItem.IsSub()
	)

	render := iStyle.NormalItem.Render

	prefix := curItem.Prefix()

	if isCurrent {
		render = func(s string) string {
			return iStyle.CurrentItem.Copy().Margin(0, 1, 0, curItem.Level).Render(prefix + s)
		}
	} else if isSelected {
		render = func(s string) string {
			return iStyle.SelectedItem.Copy().Margin(0, 1, 0, curItem.Level).Render(prefix + s)
		}
	} else if isSub {
		render = func(s string) string {
			return iStyle.SubItem.Copy().Margin(0, 1, 0, curItem.Level).Render(prefix + s)
		}
	} else {
		render = func(s string) string {
			return iStyle.NormalItem.Copy().Margin(0, 1, 0, curItem.Level).Render(prefix + s)
		}
	}

	fmt.Fprintf(w, render(content))
	//fmt.Fprintf(w, "%d: %s", curItem.id, render(title))
}
