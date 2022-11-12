package teacozy

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var fieldStyle = FieldStyle{
	Key:   lipgloss.NewStyle().Foreground(DefaultColors().Blue),
	Value: lipgloss.NewStyle().Foreground(DefaultColors().DefaultFg),
}

type FieldStyle struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

type FormData interface {
	Get(string) FieldData
	Keys() []string
}

type FieldData interface {
	FilterValue() string
	Value() string
	Key() string
	Set(string)
}

type Fields struct {
	Model     viewport.Model
	HideKeys  bool
	IsVisible bool
	Style     FieldStyle
	Data      FormData
	data      []FieldData
}

//func NewForm(data FormData) *Fields {
//  f := NewFields().SetData(data)
//  return f
//}

func NewFields() *Fields {
	return &Fields{
		Style: fieldStyle,
	}
}

func (f *Fields) Render() *Fields {
	content := f.String()
	height := lipgloss.Height(content)
	f.Model = viewport.New(TermWidth(), height)
	f.Model.SetContent(content)
	return f
}

func (f *Fields) SetData(data FormData) *Fields {
	f.Data = data
	return f
}

func (f Fields) AllFields() []FieldData {
	var fields []FieldData
	if f.Data == nil {
		return fields
	}
	for _, key := range f.Data.Keys() {
		field := f.Data.Get(key)
		fields = append(fields, field)
	}
	return fields
}

func (f *Fields) NoKeys() *Fields {
	f.HideKeys = true
	return f
}

func (m *Fields) Update(msg tea.Msg) (*Fields, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case key.Matches(msg, Keys.ExitScreen):
			cmds = append(cmds, HideCmd())
		case key.Matches(msg, Keys.EditField):
			cmds = append(cmds, EditCmd(m))
		}
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (i Fields) String() string {
	var info []string
	for _, field := range i.AllFields() {
		var line []string
		if !i.HideKeys {
			k := i.Style.Key.Render(field.Key())
			line = append(line, k, ": ")
		}

		v := i.Style.Value.Render(field.Value())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}

func (m *Fields) View() string {
	m.Model.SetContent(m.String())
	return m.Model.View()
}

func (i *Fields) Init() tea.Cmd {
	return nil
}

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}
