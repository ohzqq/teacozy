package info

import (
	"strings"

	"github.com/ohzqq/teacozy/data"
)

type Section struct {
	title  string
	keys   []string
	values []string
}

func NewSection() *Section {
	return &Section{}
}

func (s *Section) SetTitle(title string) *Section {
	s.title = title
	return s
}

func (s *Section) SetFields(fields data.Fields) *Section {
	for _, key := range fields.Keys() {
		fd := fields.Get(key)

		s.keys = append(s.keys, fd.Key())
		s.values = append(s.values, fd.Value())
	}
	return s
}

func (s Section) Render(style Style, hideKeys bool) string {
	var content []string

	if title := s.title; title != "" {
		t := style.Title.Render(title)
		content = append(content, t)
	}

	for idx, val := range s.values {
		var line []string
		if !hideKeys {
			k := style.Key.Render(s.keys[idx])
			line = append(line, k+": ")
		}

		v := style.Value.Render(val)
		line = append(line, v)

		l := strings.Join(line, "")
		content = append(content, l)
	}

	return strings.Join(content, "\n")
}
