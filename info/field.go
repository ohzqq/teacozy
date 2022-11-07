package info

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type FormData interface {
	Get(string) string
	Set(string, string)
	Keys() []string
}

type Field struct {
	idx   int
	Key   string
	Value string
}

func NewField(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func (f Field) FilterValue() string {
	return f.Value
}

type Fields struct {
	data []Field
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) SetData(data FormData) *Fields {
	for _, key := range data.Keys() {
		val := data.Get(key)
		f.Add(key, val)
	}
	return f
}

func (f Fields) Get(key string) string {
	for _, field := range f.data {
		if field.Key == key {
			return field.Value
		}
	}
	return ""
}

func (f Fields) GetField(key string) (int, Field) {
	for idx, field := range f.data {
		if field.Key == key {
			return idx, field
		}
	}
	return -1, Field{}
}

func (f *Fields) Set(key, val string) {
	if f.Has(key) {
		idx, field := f.GetField(key)
		field.Value = val
		f.data[idx] = field
	} else {
		field := NewField(key, val)
		f.data = append(f.data, field)
	}
}

func (f Fields) Keys() []string {
	var keys []string
	for _, field := range f.data {
		keys = append(keys, field.Key)
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
	field.idx = len(f.data)
	f.data = append(f.data, field)
	return nil
}
