package item

import (
	"github.com/charmbracelet/bubbles/list"
)

type Items struct {
	all         []*Item
	MultiSelect bool
}

func NewItems() Items {
	return Items{}
}

func (i *Items) Add(item *Item) *Items {
	i.all = append(i.all, item)
	return i
}

func (i *Items) Set(idx int, item *Item) {
	i.all[idx] = item
}

func (i *Items) Process() {
	var items []*Item
	idx := 0
	for _, item := range i.all {
		if i.MultiSelect {
			item.MultiSelect = true
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
	i.all = items
}

func (i *Items) All() []list.Item {
	var li []list.Item
	for _, item := range i.all {
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
	default:
		items = i.Visible()
	}
	return items
}

func (i Items) Visible() []list.Item {
	var items []list.Item
	level := 0
	for _, item := range i.all {
		if !item.IsHidden {
			items = append(items, item)
		}
		if item.HasList() && item.ListOpen {
			level++
			for _, sub := range i.GetItemList(item) {
				sub.IsHidden = false
				sub.SetLevel(level)
				items = append(items, sub)
			}
		}
	}
	return items
}

func (i Items) Selections() []*Item {
	var items []*Item
	for _, item := range i.all {
		if item.IsSelected {
			items = append(items, item)
		}
	}
	return items
}

func (i Items) GetItem(item list.Item) *Item {
	idx := item.(*Item).Index()
	return i.all[idx]
}

func (i Items) GetItemByIndex(idx int) *Item {
	var item *Item
	if idx < len(i.all) {
		item = i.all[idx]
	}
	return item
}

func (i *Items) ToggleSelectedItem(item list.Item) {
	li := item.(*Item).ToggleSelected()
	i.all[li.Index()] = li
}

func (i *Items) ToggleAllSelectedItems() {
	for _, item := range i.all {
		item.ToggleSelected()
	}
}

func (i *Items) OpenAllItemLists() {
	for _, item := range i.All() {
		li := item.(*Item)
		if li.HasList() {
			i.OpenItemList(item)
		}
	}
}

func (i *Items) OpenItemList(item list.Item) {
	li := item.(*Item)
	li.ListOpen = true
	i.Set(li.Index(), li)

	for _, sub := range i.GetItemList(item) {
		sub.Show()
		i.Set(sub.Index(), sub)
	}
}

func (i Items) GetItemList(item list.Item) []*Item {
	var items []*Item
	li := item.(*Item)
	if li.HasList() {
		t := len(li.List.all)
		items = i.all[li.idx+1 : li.idx+t+1]
	}
	return items
}
