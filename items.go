package teacozy

import (
	"github.com/charmbracelet/bubbles/list"
)

type Items struct {
	items       []*Item
	Delegate    *ItemDelegate
	MultiSelect bool
	ShowKeys    bool
}

func NewItems() Items {
	return Items{
		Delegate: NewItemDelegate(),
	}
}

func (i *Items) SetItems(items ...*Item) *Items {
	i.items = items
	return i
}

func (i *Items) List() list.Model {
	i.Process()
	w, h := TermSize()
	l := list.New(i.Visible(), i.Delegate, w, h)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.KeyMap = ListKeyMap()
	l.Styles = ListStyles()
	return l
}

func (i *Items) Add(item *Item) *Items {
	i.items = append(i.items, item)
	return i
}

func (i *Items) Set(idx int, item *Item) {
	i.items[idx] = item
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
			idx++
			sub.idx = idx
			items = append(items, sub)
		}
		idx++
	}
	i.items = items
}

func (i Items) All() []*Item {
	return i.items
}

func (i *Items) AllItems() []list.Item {
	var li []list.Item
	for _, item := range i.items {
		li = append(li, item)
	}
	return li
}

func (i Items) Display(opt string) []list.Item {
	var items []list.Item
	switch opt {
	case "selected":
		for _, item := range i.Selections() {
			items = append(items, item)
		}
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

func (i Items) Selections() []*Item {
	var items []*Item
	for _, item := range i.items {
		if item.IsSelected {
			items = append(items, item)
		}
	}
	return items
}

func (i Items) Get(item list.Item) *Item {
	idx := item.(*Item).Index()
	return i.items[idx]
}

func (i Items) GetItemByIndex(idx int) *Item {
	var item *Item
	if idx < len(i.items) {
		item = i.items[idx]
	}
	return item
}

func (i *Items) ToggleSelectedItem(idx int) {
	li := i.GetItemByIndex(idx).ToggleSelected()
	i.items[li.Index()] = li
}

func (i *Items) ToggleAllSelectedItems() {
	for _, item := range i.items {
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
		t := len(li.List.items)
		items = i.items[li.idx+1 : li.idx+t+1]
	}
	return items
}
