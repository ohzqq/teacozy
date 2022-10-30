package list

import "github.com/charmbracelet/bubbles/list"

type Items struct {
	All      []list.Item
	Selected map[int]list.Item
	Visible  []list.Item
}

func (i Items) NewList(title string) List {
	return New(title, i)
}

func (i *Items) Add(item Item) {
	i.appendItem(item)
	if item.HasList {
		for _, l := range item.Items.All {
			li := l.(Item)
			i.Add(li)
		}
	}
}

func (i *Items) appendItem(item Item) {
	item.idx = len(i.All)
	i.All = append(i.All, item)
}

type Item struct {
	idx           int
	Content       string
	Items         Items
	level         int
	IsVisible     bool
	IsSelected    bool
	IsSub         bool
	HasList       bool
	ListIsOpen    bool
	IsMultiSelect bool
}

func NewItem(content string) Item {
	return Item{Content: content}
}

func (i Item) FilterValue() string {
	return i.Content
}

func (i Item) Mark() string {
	return "- "
}
