package list

import (
	"github.com/ohzqq/bubbles/list"
)

type Items struct {
	ParseFunc func() []*Item
	ListType  list.ListType
}

type ParseItems func() []*Item

type Item struct {
	title       string
	desc        string
	filterValue string
}

type ItemOpt func(*Item)
type ItemsOpt func(*Items)

func NewItems(fn ParseItems, opts ...ItemsOpt) Items {
	items := Items{
		ParseFunc: fn,
	}
	for _, opt := range opts {
		opt(&items)
	}
	return items
}

func (items Items) NewDelegate() list.DefaultDelegate {
	del := list.NewDefaultDelegate()
	del.SetListType(items.ListType)
	return del
}

func OrderedList() ItemsOpt {
	return func(items *Items) {
		items.ListType = list.Ol
	}
}

func UnorderedList() ItemsOpt {
	return func(items *Items) {
		items.ListType = list.Ul
	}
}

func ItemsStringSlice(items []string) ParseItems {
	fn := func() []*Item {
		var li []*Item
		for _, item := range items {
			li = append(li, NewItem(item))
		}
		return li
	}
	return fn
}

func ItemsMapSlice(items []map[string]string) ParseItems {
	fn := func() []*Item {
		var li []*Item
		for _, item := range items {
			for k, v := range item {
				li = append(li, NewItem(k, Description(v)))
			}
		}
		return li
	}
	return fn
}

func ItemsMap(items map[string]string) ParseItems {
	fn := func() []*Item {
		var li []*Item
		for k, v := range items {
			li = append(li, NewItem(k, Description(v)))
		}
		return li
	}
	return fn
}

func NewItem(title string, opts ...ItemOpt) *Item {
	item := &Item{
		title:       title,
		desc:        title,
		filterValue: title,
	}

	for _, opt := range opts {
		opt(item)
	}

	return item
}

func Description(desc string) ItemOpt {
	return func(i *Item) {
		i.desc = desc
	}
}

func FilterValue(val string) ItemOpt {
	return func(i *Item) {
		i.filterValue = val
	}
}

func (i Item) FilterValue() string { return i.filterValue }
func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
