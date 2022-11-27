package info

import (
	"strings"

	"github.com/ohzqq/teacozy"
)

type Section struct {
	title  string
	keys   []string
	values []string
}

func (s *Section) SetTitle(title string) *Section {
	s.title = title
	return s
}

func (s *Section) SetFields(fields teacozy.Fields) *Section {
	for _, key := range fields.Keys() {
		fd := fields.Get(key)

		s.keys = append(s.keys, fd.Name())
		s.values = append(s.values, fd.Content())
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
