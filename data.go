package teacozy

import "github.com/charmbracelet/lipgloss"

type Fields interface {
	Get(string) Field
	Keys() []string
}

type Field interface {
	Content() string
	Key() string
	Set(string)
}

type FieldData struct {
	hideKeys bool
	Style    FieldStyle
	idx      int
	key      string
	value    string
	changed  bool
}

func NewField(key, val string) *FieldData {
	return &FieldData{
		key:   key,
		value: val,
	}
}

func (i FieldData) Key() string {
	return i.key
}

func (i *FieldData) Content() string {
	return i.value
}

func (i *FieldData) FilterValue() string {
	return i.value
}

func (i *FieldData) Set(val string) {
	i.value = val
}

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}
