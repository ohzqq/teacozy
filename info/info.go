package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list/item"
	"github.com/ohzqq/teacozy/prompt"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

var fieldStyle Style

type state int

const (
	view state = iota
	form
	edit
)

type Style struct {
	KeyStyle   lipgloss.Style
	ValueStyle lipgloss.Style
}

type Model struct {
	form  prompt.Model
	view  viewport.Model
	edit  textarea.Model
	state state
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
	m.form = m.Edit()
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

func (i *Info) Edit() prompt.Model {
	items := item.NewItems()
	for _, key := range i.Data.Keys() {
		_, f := i.Data.GetField(key)
		item := item.NewItem(f)
		items.Add(item)
	}
	m := prompt.New()
	m.Title = "Edit..."
	m.SetItems(items).MakeList()
	return m
}

func (i *Info) NoKeys() *Info {
	i.HideKeys = true
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
		if m.edit.Focused() {
			if key.Matches(msg, urkey.SaveAndExit) {
				cur := m.form.List.SelectedItem()
				field := m.form.Items.Get(cur).Data.(Field)
				val := m.edit.Value()
				m.Data.Set(field.Key, val)
				//cur.SetContent(val)
				//m.SetItem(m.List.Model.Index(), cur)
				m.edit.Blur()
				//cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			}
			m.edit, cmd = m.edit.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch m.state {
			case view:
				switch {
				case key.Matches(msg, urkey.EditField):
					cmds = append(cmds, UpdateInfoCmd())
				}
				m.view, cmd = m.view.Update(msg)
				cmds = append(cmds, cmd)
			case form:
				switch {
				case key.Matches(msg, urkey.EditField):
					cmds = append(cmds, prompt.UpdateStatusCmd("edit field"))
				}
				m.form, cmd = m.form.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case EditInfoMsg:
		//m.form = m.Edit()
		m.state = form
	case EditItemMsg:
		m.edit = textarea.New()
		m.edit.SetValue(msg.Value)
		m.edit.ShowLineNumbers = false
		m.edit.Focus()
		//cmds = append(cmds, prompt.UpdateStatusCmd("edit"))
	case UpdateContentMsg:
		m.Data.Set(msg.Key, msg.Value)
	case tea.WindowSizeMsg:
		m.view = viewport.New(msg.Width-2, msg.Height-2)
	case prompt.ReturnSelectionsMsg:
		var field Field
		if items := m.form.Items.Selections(); len(items) > 0 {
			field = items[0].Data.(Field)
		}
		cmds = append(cmds, EditItemCmd(&field))
	case prompt.UpdateStatusMsg:
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var v string
	switch m.state {
	case view:
		m.view.SetContent(m.String())
		v = m.view.View()
	case form:
		v = m.form.View()
	}
	return v
}
