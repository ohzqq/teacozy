package teacozy

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
)

const (
	check    string = "[x] "
	uncheck  string = "[ ] "
	dash     string = "- "
	openSub  string = `[+] `
	closeSub string = `[-] `
	none     string = ``
)

type ItemDelegate struct {
	MultiSelect bool
	ShowKeys    bool
	styles      ItemStyle
}

func NewItemDelegate() *ItemDelegate {
	return &ItemDelegate{
		styles: ItemStyles(),
	}
}

func (d *ItemDelegate) SetShowKeys() *ItemDelegate {
	d.ShowKeys = true
	return d
}

func (d *ItemDelegate) SetMultiSelect() *ItemDelegate {
	d.MultiSelect = true
	return d
}

func (d ItemDelegate) Height() int {
	return 1
}

func (d ItemDelegate) Spacing() int {
	return 0
}

func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
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
			if info := curItem.Fields; info.String() != "" {
				cmds = append(cmds, ShowItemInfoCmd(curItem))
			}
		case key.Matches(msg, Keys.EditField):
			cmds = append(cmds, EditFormItemCmd(curItem))
		case key.Matches(msg, Keys.ToggleItem):
			m.CursorDown()
			if curItem.HasList() {
				return ToggleItemListCmd(curItem)
			}
			if d.MultiSelect {
				return ToggleSelectedItemCmd(curItem)
			}
		}
	}
	return tea.Batch(cmds...)
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		iStyle  = &d.styles
		content string
		curItem *Item
	)

	switch i := listItem.(type) {
	case *Item:
		curItem = i
		content = i.Value()
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - iStyle.CurrentItem.GetPaddingLeft() - iStyle.CurrentItem.GetPaddingRight())
		content = padding.String(truncate.StringWithTail(content, textwidth, Ellipsis), textwidth)
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
	if curItem.HasList() {
		prefix = openSub
		if curItem.ListOpen {
			prefix = closeSub
		}
	}

	if d.ShowKeys {
		prefix = none
		key := fieldStyle.Key.Render(curItem.Key())
		content = fmt.Sprintf("%s: %s", key, content)
	}

	if curItem.Changed {
		content = "*" + content
	}

	content = prefix + content

	fmt.Fprintf(w, render(content))
	//fmt.Fprintf(w, "%d: %s", curItem.id, render(title))
}
