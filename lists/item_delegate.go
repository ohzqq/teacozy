package lists

import (
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	cozykey "github.com/ohzqq/teacozy/key"
)

type itemDelegate struct {
	IsMultiSelect bool
	keys          cozykey.KeyMap
	styles        ItemStyle
}

func NewItemDelegate(multi bool) itemDelegate {
	return itemDelegate{
		IsMultiSelect: multi,
		keys:          cozykey.DefaultKeys(),
		styles:        ItemStyles(),
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
		cur  Item
		cmds []tea.Cmd
	)

	switch i := m.SelectedItem().(type) {
	case Item:
		cur = i
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, cozykey.EditField):
			cmds = append(cmds, EditItemCmd())
		case key.Matches(msg, d.keys.ToggleItem):
			//if cur.HasList() {
			//return ToggleItemListCmd(cur)
			//return curItem.ShowListItemsCmd()
			//}
			//cur.state = itemSelected
			msg := fmt.Sprintf("%v", cur.Content)
			cmds = append(cmds, m.NewStatusMessage(msg))
			cmds = append(cmds, ToggleItemCmd())
		}
	}
	return tea.Batch(cmds...)
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
		title = strconv.Itoa(i.Idx) + i.Content
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

	prefix := curItem.Prefix()
	//if curItem.HasList() && !curItem.ListIsOpen {
	//  prefix = itemListClosed.Prefix()
	//}

	//if !d.IsMultiSelect {
	//  prefix = dash
	//}

	if isCurrent {
		render = func(s string) string {
			return iStyle.CurrentItem.Copy().Margin(0, 1, 0, curItem.level).Render(prefix + s)
		}
	} else if isSelected {
		render = func(s string) string {
			return iStyle.SelectedItem.Copy().Margin(0, 1, 0, curItem.level).Render(prefix + s)
		}
	} else if isSub {
		render = func(s string) string {
			return iStyle.SubItem.Copy().Margin(0, 1, 0, curItem.level).Render(prefix + s)
		}
	} else {
		render = func(s string) string {
			return iStyle.NormalItem.Copy().Margin(0, 1, 0, curItem.level).Render(prefix + s)
		}
	}

	fmt.Fprintf(w, render(title))
	//fmt.Fprintf(w, "%d: %s", curItem.id, render(title))
}
