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

func NewItemDelegate(items *Items) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var (
			curItem *Item
		)

		sel := m.SelectedItem()
		if item, ok := sel.(*Item); ok {
			curItem = items.GetItemByIndex(item.Index())
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, key.ToggleItemList):
				var i *Item
				switch {
				case curItem.HasChildren():
					i = curItem
				case curItem.IsSub():
					i = curItem.Parent
				default:
					return nil
				}
				switch i.ShowChildren {
				case true:
					m.Select(i.Index())
					items.CloseItemList(i.Index())
				default:
					m.CursorDown()
					items.OpenItemList(i.Index())
				}
				m.SetItems(items.Visible())
			case key.Matches(msg, key.ToggleItem):
				m.CursorDown()
				if curItem.HasChildren() {
					switch curItem.ShowChildren {
					case true:
						m.Select(curItem.Index())
						items.CloseItemList(curItem.Index())
					default:
						items.OpenItemList(curItem.Index())
					}
					m.SetItems(items.Visible())
					return nil
				}
				if items.MultiSelect {
					items.ToggleSelectedItem(curItem.Index())
					return nil
				}
			}
		}
		return nil
	}
	d.ShowDescription = false
	d.SetSpacing(0)
	d.Styles = style.NewDefaultItemStyles()

	return d
}

func (d Items) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		content string
		curItem *Item
	)

	switch i := listItem.(type) {
	case *Item:
		curItem = i
		content = i.Content()
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - d.Style.Current.GetPaddingLeft() - d.Style.Current.GetPaddingRight())
		content = padding.String(truncate.StringWithTail(content, textwidth, style.Ellipsis), textwidth)
	}

	var (
		isCurrent = index == m.Index()
	)

	itemStyle := d.Style.Normal
	switch {
	case isCurrent:
		itemStyle = d.Style.Current
	case curItem.IsSelected:
		itemStyle = d.Style.Selected
	case curItem.IsSub():
		itemStyle = d.Style.Sub
	}
	itemStyle = itemStyle.Copy().Margin(0, 1, 0, curItem.Level)

	fmt.Fprintf(w, itemStyle.Render(curItem.Prefix()+content))
}
