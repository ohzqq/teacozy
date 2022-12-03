package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

const (
	check    string = "[x]"
	uncheck  string = "[ ]"
	dash     string = " > "
	openSub  string = `[+]`
	closeSub string = `[-]`
	none     string = ``
)

type Items struct {
	flat        []*Item
	items       []*Item
	MultiSelect bool
	ShowKeys    bool
	Style       style.ItemStyle
	UpdateFuncs map[*key.Key]ItemUpdateFunc
	list.DefaultDelegate
}

type ItemUpdateFunc func(*Item, *Items, *list.Model) tea.Cmd

func DefaultItemFuncs() map[*key.Key]ItemUpdateFunc {
	x := key.NewKey("x", "toggle item list")
	t := key.NewKey(" ", "toggle")

	funcs := make(map[*key.Key]ItemUpdateFunc)
	funcs[x] = ToggleItemList
	funcs[t] = ToggleItem

	return funcs
}

func NewItems() *Items {
	items := &Items{
		Style:       style.ItemStyles(),
		UpdateFuncs: DefaultItemFuncs(),
	}
	items.DefaultDelegate = NewItemDelegate(items)

	return items
}

func (i *Items) SetItems(items ...*Item) *Items {
	i.flat = items
	i.items = items
	return i
}

func (i *Items) SetUpdateFuncs(funcs map[*key.Key]ItemUpdateFunc) *Items {
	i.UpdateFuncs = funcs
	return i
}

func (i *Items) AddUpdateFunc(k, h string, fn ItemUpdateFunc) *Items {
	kp := key.NewKey(k, h)
	i.UpdateFuncs[kp] = fn
	return i
}

func (i *Items) UpdateItem(msg tea.Msg, m *list.Model) tea.Cmd {
	var (
		curItem *Item
	)

	sel := m.SelectedItem()
	if item, ok := sel.(*Item); ok {
		curItem = i.GetItemByIndex(item.Index())
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for k, fn := range i.UpdateFuncs {
			if k.Matches(msg) {
				return fn(curItem, i, m)
			}
		}
	}
	return nil
}

func (i *Items) Process() {
	var items []*Item
	idx := 0
	for _, item := range i.items {
		if i.MultiSelect {
			item.SetMultiSelect()
		}
		item.idx = idx
		items = append(items, item)
		for _, sub := range item.Flatten() {
			sub.Parent = item
			idx++
			sub.idx = idx
			items = append(items, sub)
		}
		idx++
	}
	i.flat = items
}

func (i Items) Flat() []*Item {
	return i.flat
}

func (i Items) All() []*Item {
	return i.items
}

func (i *Items) AllItems() []list.Item {
	var li []list.Item
	for _, item := range i.flat {
		li = append(li, item)
	}
	return li
}

func (i *Items) Add(item *Item) *Items {
	i.flat = append(i.flat, item)
	i.items = append(i.items, item)
	return i
}

func (i *Items) Set(idx int, item *Item) {
	i.flat[idx] = item
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

func (i Items) Visible() []list.Item {
	var items []list.Item
	level := 0
	for _, item := range i.Flat() {
		if !item.IsHidden {
			items = append(items, item)
		}
		if item.HasChildren() && item.ShowChildren {
			level++
			for _, sub := range i.GetItemList(item) {
				sub.Parent = item
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
	for _, item := range i.Flat() {
		if item.IsSelected {
			items = append(items, item)
		}
	}
	return items
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

func (i *Items) SelectAllItems() {
	for _, item := range i.flat {
		item.Select()
	}
}

func (i *Items) DeselectAllItems() {
	for _, item := range i.flat {
		item.Deselect()
	}
}
func (i *Items) OpenAllItemLists() {
	for _, item := range i.AllItems() {
		li := item.(*Item)
		if li.HasChildren() {
			i.OpenItemList(li.Index())
		}
	}
}

func (i *Items) OpenItemList(idx int) {
	li := i.GetItemByIndex(idx)
	li.ShowChildren = true
	i.Set(li.Index(), li)

	for _, sub := range i.GetItemList(li) {
		sub.Show()
		i.Set(sub.Index(), sub)
	}
}

func (i *Items) CloseItemList(idx int) {
	li := i.GetItemByIndex(idx)
	li.ShowChildren = false
	i.Set(li.Index(), li)

	for _, sub := range i.GetItemList(li) {
		sub.Hide()
		i.Set(sub.Index(), sub)
		if sub.HasChildren() {
			i.CloseItemList(sub.Index())
		}
	}
}

func (i Items) GetItemList(item list.Item) []*Item {
	var items []*Item
	li := item.(*Item)
	if li.HasChildren() {
		t := li.TotalChildren()
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
