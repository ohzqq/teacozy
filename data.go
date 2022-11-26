package teacozy

type Fields interface {
	Get(string) Field
	Keys() []string
}

type Field interface {
	Value() string
	Key() string
	Set(string)
}
