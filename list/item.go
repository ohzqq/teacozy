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

func (li Items) Display(opt string) Items {
	switch opt {
	case "selected":
		return li.GetSelected()
	default:
		var items Items
		level := 0
		for _, item := range li {
			i := item.(Item)
			if !i.IsHidden {
				items = append(items, i)
			}
			if i.HasList && i.ListIsOpen() {
				level++
				for _, sub := range li.GetSubList(i) {
					s := sub.(Item)
					s.SetIsSub().SetLevel(level)
					items = append(items, s)
				}
			}
		}
		return items
	}
}

func (li Items) GetSubList(i list.Item) Items {
	item := i.(Item)
	if item.HasList {
		t := len(item.Items)
		return li[item.id+1 : item.id+t+1]
	}
	return Items{}
}

func (i Items) Add(item Item) Items {
	i = i.appendItem(item)
	if item.HasList {
		for _, l := range item.Items {
			li := l.(Item)
			li.IsHidden = true
			i = i.Add(li)
		}
	}
	return i
}

func (i Items) appendItem(item Item) Items {
	item.SetId(len(i))
	//item.id = len(i)
	i = append(i, item)
	return i
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
	if item.HasList {
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

func (li Items) NewList(title string, state listType) *Model {
	l := New(title)
	l.state = state
	return l
}

type Item struct {
	data       list.Item
	id         int
	isSelected bool
	HasList    bool
	IsSub      bool
	listOpen   bool
	IsVisible  bool
	IsHidden   bool
	mark       string
	IsMulti    bool
	state      ItemState
	Level      int
	Items      Items
	List       *Model
	Content    string
}

func NewItem(item list.Item) Item {
	return Item{
		data:    item,
		Content: item.FilterValue(),
	}
}

func (i *Item) SetContent(content string) {
	i.Content = content
}

func (i Item) FilterValue() string {
	return i.Content
}

func (i Item) State() ItemState {
	if i.HasList && i.state == ItemListClosed {
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

func (i Item) ID() int {
	return i.id
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
	i.Level = l
	return i
}

func (i *Item) SetIsSub() *Item {
	i.IsSub = true
	return i
}

func (i *Item) SetList(l *Model) *Item {
	i.List = l
	return i
}
