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
		textwidth := uint(m.Width() - iStyle.Current.GetPaddingLeft() - iStyle.Current.GetPaddingRight())
		content = padding.String(truncate.StringWithTail(content, textwidth, style.Ellipsis), textwidth)
	}

	var (
		isCurrent = index == m.Index()
	)

	s := curItem.Style(isCurrent)
	fmt.Fprintf(w, s.Render(curItem.Prefix()+content))
}

func (d Items) Height() int {
	return 1
}

func (d Items) Spacing() int {
	return 0
}
