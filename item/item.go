package item

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
)

const (
	check    string = "[x] "
	uncheck  string = "[ ] "
	dash     string = "- "
	openSub  string = `[+] `
	closeSub string = `[-] `
)

type Item struct {
	idx        int
	data       list.Item
	IsSelected bool
	hasList    bool
	ListOpen   bool
	IsHidden   bool
	IsMulti    bool
	Level      int
	List       Items
	Label      string
	Content    string
	//Info       *list.InfoWidget
}

func NewItem(item list.Item) Item {
	return Item{
		data:    item,
		Content: item.FilterValue(),
		//Info:    list.NewInfoWidget(),
	}
}

func NewDefaultItem(content string) Item {
	item := Item{
		Content: content,
	}
	i := NewItem(item)
	return i
}

func (i *Item) SetLabel(label string) {
	i.Label = label
}

func (i Item) Flatten() []Item {
	var items []Item
	if i.HasList() {
		for _, item := range i.List.all {
			item.IsHidden = true
			items = append(items, item)
			if item.HasList() {
				items = append(items, item.Flatten()...)
			}
		}
	}
	return items
}

//func (i Item) DisplayInfo() string {
//  return i.Info.String()
//}

func (i *Item) SetContent(content string) {
	i.Content = content
}

func (i Item) FilterValue() string {
	return i.Content
}

func (i Item) HasList() bool {
	has := len(i.List.all) > 0
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
		if i.ListOpen {
			return closeSub
		}
		return openSub
	} else {
		if i.IsSelected {
			return check
		}
		return uncheck
	}

	return dash
}

func (i *Item) SetIndex(id int) *Item {
	i.idx = id
	return i
}

func (i Item) Index() int {
	return i.idx
}

func (i *Item) ToggleSelected() Item {
	i.IsSelected = !i.IsSelected
	return *i
}

func (i *Item) ToggleList() Item {
	i.ListOpen = !i.ListOpen
	return *i
}

func (i *Item) SetLevel(l int) *Item {
	i.Level = l
	return i
}

func (i *Item) IsSub() bool {
	return i.Level != 0
}
