package form

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/ohzqq/teacozy/data"
	"github.com/ohzqq/teacozy/info"
)

type Fields struct {
	fields []*Field
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) Add(fd data.Field) {
	field := NewField(fd)
	f.fields = append(f.fields, field)
}

func (f *Fields) SetData(fields data.Fields) {
	for i, key := range fields.Keys() {
		fd := fields.Get(key)
		field := data.NewField(fd.Key(), fd.Value())
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

func (f *Fields) UndoChanges() *Fields {
	for _, item := range f.fields {
		item.Undo()
	}
	return f
}

func (f *Fields) Items() []list.Item {
	var li []list.Item
	for _, item := range f.fields {
		li = append(li, item)
	}
	return li
}

func (f *Fields) StringMap() map[string]string {
	hash := make(map[string]string)
	for _, item := range f.fields {
		hash[item.Key()] = item.Value()
	}
	return hash
}

func (f *Fields) StringMapChanges() map[string]string {
	hash := make(map[string]string)
	for _, item := range f.fields {
		if item.Changed {
			hash[item.Key()] = item.Value()
		}
	}
	return hash
}

func (f Fields) Get(key string) data.Field {
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
	data    *data.FieldData
}

func NewField(field data.Field) *Field {
	return &Field{
		key:   field.Key(),
		value: field.Value(),
		data:  data.NewField(field.Key(), field.Value()),
	}
}

func (f *Field) Update() {
	f.Changed = true
}

func (f *Field) Save() {
	if f.Value() != f.data.Val {
		f.data.Val = f.value
	}
}

func (f *Field) Undo() {
	f.Changed = false
	f.Set(f.data.Val)
}

// To satisfy field interface
func (f Field) Key() string {
	return f.key
}

func (f Field) Value() string {
	return f.value
}

func (f *Field) Set(val string) {
	f.value = val
}

// To satisfy list item interface
func (f Field) Title() string {
	return f.key
}

func (f Field) Description() string {
	return f.value
}

func (f Field) FilterValue() string {
	return f.value
}
