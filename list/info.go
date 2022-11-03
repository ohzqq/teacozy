package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
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
	w := InfoWidget{
		KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
	}

	w.AddString("", "")

	return &w
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
				if str := key.String(); str != "" {
					k := i.KeyStyle.Render(str)
					line = append(line, k, ": ")
				}
			}
			if str := val.String(); str != "" {
				v := i.KeyStyle.Render(str)
				line = append(line, v)
			}
		}
		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}

func (m *Model) UpdateInfoWidget(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, urkey.Info):
			m.ToggleInfo()
			cmds = append(cmds, SetFocusedViewCmd("list"))
		}
	}

	m.info, cmd = m.info.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
