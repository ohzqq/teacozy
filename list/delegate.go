package list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

// Item Delegate Interface
func (d *Items) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var (
		curItem *Item
	)

	sel := m.SelectedItem()
	if item, ok := sel.(*Item); ok {
		curItem = d.GetItemByIndex(item.Index())
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.InfoKey):
			if sel != nil {
				if curItem.HasMeta() {
					return ShowItemInfoCmd(curItem)
				}
			}
		case key.Matches(msg, key.ToggleItemList):
			switch {
			case curItem.HasChildren():
				return ToggleItemChildrenCmd(curItem)
			case curItem.IsSub():
				return ToggleItemChildrenCmd(curItem.Parent)
			}
		case key.Matches(msg, key.ToggleItem):
			m.CursorDown()
			if curItem.HasChildren() {
				return ToggleItemChildrenCmd(curItem)
			}
			if d.MultiSelect {
				return ToggleSelectedItemCmd(curItem)
			}
		}
	}
	return nil
}

func (d Items) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		iStyle  = &d.styles
		content string
		curItem *Item
	)

	switch i := listItem.(type) {
	case *Item:
		curItem = i
		content = i.Content()
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - iStyle.CurrentItem.GetPaddingLeft() - iStyle.CurrentItem.GetPaddingRight())
		content = padding.String(truncate.StringWithTail(content, textwidth, style.Ellipsis), textwidth)
	}

	var (
		isCurrent  = index == m.Index()
		isSelected = curItem.IsSelected
	)

	render := iStyle.NormalItem.Render

	if isCurrent {
		render = func(s string) string {
			return iStyle.CurrentItem.Copy().Margin(0, 1, 0, curItem.Level).Render(s)
		}
	} else if isSelected {
		render = func(s string) string {
			return iStyle.SelectedItem.Copy().Margin(0, 1, 0, curItem.Level).Render(s)
		}
	} else if curItem.IsSub() {
		render = func(s string) string {
			return iStyle.SubItem.Copy().Margin(0, 1, 0, curItem.Level).Render(s)
		}
	} else {
		render = func(s string) string {
			return iStyle.NormalItem.Copy().Margin(0, 1, 0, curItem.Level).Render(s)
		}
	}

	prefix := dash
	if d.MultiSelect {
		prefix = uncheck
		if isSelected {
			prefix = check
		}
	}

	if curItem.HasChildren() {
		prefix = openSub
		if curItem.ShowChildren {
			prefix = closeSub
		}
	}

	if d.ShowKeys {
		prefix = none
		key := curItem.Name()
		//key := fieldStyle.Key.Render(curItem.Name())
		content = fmt.Sprintf("%s: %s", key, content)
	}

	content = prefix + content

	fmt.Fprintf(w, render(content))
	//fmt.Fprintf(w, "%d: %s", curItem.id, render(title))
}

func (d Items) Height() int {
	return 1
}

func (d Items) Spacing() int {
	return 0
}
