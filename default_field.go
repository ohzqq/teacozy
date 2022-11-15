package teacozy

type Field struct {
	*Item
	key   string
	value string
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
