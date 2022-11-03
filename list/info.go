package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
)

type infoStr string

func (i infoStr) String() string {
	return string(i)
}

type InfoWidget struct {
	content    []map[fmt.Stringer]fmt.Stringer
	HideKeys   bool
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
}

func NewInfoWidget() *InfoWidget {
	return &InfoWidget{
		KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
	}
}

func (i *InfoWidget) NoKeys() *InfoWidget {
	i.HideKeys = true
	return i
}

func (i *InfoWidget) AddString(key, val string) {
	i.Add(infoStr(key), infoStr(val))
}

func (i *InfoWidget) Add(key, val fmt.Stringer) {
	content := make(map[fmt.Stringer]fmt.Stringer)
	content[key] = val
	i.content = append(i.content, content)
}

func (i *InfoWidget) Set(content ...map[fmt.Stringer]fmt.Stringer) *InfoWidget {
	i.content = content
	return i
}

func (i InfoWidget) String() string {
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
