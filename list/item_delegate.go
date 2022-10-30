package list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	cozykey "github.com/ohzqq/teacozy/key"
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
	keys        cozykey.KeyMap
	styles      ItemStyle
}

func NewItemDelegate(multi bool) itemDelegate {
	return itemDelegate{
		MultiSelect: multi,
		keys:        cozykey.DefaultKeys(),
		styles:      ItemStyles(),
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
	//var cur Item

	//switch i := m.SelectedItem().(type) {
	//case Item:
	//  cur = i
	//}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.ToggleItem):
			//if cur.HasList() {
			//return ToggleItemListCmd(cur)
			//return curItem.ShowListItemsCmd()
			//}
			return m.NewStatusMessage("item toggled")
			//return toggleItemCmd(cur)
		}
	}
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		iStyle  = &d.styles
		title   string
		curItem Item
	)

	switch i := listItem.(type) {
	case Item:
		curItem = i
		title = i.Content
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - iStyle.CurrentItem.GetPaddingLeft() - iStyle.CurrentItem.GetPaddingRight())
		title = padding.String(truncate.StringWithTail(title, textwidth, Ellipsis), textwidth)
	}

	var (
		isCurrent  = index == m.Index()
		isSelected = curItem.IsSelected
		isSub      = curItem.IsSub
	)

	render := iStyle.NormalItem.Render

	mark := curItem.Mark()
	if curItem.HasList && !curItem.ListIsOpen {
		//mark = itemListClosed.Mark()
		mark = "- "
	}

	if isCurrent {
		render = func(s string) string {
			return iStyle.CurrentItem.Copy().Margin(0, 1, 0, curItem.level).Render(mark + s)
		}
	} else if isSelected {
		render = func(s string) string {
			return iStyle.SelectedItem.Copy().Margin(0, 1, 0, curItem.level).Render(mark + s)
		}
	} else if isSub {
		render = func(s string) string {
			return iStyle.SubItem.Copy().Margin(0, 1, 0, curItem.level).Render(mark + s)
		}
	} else {
		render = func(s string) string {
			return iStyle.NormalItem.Copy().Margin(0, 1, 0, curItem.level).Render(mark + s)
		}
	}

	fmt.Fprintf(w, render(title))
	//fmt.Fprintf(w, "%d: %s", curItem.id, render(title))
}
