package teacozy

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
)

type FormData struct {
	Model     viewport.Model
	IsVisible bool
	hideKeys  bool
	Style     FieldStyle
	Data      []Field
	fields    []*FieldData
	data      *Items
}

func NewFields() *FormData {
	return &FormData{
		Style: fieldStyle,
		data:  NewItems(),
	}
}

func DisplayFields(fields *FormData, w, h int) *Info {
	info := NewInfo().SetData(fields)
	info.SetSize(w, h)
	return info
}

func (f *FormData) NewField(key, val string) *FormData {
	field := NewField(key, val)
	f.Data = append(f.Data, field)
	item := NewItem().SetData(field)
	f.data.Add(item)
	return f
}

func (f *FormData) Add(field Field) *FormData {
	f.Data = append(f.Data, field)
	return f
}

func (f FormData) HasData() bool {
	return len(f.Data) > 0
}

func (f FormData) Items() *Items {
	return f.data
}

func (f *FormData) SetData(data Fields) *FormData {
	for i, key := range data.Keys() {
		fd := data.Get(key)
		f.Data = append(f.Data, fd)
		item := NewItem().SetData(fd)
		f.data.Add(item)
		field := NewField(fd.Name(), fd.Content())
		field.idx = i
		f.fields = append(f.fields, field)
	}
	return f
}

func (f FormData) Get(key string) Field {
	for _, field := range f.Data {
		if field.Name() == key {
			return field
		}
	}
	return nil
}

func (f FormData) Keys() []string {
	var keys []string
	for _, field := range f.Data {
		keys = append(keys, field.Name())
	}
	return keys
}

func (f *FormData) Edit() *List {
	form := NewList().SetTitle("edit")
	if len(f.Data) > 0 {
		for _, field := range f.All() {
			i := NewItem().SetData(field)
			form.Add(i)
		}
	}
	return form.Edit()
}

func (f FormData) All() []Field {
	return f.Data
}

func (f *FormData) HideKeys() *FormData {
	f.hideKeys = true
	return f
}

func (i FormData) String() string {
	var info []string
	for _, field := range i.All() {
		var line []string
		if !i.hideKeys {
			k := i.Style.Key.Render(field.Name())
			line = append(line, k, ": ")
		}

		v := i.Style.Value.Render(field.Content())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}

	return strings.Join(info, "\n")
}
