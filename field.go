package teacozy

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

var fieldStyle = FieldStyle{
	Key:   lipgloss.NewStyle().Foreground(DefaultColors().Blue),
	Value: lipgloss.NewStyle().Foreground(DefaultColors().DefaultFg),
}

type FieldStyle struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

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
	HideKeys  bool
	IsVisible bool
	Style     FieldStyle
	Data      FormData
	data      []FieldData
}

type Field struct {
	*Item
	key   string
	value string
}

//func NewForm(data FormData) *Fields {
//  f := NewFields().SetData(data)
//  return f
//}

func NewFields() *Fields {
	return &Fields{
		Style: fieldStyle,
	}
}

func NewField(key, val string) *Field {
	return &Field{key: key, value: val}
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
func (f *Fields) NewField(key, val string) *Fields {
	item := NewField(key, val)
	f.data = append(f.data, item)
	return f
}

func (f *Fields) Add(field FieldData) *Fields {
	f.data = append(f.data, field)
	return f
}

func (f *Fields) SetData(data FormData) *Fields {
	f.Data = data
	for _, key := range data.Keys() {
		f.data = append(f.data, data.Get(key))
	}
	return f
}

func (f Fields) Get(key string) FieldData {
	for _, field := range f.data {
		if field.Key() == key {
			return field
		}
	}
	return nil
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.data {
		keys = append(keys, field.Key())
	}
	return keys
}

func (f *Fields) Render() *Fields {
	content := f.String()
	height := lipgloss.Height(content)
	f.Model = viewport.New(TermWidth(), height)
	f.Model.SetContent(content)
	return f
}

func (f *Fields) Info() *Info {
	return NewInfo(f)
}

func (f *Fields) Display() viewport.Model {
	content := f.String()
	height := lipgloss.Height(content)
	return viewport.New(TermWidth(), height)
}

func (f *Fields) Edit() *List {
	items := NewItems()
	if len(f.data) > 0 {
		for _, field := range f.All() {
			i := NewItem().SetData(field)
			items.Add(i)
		}
	}
	form := NewList("Edit...", items)
	form.SetShowKeys()
	form.isForm = true
	return form
}

func (f Fields) All() []FieldData {
	return f.data
}

func (f *Fields) NoKeys() *Fields {
	f.HideKeys = true
	return f
}

func (i Fields) String() string {
	var info []string
	for _, field := range i.All() {
		var line []string
		if !i.HideKeys {
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

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}
