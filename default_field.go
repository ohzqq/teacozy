package teacozy

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type DefaultFields struct {
	data []FieldData
}

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

func NewDefaultField(key, val string) *Field {
	return &Field{key: key, value: val}
}

func (f DefaultFields) Get(key string) FieldData {
	for _, field := range f.data {
		if field.Key() == key {
			return field
		}
	}
	return &Item{}
}

func (f DefaultFields) Keys() []string {
	var keys []string
	for _, field := range f.data {
		keys = append(keys, field.Key())
	}
	return keys
}

func (f DefaultFields) Has(key string) bool {
	return slices.Contains(f.Keys(), key)
}

func (f *DefaultFields) Add(key, val string) error {
	if f.Has(key) {
		return fmt.Errorf("keys must be unique")
	}
	field := NewDefaultField(key, val)
	f.data = append(f.data, field)
	return nil
}
