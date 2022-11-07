package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/slices"
)

var fieldStyle Style

type Style struct {
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
}

type Model struct {
	view viewport.Model
	*Info
}

func New(data FormData) *Model {
	fieldStyle = Style{
		KeyStyle:   lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
		ValueStyle: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
	}
	fields := NewFields().SetData(data)
	m := Model{
		Info: &Info{
			Data: fields,
		},
	}
	height := lipgloss.Height(m.String())
	m.view = viewport.New(util.TermWidth(), height)
	m.view.SetContent(m.String())
	return &m
}

type FormData interface {
	Get(string) string
	Set(string, string)
	Keys() []string
}

type Info struct {
	Data     *Fields
	HideKeys bool
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

type Field struct {
	Key   string
	Value string
}

type UpdateContentMsg struct {
	Field
}

func UpdateContentCmd(key, val string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Field: NewField(key, val)}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case key.Matches(msg, urkey.EditField):
			cmds = append(cmds, UpdateContentCmd("one", "edit"))
		}
	case UpdateContentMsg:
		m.Data.Set(msg.Key, msg.Value)
	case tea.WindowSizeMsg:
		m.view = viewport.New(msg.Width-2, msg.Height-2)
	}
	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (i *Info) NoKeys() *Info {
	i.HideKeys = true
	return i
}

func NewField(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func (f *Fields) Add(key, val string) {
	field := NewField(key, val)
	f.data = append(f.data, field)
}

func (i *Info) Set(f ...map[string]string) *Info {
	var fields Fields
	for _, field := range f {
		for k, v := range field {
			fields.data = append(fields.data, NewField(k, v))
		}
	}
	i.Data = &fields
	return i
}

func (i Info) String() string {
	var info []string
	for _, key := range i.Data.Keys() {
		var line []string
		if !i.HideKeys {
			k := fieldStyle.KeyStyle.Render(key)
			line = append(line, k, ": ")
		}

		val := i.Data.Get(key)
		v := fieldStyle.ValueStyle.Render(val)
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	m.view.SetContent(m.String())
	return m.view.View()
}
