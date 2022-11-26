package form

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/ohzqq/teacozy"
)

type Fields struct {
	fields []*Field
	Data   []teacozy.Field
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) Add(fd teacozy.Field) {
	field := NewField(fd)
	f.fields = append(f.fields, field)
	f.Data = append(f.Data, fd)
}

func (f *Fields) SetData(data teacozy.Fields) {
	for i, key := range data.Keys() {
		fd := data.Get(key)
		f.Data = append(f.Data, fd)
		field := teacozy.NewField(fd.Name(), fd.Content())
		ff := NewField(field)
		ff.idx = i
		f.fields = append(f.fields, ff)
	}
}

func (i *Fields) SaveChanges() *Fields {
	for _, item := range i.fields {
		item.Save()
	}
	return i
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
		if field.Name() == key {
			return field
		}
	}
	return nil
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.fields {
		keys = append(keys, field.Name())
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
		key:   data.Name(),
		value: data.Content(),
		data:  teacozy.NewField(data.Name(), data.Content()),
	}
}

func (i *Field) Update() {
	i.Changed = true
}

func (i *Field) Save() {
	if i.Content() != i.data.Val {
		i.data.Val = i.value
	}
}

func (i *Field) Undo() {
	i.Changed = false
	i.Set(i.data.Val)
}

// To satisfy field interface
func (i Field) Name() string {
	return i.key
}

func (i Field) Content() string {
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
