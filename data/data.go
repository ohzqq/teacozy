package data

type Fields interface {
	Get(string) Field
	Keys() []string
}

type Field interface {
	Value() string
	Key() string
	Set(string)
}

type FieldData struct {
	hideKeys bool
	label    string
	Val      string
	changed  bool
}

func NewField(key, val string) *FieldData {
	return &FieldData{
		label: key,
		Val:   val,
	}
}

func (i FieldData) Key() string {
	return i.label
}

func (i *FieldData) Value() string {
	return i.Val
}

func (i *FieldData) FilterValue() string {
	return i.Val
}

func (i *FieldData) Set(val string) {
	i.Val = val
}
