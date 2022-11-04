package item

import (
	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/exp/slices"
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

func (i Items) All() []list.Item {
	var items []*Item
	var li []list.Item
	idx := 0
	for _, item := range i.all {
		if i.MultiSelect {
			item.MultiSelect = true
		}
		item.idx = idx
		items = append(items, item)
		li = append(li, item)
		for _, sub := range item.Flatten() {
			idx++
			sub.idx = idx
			items = append(items, sub)
			li = append(li, sub)
		}
		idx++
	}
	i.all = items
	return li
}

func (i Items) Visible() []list.Item {
	var items []list.Item
	for _, item := range i.all {
		if !item.IsHidden {
			items = append(items, item)
		}
	}
	return items
}

func (items *Items) GetItemIndex(i list.Item) int {
	content := i.FilterValue()
	fn := func(item list.Item) bool {
		c := item.FilterValue()
		return content == c
	}
	return slices.IndexFunc(items.All(), fn)
}

func (i Items) GetItem(item list.Item) *Item {
	idx := i.GetItemIndex(item)
	return i.all[idx]
}

func (i Items) GetItemByIndex(idx int) *Item {
	var item *Item
	if idx < len(i.all) {
		item = i.all[idx]
	}
	return item
}

func (i Items) ToggleSelected(item list.Item) {
	li := item.(*Item)
	li.ToggleSelected()
	i.all[li.Index()] = li
}
