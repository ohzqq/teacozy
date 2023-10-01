package list

import (
	"github.com/ohzqq/bubbles/list"
)

// Items holds the values to configure list.DefaultDelegate.
type Items struct {
	ParseFunc func() []*Item
	ListType  list.ListType
}

// ParseItems is a func to return a slice of Item.
type ParseItems func() []*Item

// Item is a struct to satisfy list.DefaultItem.
type Item struct {
	title       string
	desc        string
	filterValue string
}

// ItemOption sets options for Items.
type ItemOption func(*Items)

// NewItems initializes an Items.
func NewItems(fn ParseItems, opts ...ItemOption) Items {
	items := Items{
		ParseFunc: fn,
	}
	for _, opt := range opts {
		opt(&items)
	}
	return items
}

// NewDelegate returns a list.DefaultDelegate.
func (items Items) NewDelegate() list.DefaultDelegate {
	del := list.NewDefaultDelegate()
	del.SetListType(items.ListType)
	return del
}

// OrderedList sets the list.DefaultDelegate ListType.
func OrderedList() ItemOption {
	return func(items *Items) {
		items.ListType = list.Ol
	}
}

// UnrderedList sets the list.DefaultDelegate ListType.
func UnorderedList() ItemOption {
	return func(items *Items) {
		items.ListType = list.Ul
	}
}

// ItemsStringSlice returns a ParseItems for a slice of strings.
func ItemsStringSlice(items []string) ParseItems {
	fn := func() []*Item {
		li := make([]*Item, len(items))
		for i, item := range items {
			li[i] = NewItem(item)
		}
		return li
	}
	return fn
}

// ItemsMapSlice returns a ParseItems for a slice of map[string]string.
func ItemsMapSlice(items []map[string]string) ParseItems {
	fn := func() []*Item {
		li := make([]*Item, len(items))
		for i, item := range items {
			for k, v := range item {
				li[i] = NewItem(k).SetDescription(v)
			}
		}
		return li
	}
	return fn
}

// ItemsMap returns a ParseItems for a map[string]string.
func ItemsMap(items map[string]string) ParseItems {
	fn := func() []*Item {
		var li []*Item
		for k, v := range items {
			li = append(li, NewItem(k).SetDescription(v))
		}
		return li
	}
	return fn
}

// NewItem returns an Item struct.
func NewItem(title string) *Item {
	return &Item{
		title:       title,
		desc:        title,
		filterValue: title,
	}
}

// SetDescription sets the Item description.
func (i *Item) SetDescription(desc string) *Item {
	i.desc = desc
	return i
}

// SetFilterValue sets the Item filter value.
func (i *Item) SetFilterValue(val string) *Item {
	i.filterValue = val
	return i
}

func (i Item) FilterValue() string { return i.filterValue }
func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
