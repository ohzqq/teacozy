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
	FilterValue() string
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

//func NewForm(data FormData) *Fields {
//  f := NewFields().SetData(data)
//  return f
//}

func NewFields() *Fields {
	return &Fields{
		Style: fieldStyle,
	}
}

func (f *Fields) Render() *Fields {
	content := f.String()
	height := lipgloss.Height(content)
	f.Model = viewport.New(TermWidth(), height)
	f.Model.SetContent(content)
	return f
}

func (f *Fields) Info() *Info {
	return NewInfo(f.Data)
}

func (f *Fields) Display() viewport.Model {
	content := f.String()
	height := lipgloss.Height(content)
	return viewport.New(TermWidth(), height)
}

func (f *Fields) Edit() *List {
	items := NewItems()
	if f.Data != nil {
		for _, field := range f.All() {
			items.Add(NewItem(field))
		}
	}
	form := NewList("Edit...", items)
	return form
}

func (f *Fields) SetData(data FormData) *Fields {
	f.Data = data
	return f
}

func (f Fields) All() []FieldData {
	var fields []FieldData
	if f.Data == nil {
		return fields
	}
	for _, key := range f.Data.Keys() {
		field := f.Data.Get(key)
		fields = append(fields, field)
	}
	return fields
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
