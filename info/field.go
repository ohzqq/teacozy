package info

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type FormData interface {
	Get(string) Field
	Set(string, string)
	Keys() []string
}

type Field interface {
	FilterValue() string
	Value() string
	Key() string
	Set(string)
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

type Fields struct {
	data []Field
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) SetData(data FormData) *Fields {
	for _, key := range data.Keys() {
		field := data.Get(key)
		f.Add(field.Key(), field.Value())
	}
	return f
}

func (f Fields) Get(key string) Field {
	for _, field := range f.data {
		if field.Key() == key {
			return field
		}
	}
	return &DefaultField{}
}

func (f *Fields) Set(key, val string) {
	if f.Has(key) {
		ff := f.Get(key)
		ff.Set(val)
	} else {
		f.Add(key, val)
	}
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.data {
		keys = append(keys, field.Key())
	}
	return keys
}

func (f Fields) Has(key string) bool {
	return slices.Contains(f.Keys(), key)
}

func (f *Fields) Add(key, val string) error {
	if f.Has(key) {
		return fmt.Errorf("keys must be unique")
	}
	field := NewDefaultField(key, val)
	f.data = append(f.data, field)
	return nil
}
