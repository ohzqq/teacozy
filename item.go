package teacozy

import (
	"github.com/charmbracelet/bubbles/textarea"
)

type Item struct {
	idx           int
	IsHidden      bool
	IsSelected    bool
	MultiSelect   bool
	ShowChildren  bool
	Level         int
	TotalChildren int
	Parent        *Item
	Children      Items
	Changed       bool
	hasFields     bool
	Fields        *Fields
	Data          FieldData
}

func NewItem() *Item {
	return &Item{
		Data:   &Field{},
		Fields: NewFields(),
	}
}

// item info
func (i Item) DisplayFields() string {
	return i.Fields.String()
}

func (i Item) HasFields() bool {
	return i.Fields.HasData()
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
	i.Fields = NewFields().Add(data)
	return i
}

func (i *Item) SetKey(key string) *Item {
	if field, ok := i.Data.(*Field); ok {
		field.key = key
	}
	return i
}

func (i *Item) SetValue(val string) *Item {
	i.Data.Set(val)
	return i
}

func (i Item) ListDepth() int {
	depth := 0
	if i.HasChildren() {
		depth++
		for _, item := range i.Children.flat {
			if item.HasChildren() {
				depth++
			}
		}
	}
	return depth
}

func (i Item) ListLength() int {
	return len(i.Flatten())
}

func (i Item) HasChildren() bool {
	has := len(i.Children.flat) > 0
	return has
}

func (i *Item) Flatten() []*Item {
	var items []*Item
	if i.HasChildren() {
		for _, item := range i.Children.flat {
			if i.MultiSelect {
				item.SetMultiSelect()
			}
			item.IsHidden = true
			items = append(items, item)
			if item.HasChildren() {
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
	i.Data.Set(content)
}

func (i Item) FilterValue() string {
	return i.Data.Value()
}

func (i Item) Key() string {
	return i.Data.Key()
}

func (i *Item) Edit() textarea.Model {
	input := textarea.New()
	input.SetValue(i.Value())
	input.ShowLineNumbers = false
	return input
}

func (i *Item) ToggleSelected() *Item {
	i.IsSelected = !i.IsSelected
	return i
}

func (i *Item) Select() *Item {
	i.IsSelected = true
	return i
}

func (i *Item) Deselect() *Item {
	i.IsSelected = false
	return i
}

func (i *Item) ToggleList() *Item {
	i.ShowChildren = !i.ShowChildren
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
	i.ShowChildren = true
	return i
}

func (i *Item) Close() *Item {
	i.ShowChildren = false
	return i
}

func (i *Item) SetLevel(l int) *Item {
	i.Level = l
	return i
}

func (i *Item) IsSub() bool {
	return i.Level != 0
}
