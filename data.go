package teacozy

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

func (i FieldData) Name() string {
	return i.label
}

func (i *FieldData) Content() string {
	return i.Val
}

func (i *FieldData) FilterValue() string {
	return i.Val
}

func (i *FieldData) Set(val string) {
	i.Val = val
}
