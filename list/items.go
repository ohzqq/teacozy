package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/ohzqq/teacozy/style"
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
	items       []*Item
	MultiSelect bool
	ShowKeys    bool
	styles      style.ItemStyle
}

func NewItems() *Items {
	return &Items{
		styles: style.ItemStyles(),
	}
}

func (i *Items) SetItems(items ...*Item) *Items {
	i.flat = items
	i.items = items
	return i
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
		fmt.Println(item.HasChildren())
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
