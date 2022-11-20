package teacozy

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type FormData interface {
	Get(string) FieldData
	Keys() []string
}

type FieldData interface {
	Value() string
	Key() string
	Set(string)
}

type Fields struct {
	Model     viewport.Model
	hideKeys  bool
	IsVisible bool
	Style     FieldStyle
	Data      []FieldData
}

func NewFields() *Fields {
	return &Fields{
		Style: fieldStyle,
	}
}

func DisplayFields(fields *Fields, w, h int) *Info {
	info := NewInfo().SetData(fields)
	info.SetSize(w, h)
	return info
}

func (f *Fields) NewField(key, val string) *Fields {
	item := NewField(key, val)
	f.Data = append(f.Data, item)
	return f
}

func (f *Fields) Add(field FieldData) *Fields {
	f.Data = append(f.Data, field)
	return f
}

func (f Fields) HasData() bool {
	return len(f.Data) > 0
}

func (f *Fields) SetData(data FormData) *Fields {
	for _, key := range data.Keys() {
		f.Data = append(f.Data, data.Get(key))
	}
	return f
}

func (f Fields) Get(key string) FieldData {
	for _, field := range f.Data {
		if field.Key() == key {
			return field
		}
	}
	return nil
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.Data {
		keys = append(keys, field.Key())
	}
	return keys
}

func (f *Fields) Edit() *List {
	form := NewList().SetTitle("edit")
	if len(f.Data) > 0 {
		for _, field := range f.All() {
			i := NewItem().SetData(field)
			form.Add(i)
		}
	}
	return form.Edit()
}

func (f Fields) All() []FieldData {
	return f.Data
}

func (f *Fields) HideKeys() *Fields {
	f.hideKeys = true
	return f
}

func (i Fields) String() string {
	var info []string
	for _, field := range i.All() {
		var line []string
		if !i.hideKeys {
			k := i.Style.Key.Render(field.Key())
			line = append(line, k, ": ")
		}

		v := i.Style.Value.Render(field.Value())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}

	return strings.Join(info, "\n")
}

type Field struct {
	key   string
	value string
}

func NewField(key, val string) *Field {
	return &Field{
		key:   key,
		value: val,
	}
}

func (i Field) Key() string {
	return i.key
}

func (i *Field) Value() string {
	return i.value
}

func (i *Field) FilterValue() string {
	return i.value
}

func (i *Field) Set(val string) {
	i.value = val
}

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}
