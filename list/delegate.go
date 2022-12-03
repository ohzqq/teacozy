package list

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
	"github.com/ohzqq/teacozy/style"
)

func NewItemDelegate(items *Items) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = items.UpdateItem
	d.ShowDescription = false
	d.SetSpacing(0)
	d.Styles = style.NewDefaultItemStyles()

	return d
}

func (d Items) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		content string
		desc    string
		curItem *Item
	)

	switch i := listItem.(type) {
	case *Item:
		curItem = i
		content = i.Title()
		desc = i.Description()
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - d.Style.Current.GetPaddingLeft() - d.Style.Current.GetPaddingRight())
		content = padding.String(truncate.StringWithTail(content, textwidth, style.Ellipsis), textwidth)

		if d.ShowDescription {
			var lines []string
			for i, line := range strings.Split(desc, "\n") {
				if i >= d.Height()-1 {
					break
				}
				lines = append(lines, truncate.StringWithTail(line, textwidth, style.Ellipsis))
			}
			desc = strings.Join(lines, "\n")
		}
	}

	var (
		isCurrent   = index == m.Index()
		prefix      = curItem.Prefix()
		prefixWidth = lipgloss.Width(prefix)
	)

	switch {
	case isCurrent:
		s := d.Style.Current.Copy().Reverse(true)
		prefix = s.Render(prefix)
	case curItem.IsSelected:
		s := d.Style.Selected
		prefix = s.Render(prefix)
		content = s.Render(content)
		desc = s.Render(desc)
	case curItem.IsSub():
		s := d.Style.Sub
		prefix = s.Render(prefix)
	default:
		s := d.Style.Normal
		prefix = d.Style.Current.Render(prefix)
		content = s.Render(content)
		desc = s.Render(desc)
	}
	prefix = indent.String(prefix, uint(curItem.Level))

	if d.ShowDescription && !curItem.HasChildren() {
		desc = indent.String(desc, uint(curItem.Level+prefixWidth))
		fmt.Fprintf(w, "%s%s\n%s", prefix, content, desc)
		return
	}

	fmt.Fprintf(w, "%s%s", prefix, content)
}

func ToggleItem(curItem *Item, items *Items, m *list.Model) tea.Cmd {
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
		m.CursorDown()
		items.ToggleSelectedItem(curItem.Index())
		return nil
	}
	return nil
}

func ToggleItemList(curItem *Item, items *Items, m *list.Model) tea.Cmd {
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

	return nil
}
