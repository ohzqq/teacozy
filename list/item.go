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

func (s ItemState) Prefix() string {
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

func (li Items) Flatten() Items {
	var items Items
	for _, item := range li {
		i := item.(Item)
		items = append(items, i)
		if i.HasList() {
			items = append(items, i.Items.Flatten()...)
		}
	}
	return items
}

func FlattenItems(li []Item) []list.Item {
	var items []list.Item
	for _, item := range li {
		items = append(items, item)
		if item.HasList() {
			items = append(items, FlattenItems(item.items)...)
		}
	}
	return items
}

func (li Items) Display(opt string) Items {
	switch opt {
	case "selected":
		return li.Selected()
	default:
		return li.Visible()
	}
}

func (li Items) Visible() Items {
	var items Items
	level := 0
	for _, item := range li {
		i := item.(Item)
		if !i.IsHidden {
			items = append(items, i)
		}
		if i.HasList() && i.listOpen {
			level++
			for _, sub := range li.GetSubList(i) {
				s := sub.(Item)
				s.IsHidden = false
				s.SetLevel(level)
				items = append(items, s)
			}
		}
	}
	return items
}

func (li Items) GetSubList(i list.Item) Items {
	item := i.(Item)
	if item.HasList() {
		t := len(item.Items)
		return li[item.id+1 : item.id+t+1]
	}
	return Items{}
}

func (i Items) Add(item Item) Items {
	i = i.appendItem(item)
	if item.HasList() {
		for _, l := range item.Items {
			li := l.(Item)
			li.IsHidden = true
			i = i.appendItem(li)
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

func (li Items) Selected() Items {
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
	return item
}

func (li Items) CloseList(idx int) Item {
	item := li.Get(idx)
	return item
}

func (li Items) ToggleList(idx int) Item {
	item := li.Get(idx)

	if item.HasList() {
		item.listOpen = !item.listOpen
	}

	li[idx] = item

	return li.Get(idx)
}

func (li Items) Select(idx int) Item {
	item := li.Get(idx)
	return item
}

func (li Items) Set(i list.Item) {
	li[i.(Item).Index()] = i
}

func (li Items) Deselect(idx int) Item {
	item := li.Get(idx)
	return item
}

func (li Items) NewList(title string, state listType) *Model {
	l := New(title)
	l.state = state
	return l
}

type Item struct {
	id         int
	isSelected bool
	hasList    bool
	listOpen   bool
	IsVisible  bool
	IsHidden   bool
	mark       string
	IsMulti    bool
	Level      int
	items      []Item
	Items      Items
	Content    string
	Info       *InfoWidget
}

func NewItem(item list.Item) Item {
	return Item{
		Content: item.FilterValue(),
		Info:    NewInfoWidget(),
	}
}

func (i Item) DisplayInfo() string {
	return i.Info.String()
}

func (i *Item) SetContent(content string) {
	i.Content = content
}

func (i Item) FilterValue() string {
	return i.Content
}

func (i Item) HasList() bool {
	has := len(i.Items) > 0
	return has
}

func (i *Item) Edit() textarea.Model {
	input := textarea.New()
	input.SetValue(i.Content)
	input.ShowLineNumbers = false
	return input
}

func (i Item) Prefix() string {
	if i.HasList() {
		if i.listOpen {
			return closeSub
		}
		return openSub
	} else {
		if i.IsSelected() {
			return check
		}
		return uncheck
	}

	return dash
}

func (i *Item) SetId(id int) *Item {
	i.id = id
	return i
}

func (i Item) Index() int {
	return i.id
}

func (i Item) IsSelected() bool {
	if i.isSelected {
		return true
	}
	return false
}

func (i *Item) ToggleSelected() Item {
	i.isSelected = !i.isSelected
	return *i
}

func (i *Item) ToggleList() Item {
	i.listOpen = !i.listOpen
	return *i
}

func (i *Item) SetLevel(l int) *Item {
	i.Level = l
	return i
}

func (i *Item) IsSub() bool {
	return i.Level != 0
}
