package teacozy

import "github.com/charmbracelet/lipgloss"

type Fields interface {
	Get(string) Field
	Keys() []string
}

type Field interface {
	Content() string
	Name() string
	Set(string)
}

type FieldData struct {
	hideKeys bool
	Style    FieldStyle
	idx      int
	Key      string
	Value    string
	changed  bool
}

func NewField(key, val string) *FieldData {
	return &FieldData{
		Key:   key,
		Value: val,
	}
}

func (i FieldData) Name() string {
	return i.Key
}

func (i *FieldData) Content() string {
	return i.Value
}

func (i *FieldData) FilterValue() string {
	return i.Value
}

func (i *FieldData) Set(val string) {
	i.Value = val
}

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}
