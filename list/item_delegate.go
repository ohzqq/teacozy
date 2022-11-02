package list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

const (
	check    string = "[x] "
	uncheck  string = "[ ] "
	dash     string = "- "
	openSub  string = `[+] `
	closeSub string = `[-] `
)

type itemDelegate struct {
	MultiSelect bool
	keys        urkey.KeyMap
	styles      style.ItemStyle
}

func NewItemDelegate(multi bool) itemDelegate {
	return itemDelegate{
		MultiSelect: multi,
		styles:      style.ItemStyles(),
		keys:        urkey.DefaultKeys(),
	}
}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) ShortHelp() []key.Binding {
	return d.keys.ShortHelp()
}

func (d itemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{d.keys.Enter},
	}
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var (
		curItem Item
		cmds    []tea.Cmd
	)

	switch i := m.SelectedItem().(type) {
	case Item:
		curItem = i
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, urkey.EditField):
			cmds = append(cmds, EditItemCmd())
		case key.Matches(msg, urkey.ToggleItem):
			if curItem.HasList() {
				return ToggleItemListCmd(curItem)
			}
			return toggleItemCmd(curItem)
		}
	}
	return tea.Batch(cmds...)
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		iStyle  = &d.styles
		content string
		curItem Item
	)

	switch i := listItem.(type) {
	case Item:
		curItem = i
		c := fmt.Sprintf("%d: %s", i.id, i.FilterValue())
		content = c
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - iStyle.CurrentItem.GetPaddingLeft() - iStyle.CurrentItem.GetPaddingRight())
		content = padding.String(truncate.StringWithTail(content, textwidth, style.Ellipsis), textwidth)
	}

	var (
		isCurrent  = index == m.Index()
		isSelected = curItem.IsSelected()
		isSub      = curItem.IsSub()
	)

	render := iStyle.NormalItem.Render

	prefix := curItem.Prefix()
	//if curItem.HasList() && !curItem.ListIsOpen() {
	//  prefix = ItemListClosed.Prefix()
	//}

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
