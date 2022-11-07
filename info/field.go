package info

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type FormData interface {
	Get(string) FormField
	Set(string, string)
	Keys() []string
}

type FormField interface {
	FilterValue() string
	Value() string
	Key() string
	Set(string)
}

type Field struct {
	idx   int
	key   string
	value string
}

func NewField(key, val string) Field {
	return Field{
		key:   key,
		value: val,
	}
}

func (f Field) FilterValue() string {
	return f.value
}

func (f Field) Value() string {
	return f.value
}

func (f *Field) Set(val string) {
	f.value = val
}

func (f Field) Key() string {
	return f.key
}

type Fields struct {
	data   []Field
	fields []FormField
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) SetData(data FormData) *Fields {
	for _, key := range data.Keys() {
		val := data.Get(key)
		f.Add(key, val.Value())
	}
	return f
}

func (f Fields) Get(key string) FormField {
	for _, field := range f.fields {
		if field.Key() == key {
			return field
		}
	}
	return &Field{}
}

func (f Fields) GetField(key string) (int, Field) {
	for idx, field := range f.data {
		if field.key == key {
			return idx, field
		}
	}
	return -1, Field{}
}

func (f *Fields) Set(key, val string) {
	if f.Has(key) {
		ff := f.Get(key)
		ff.Set(val)
	} else {
		field := NewField(key, val)
		f.data = append(f.data, field)
	}
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.data {
		keys = append(keys, field.key)
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
	field := NewField(key, val)
	f.data = append(f.data, field)
	f.fields = append(f.fields, &field)
	return nil
}
