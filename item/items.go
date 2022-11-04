package item

import (
	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/exp/slices"
)

type Items struct {
	AllItems    []*Item
	MultiSelect bool
}

func NewItems() Items {
	return Items{}
}

func (i *Items) Add(item *Item) *Items {
	i.AllItems = append(i.AllItems, item)
	return i
}

func (i *Items) Set(idx int, item *Item) {
	i.AllItems[idx] = item
}

func (i *Items) All() []list.Item {
	var items []*Item
	var li []list.Item
	idx := 0
	for _, item := range i.AllItems {
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
	i.AllItems = items
	return li
}

func (i Items) Visible() []list.Item {
	var items []list.Item
	for _, item := range i.AllItems {
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
	return i.AllItems[idx]
}

func (i Items) GetItemByIndex(idx int) *Item {
	var item *Item
	if idx < len(i.AllItems) {
		item = i.AllItems[idx]
	}
	return item
}

func (i Items) ToggleSelected(item list.Item) {
	li := item.(*Item).ToggleSelected()
	i.AllItems[li.Index()] = li
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
		t := len(li.List.AllItems)
		items = i.AllItems[li.idx+1 : li.idx+t+1]
	}
	return items
}
