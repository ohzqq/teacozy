package item

import "github.com/charmbracelet/bubbles/list"

type Items struct {
	all         []Item
	MultiSelect bool
}

func NewItems() Items {
	return Items{}
}

func (i *Items) Add(item Item) *Items {
	i.all = append(i.all, item)
	return i
}

func (i Items) All() []Item {
	var items []Item
	for _, item := range i.all {
		if i.MultiSelect {
			item.MultiSelect = true
		}
		items = append(items, item)
		items = append(items, item.Flatten()...)
	}
	i.all = items
	return i.all
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
