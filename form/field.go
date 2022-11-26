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

func (f *Fields) Add(field *Field) {
	f.fields = append(f.fields, field)
	f.Data = append(f.Data, field)
}

func (f *Fields) SetData(data teacozy.Fields) {
	for i, key := range data.Keys() {
		fd := data.Get(key)
		f.Data = append(f.Data, fd)
		field := NewField()
		field.SetKey(fd.Name())
		field.SetValue(fd.Content())
		field.idx = i
		field.Data = field
		f.fields = append(f.fields, field)
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
	for _, field := range f.Data {
		if field.Name() == key {
			return field
		}
	}
	return nil
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.Data {
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
	Data    teacozy.Field
}

func NewField() *Field {
	return &Field{
		Data: &Field{},
		data: &teacozy.FieldData{},
	}
}

func (i *Field) SetData(data teacozy.Field) {
	i.key = data.Name()
	i.value = data.Content()
	i.data = teacozy.NewField(data.Name(), data.Content())
	i.Data = data
}

func (i *Field) SetKey(key string) {
	i.key = key
}

func (i *Field) SetValue(val string) {
	i.data.Value = val
	i.value = val
}

func (i *Field) Update() {
	i.Changed = true
}

func (i Field) Name() string {
	return i.key
}

func (i Field) Content() string {
	return i.value
}

func (i Field) Title() string {
	return i.key
}

func (i Field) Description() string {
	return i.value
}

func (i Field) FilterValue() string {
	return i.value
}

func (i *Field) Set(val string) {
	i.value = val
}

func (i *Field) Save() {
	if i.Content() != i.data.Value {
		i.data.Value = i.value
	}
}

func (i *Field) Undo() {
	i.Changed = false
	i.Set(i.data.Value)
}
