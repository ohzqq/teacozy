package form

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/info"
)

type Fields struct {
	fields []*Field
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) Add(fd teacozy.Field) {
	field := NewField(fd)
	f.fields = append(f.fields, field)
}

func (f *Fields) SetData(data teacozy.Fields) {
	for i, key := range data.Keys() {
		fd := data.Get(key)
		field := teacozy.NewField(fd.Key(), fd.Value())
		ff := NewField(field)
		ff.idx = i
		f.fields = append(f.fields, ff)
	}
}

func (f *Fields) PreviewChanges() *info.Section {
	return info.NewSection().SetFields(f)
}

func (f *Fields) ConfirmChanges() *info.Section {
	return f.PreviewChanges().SetTitle(`save and exit? y\n`)
}

func (f *Fields) SaveChanges() *Fields {
	for _, item := range f.fields {
		item.Save()
	}
	return f
}

func (i *Fields) UndoChanges() *Fields {
	for _, item := range i.fields {
		item.Undo()
	}
	return i
}

func (i *Fields) Items() []list.Item {
	var li []list.Item
	for _, item := range i.fields {
		li = append(li, item)
	}
	return li
}

func (f Fields) Get(key string) teacozy.Field {
	for _, field := range f.fields {
		if field.Key() == key {
			return field
		}
	}
	return nil
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.fields {
		keys = append(keys, field.Key())
	}
	return keys
}

type Field struct {
	key     string
	value   string
	Changed bool
	idx     int
	data    *teacozy.FieldData
}

func NewField(data teacozy.Field) *Field {
	return &Field{
		key:   data.Key(),
		value: data.Value(),
		data:  teacozy.NewField(data.Key(), data.Value()),
	}
}

func (i *Field) Update() {
	i.Changed = true
}

func (i *Field) Save() {
	if i.Value() != i.data.Val {
		i.data.Val = i.value
	}
}

func (i *Field) Undo() {
	i.Changed = false
	i.Set(i.data.Val)
}

// To satisfy field interface
func (i Field) Key() string {
	return i.key
}

func (i Field) Value() string {
	return i.value
}

func (i *Field) Set(val string) {
	i.value = val
}

// To satisfy list item interface
func (i Field) Title() string {
	return i.key
}

func (i Field) Description() string {
	return i.value
}

func (i Field) FilterValue() string {
	return i.value
}
