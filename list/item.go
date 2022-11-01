package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
)

type ItemState int

//go:generate go run golang.org/x/tools/cmd/stringer -type=ItemState
const (
	ItemNotSelected ItemState = iota + 1
	ItemSelected
	ItemListOpen
	ItemListClosed
)

func (s ItemState) Mark() string {
	switch s {
	case ItemNotSelected:
		return uncheck
	case ItemSelected:
		return check
	case ItemListOpen:
		return closeSub
	case ItemListClosed:
		return openSub
	default:
		return dash
	}
}

type Items []list.Item

func (li Items) Get(idx int) Item {
	if idx < len(li) {
		return li[idx].(Item)
	}
	return Item{}
}

func (li Items) GetSelected() Items {
	var items Items
	for _, item := range li {
		if i, ok := item.(Item); ok && i.IsSelected() {
			items = append(items, i)
		}
	}
	return items
}

func (li Items) SelectAll() Items {
	var items Items
	for _, i := range li {
		item := i.(Item)
		item.isSelected = true
		items = append(items, item)
	}
	return items
}

func (li Items) SetSub(level int) Items {
	var items Items
	for _, i := range li {
		item := i.(Item)
		item.SetLevel(level)
		items = append(items, item)
	}
	return items
}

func (li Items) OpenList(idx int) Item {
	item := li.Get(idx)
	item.state = ItemListOpen
	return item
}

func (li Items) CloseList(idx int) Item {
	item := li.Get(idx)
	item.state = ItemListClosed
	return item
}

func (li Items) ToggleList(idx int) Item {
	item := li.Get(idx)

	var i Item
	if item.HasList() {
		i = li.OpenList(idx)
		if item.ListIsOpen() {
			i = li.CloseList(idx)
		}
	}

	li[idx] = i

	return li.Get(idx)
}

func (li Items) Select(idx int) Item {
	item := li.Get(idx)
	item.state = ItemSelected
	return item
}

func (li Items) Deselect(idx int) Item {
	item := li.Get(idx)
	item.state = ItemNotSelected
	return item
}

func (li Items) Toggle(idx int) Item {
	item := li.Get(idx)

	i := li.Select(idx)

	if item.IsSelected() {
		i = li.Deselect(idx)
	}

	li[idx] = i

	return li.Get(idx)
}

func (li Items) NewList(title string, state listType) *List {
	l := New(title)
	l.state = state
	return l
}

type Item struct {
	defaultItem
	data       list.Item
	id         int
	isSelected bool
	hasList    bool
	isSub      bool
	listOpen   bool
	isVisible  bool
	mark       string
	isMulti    bool
	state      ItemState
	level      int
	items      Items
	list       *List
	Content    string
}

func (i *Item) SetContent(content string) {
	i.Content = content
}

type defaultItem struct {
	title       string
	filterValue string
}

func (i defaultItem) FilterValue() string {
	return i.filterValue
}

func (i defaultItem) Title() string {
	return i.title
}

type ListItem interface {
	FilterValue() string
	Title() string
}

func NewListItem(i ListItem) Item {
	return Item{
		data:        i,
		Content:     i.FilterValue(),
		defaultItem: newDefaultItem(i.Title(), i.FilterValue()),
	}
}

func newDefaultItem(title, fv string) defaultItem {
	return defaultItem{
		title:       title,
		filterValue: fv,
	}
}

func NewDefaultItem(title, fv string) Item {
	i := newDefaultItem(title, fv)
	return Item{
		data:        i,
		defaultItem: i,
	}
}

func (i Item) State() ItemState {
	if i.HasList() && i.state == ItemListClosed {
		return ItemListClosed
	}
	return i.state
}

func (i *Item) Edit() textarea.Model {
	input := textarea.New()
	input.SetValue(i.Content)
	input.ShowLineNumbers = false
	return input
}

func (i Item) Mark() string {
	return i.state.Mark()
}

func (i *Item) SetId(id int) *Item {
	i.id = id
	return i
}

func (i Item) Data() list.Item {
	return i.data
}

func (i Item) IsSelected() bool {
	if i.isSelected || i.state == ItemSelected {
		return true
	}
	return false
}

func (i Item) ListIsOpen() bool {
	return i.state == ItemListOpen
}

func (i Item) IsSub() bool {
	return i.isSub
}

func (i Item) HasList() bool {
	return i.hasList
}

func (i *Item) Toggle() {
	i.isSelected = !i.isSelected
	switch i.IsSelected() {
	case true:
		i.state = ItemNotSelected
	case false:
		i.state = ItemSelected
	}
}

func (i *Item) SetLevel(l int) *Item {
	i.level = l
	return i
}

func (i *Item) SetIsSub() *Item {
	i.isSub = true
	return i
}

func (i *Item) SetList(l *List) *Item {
	i.list = l
	return i
}
