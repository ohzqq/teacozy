package teacozy

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type DefaultFields struct {
	data []Field
}

type DefaultField struct {
	key   string
	value string
}

func NewDefaultField(key, val string) *DefaultField {
	return &DefaultField{
		key:   key,
		value: val,
	}
}

func (f DefaultFields) Get(key string) Field {
	for _, field := range f.data {
		if field.Key() == key {
			return field
		}
	}
	return &DefaultField{}
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

func (f DefaultField) FilterValue() string {
	return f.value
}

func (f DefaultField) Value() string {
	return f.value
}

func (f *DefaultField) Set(val string) {
	f.value = val
}

func (f DefaultField) Key() string {
	return f.key
}
