package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
)

type Info struct {
	content    []map[fmt.Stringer]fmt.Stringer
	HideKeys   bool
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
}

func NewInfo() *Info {
	return &Info{
		KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
	}
}

func (i *Info) NoKeys() *Info {
	i.HideKeys = true
	return i
}

func (i *Info) Add(key, val fmt.Stringer) {
	content := make(map[fmt.Stringer]fmt.Stringer)
	content[key] = val
	i.content = append(i.content, content)
}

func (i *Info) Set(content ...map[fmt.Stringer]fmt.Stringer) *Info {
	i.content = content
	return i
}

func (i Info) String() string {
	var info []string
	for _, pair := range i.content {
		var line []string
		for key, val := range pair {
			if !i.HideKeys {
				k := i.KeyStyle.Render(key.String())
				line = append(line, k)
			}
			v := i.ValueStyle.Render(val.String())
			line = append(line, v)
		}
		l := strings.Join(line, ": ")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}
