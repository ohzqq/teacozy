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

type Items struct {
	flat        []*Item
	MultiSelect bool
	ShowKeys    bool
	styles      ItemStyle
}

func NewItems() *Items {
	return &Items{
		styles: ItemStyles(),
	}
}

func (i *Items) SetItems(items ...*Item) *Items {
	i.flat = items
	return i
}

func (i *Items) List() list.Model {
	i.Process()
	w, h := TermSize()
	l := NewListModel(w, h, i)
	return l
}

func (i *Items) Add(item *Item) *Items {
	i.flat = append(i.flat, item)
	return i
}

func (i *Items) Set(idx int, item *Item) {
	i.flat[idx] = item
}

func (i *Items) Process() {
	var items []*Item
	idx := 0
	for _, item := range i.All() {
		if i.MultiSelect {
			item.SetMultiSelect()
		}
		if item.HasFields() {
			item.hasFields = true
		}
		item.idx = idx
		items = append(items, item)
		for _, sub := range item.Flatten() {
			idx++
			sub.idx = idx
			items = append(items, sub)
		}
		idx++
	}
	i.flat = items
}

func (i Items) All() []*Item {
	return i.flat
}

func (i *Items) AllItems() []list.Item {
	var li []list.Item
	for _, item := range i.flat {
		li = append(li, item)
	}
	return li
}

func (i Items) Display(opt string) []list.Item {
	var items []list.Item
	switch opt {
	case "selected":
		items = i.Selections()
		//for _, item := range i.Selections() {
		//  items = append(items, item)
		//}

	case "all":
		items = i.AllItems()
	default:
		items = i.Visible()
	}
	return items
}

func (i Items) Visible() []list.Item {
	var items []list.Item
	level := 0
	for _, item := range i.All() {
		if !item.IsHidden {
			items = append(items, item)
		}
		if item.HasList() && item.ListOpen {
			level++
			for _, sub := range i.GetItemList(item) {
				sub.Hide()
				sub.SetLevel(level)
				items = append(items, sub)
			}
		}
	}
	return items
}

func (i Items) Selections() []list.Item {
	var items []list.Item
	for _, item := range i.All() {
		if item.IsSelected {
			items = append(items, item)
		}
	}
	return items
}

func (i Items) Get(item list.Item) *Item {
	idx := item.(*Item).Index()
	return i.flat[idx]
}

func (i Items) GetItemByIndex(idx int) *Item {
	var item *Item
	if idx < len(i.flat) {
		item = i.flat[idx]
	}
	return item
}

func (i *Items) ToggleSelectedItem(idx int) {
	li := i.GetItemByIndex(idx).ToggleSelected()
	i.flat[li.Index()] = li
}

func (i *Items) ToggleAllSelectedItems() {
	for _, item := range i.flat {
		item.ToggleSelected()
	}
}

func (i *Items) OpenAllItemLists() {
	for _, item := range i.AllItems() {
		li := item.(*Item)
		if li.HasList() {
			i.OpenItemList(li.Index())
		}
	}
}

func (i *Items) OpenItemList(idx int) {
	li := i.GetItemByIndex(idx)
	li.ListOpen = true
	i.Set(li.Index(), li)

	for _, sub := range i.GetItemList(li) {
		sub.Show()
		i.Set(sub.Index(), sub)
	}
}

func (i *Items) CloseItemList(idx int) {
	li := i.GetItemByIndex(idx)
	li.ListOpen = false
	i.Set(li.Index(), li)

	for _, sub := range i.GetItemList(li) {
		sub.Hide()
		i.Set(sub.Index(), sub)
		if sub.HasList() {
			i.CloseItemList(sub.Index())
		}
	}
}

func (i Items) GetItemList(item list.Item) []*Item {
	var items []*Item
	li := item.(*Item)
	if li.HasList() {
		t := len(li.List.flat)
		items = i.flat[li.idx+1 : li.idx+t+1]
	}
	return items
}

func (d *Items) SetShowKeys() *Items {
	d.ShowKeys = true
	return d
}

func (d *Items) SetMultiSelect() *Items {
	d.MultiSelect = true
	return d
}

func (d Items) Height() int {
	return 1
}

func (d Items) Spacing() int {
	return 0
}

func (d *Items) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var (
		curItem *Item
		cmds    []tea.Cmd
	)

	sel := m.SelectedItem()
	if item, ok := sel.(*Item); ok {
		curItem = d.GetItemByIndex(item.Index())
	}
	//switch i := m.SelectedItem().(type) {
	//case *Item:
	//  curItem = i
	//}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Info):
			if sel != nil {
				if curItem.HasFields() {
					cmds = append(cmds, ShowItemInfoCmd(curItem))
				}
			}
		case key.Matches(msg, Keys.EditField):
			cmds = append(cmds, EditFormItemCmd(curItem))
		case key.Matches(msg, Keys.ToggleItem):
			m.CursorDown()
			if curItem.HasList() {
				return ToggleItemListCmd(curItem)
			}
			if d.MultiSelect {
				d.ToggleSelectedItem(curItem.Index())
				//return ToggleSelectedItemCmd(curItem)
			}
		}
	}
	return tea.Batch(cmds...)
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
