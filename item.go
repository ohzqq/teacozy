package teacozy

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
	none     string = ``
)

type Item struct {
	Item              list.Item
	idx               int
	IsSelected        bool
	ListOpen          bool
	IsHidden          bool
	MultiSelect       bool
	Level             int
	List              Items
	TotalSubListItems int
	key               string
	value             string
	Form              FormData
	Fields            *Fields
	Data              FieldData
}

func NewItem() *Item {
	return &Item{
		Data: &Field{},
	}
}

// item info
func (i Item) DisplayFields() string {
	return i.Fields.String()
}

func (i *Item) SetFields(f *Fields) {
	i.Fields = f
}

func (i *Item) EditFields() *List {
	return i.Fields.Edit()
}

func (i *Item) SetMultiSelect() *Item {
	i.MultiSelect = true
	return i
}

func (i *Item) SetData(data FieldData) *Item {
	i.Data = data
	i.Item = i
	i.Fields = NewFields().Add(data)
	return i
}

func (i *Item) SetKey(key string) *Item {
	if field, ok := i.Data.(*Field); ok {
		field.key = key
	}
	i.key = key
	return i
}

func (i *Item) SetValue(val string) *Item {
	i.Data.Set(val)
	//i.value = val
	return i
}

func (i Item) ListDepth() int {
	depth := 0
	if i.HasList() {
		depth++
		for _, item := range i.List.all {
			if item.HasList() {
				depth++
			}
		}
	}
	return depth
}

func (i Item) ListLength() int {
	return len(i.Flatten())
}

func (i Item) Flatten() []*Item {
	var items []*Item
	if i.HasList() {
		for _, item := range i.List.all {
			if i.MultiSelect {
				item.SetMultiSelect()
			}
			item.IsHidden = true
			items = append(items, item)
			if item.HasList() {
				items = append(items, item.Flatten()...)
			}
		}
	}
	return items
}

func (i Item) Index() int {
	return i.idx
}

func (i Item) Get(key string) FieldData {
	return i.Fields.Get(key)
}

func (i Item) Keys() []string {
	return i.Fields.Keys()
}

func (i Item) Value() string {
	return i.Data.Value()
}

func (i *Item) Set(content string) {
	i.value = content
	i.Data.Set(content)
}

func (i Item) FilterValue() string {
	return i.Data.Value()
}

func (i Item) Key() string {
	return i.Data.Key()
}

func (i Item) HasList() bool {
	has := len(i.List.all) > 0
	return has
}

func (i *Item) Edit() textarea.Model {
	input := textarea.New()
	input.SetValue(i.value)
	input.ShowLineNumbers = false
	return input
}

func (i Item) Prefix() string {
	if i.HasList() {
		if i.ListOpen {
			return closeSub
		}
		return openSub
	}

	if i.MultiSelect {
		if i.IsSelected {
			return check
		}
		return uncheck
	}

	return dash
}

func (i *Item) ToggleSelected() *Item {
	i.IsSelected = !i.IsSelected
	return i
}

func (i *Item) ToggleList() *Item {
	i.ListOpen = !i.ListOpen
	return i
}

func (i *Item) Show() *Item {
	i.IsHidden = false
	return i
}

func (i *Item) Hide() *Item {
	i.IsHidden = true
	return i
}

func (i *Item) Open() *Item {
	i.ListOpen = true
	return i
}

func (i *Item) Close() *Item {
	i.ListOpen = false
	return i
}

func (i *Item) SetLevel(l int) *Item {
	i.Level = l
	return i
}

func (i *Item) IsSub() bool {
	return i.Level != 0
}
