package list

import "github.com/ohzqq/teacozy/style"

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

func (i Items) Flat() []*Item {
	return i.flat
}

func (i Items) All() []*Item {
	return i.items
}
