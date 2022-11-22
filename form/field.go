package form

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/ohzqq/teacozy"
)

type Fields struct {
	fields []*Field
	Data   []teacozy.FieldData
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) Add(field *Field) {
	f.fields = append(f.fields, field)
	f.Data = append(f.Data, field)
}

func (f *Fields) SetData(data teacozy.FormData) {
	for i, key := range data.Keys() {
		fd := data.Get(key)
		f.Data = append(f.Data, fd)
		field := NewField()
		field.SetKey(fd.Key())
		field.SetValue(fd.Value())
		field.idx = i
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

type Field struct {
	key     string
	value   string
	changed bool
	idx     int
	Data    teacozy.FieldData
}

func NewField() *Field {
	return &Field{
		Data: &Field{},
	}
}

func (i *Field) SetData(data teacozy.FieldData) {
	i.key = data.Key()
	i.value = data.Value()
	i.Data = data
}

func (i Field) Key() string {
	return i.key
}

func (i *Field) SetKey(key string) {
	if field, ok := i.Data.(*Field); ok {
		//field.SetKey(key)
		field.key = key
	}
	i.key = key
}

func (i *Field) SetValue(val string) {
	if field, ok := i.Data.(*Field); ok {
		//field.Set(val)
		field.value = val
	}
	i.value = val
	//i.Set(val)
}

func (i *Field) Changed() {
	i.changed = true
}

func (i Field) Value() string {
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
	if i.Value() != i.Data.Value() {
		i.Data.Set(i.value)
	}
}

func (i *Field) Undo() {
	i.changed = false
	i.Set(i.Data.Value())
}
