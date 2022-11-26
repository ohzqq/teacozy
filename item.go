package teacozy

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	idx          int
	IsHidden     bool
	IsSelected   bool
	MultiSelect  bool
	ShowChildren bool
	showKey      bool
	Level        int
	Parent       *Item
	Children     Items
	changed      bool
	key          string
	value        string
	hasFields    bool
	Fields       *FormData
	Data         Field
}

func NewItem() *Item {
	return &Item{
		Data:   &FieldData{},
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

func (i *Item) SetFields(f *FormData) {
	i.Fields = f
}

func (i *Item) EditFields() *List {
	return i.Fields.Edit()
}

func (i *Item) SetMultiSelect() *Item {
	i.MultiSelect = true
	return i
}

func (i *Item) SetData(data Field) *Item {
	i.key = data.Name()
	i.value = data.Content()
	i.Data = data
	i.Fields = NewFields().Add(data)
	return i
}

func (i *Item) SetKey(key string) *Item {
	if field, ok := i.Data.(*FieldData); ok {
		field.Key = key
	}
	i.key = key
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

func (i Item) HasChildren() bool {
	has := len(i.Children.flat) > 0
	return has
}

func (i Item) TotalChildren() int {
	if i.HasChildren() {
		return len(i.Children.flat)
	}
	return 0
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

func (i Item) Get(key string) Field {
	return i.Fields.Get(key)
}

func (i Item) Keys() []string {
	return i.Fields.Keys()
}

func (i *Item) Set(content string) {
	i.value = content
}

func (i *Item) Changed() *Item {
	i.changed = true
	return i
}

func (i *Item) ChangedCmd() tea.Cmd {
	return func() tea.Msg {
		i.Changed()
		return ItemChangedMsg{Item: i}
	}
}

func (i *Item) Save() {
	if i.value != i.Data.Content() {
		i.Data.Set(i.value)
	}
}

func (i *Item) Undo() {
	i.changed = false
	i.value = i.Data.Content()
}

func (i Item) Value() string {
	//return i.Data.Value()
	return i.value
}

func (i Item) FilterValue() string {
	return i.value
}

func (i Item) Name() string {
	return i.key
}

func (i Item) String() string {
	var item string
	if i.showKey {
		item = i.Name() + ": "
	}
	return item + i.Value()
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
