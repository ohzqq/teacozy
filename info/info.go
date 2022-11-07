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
			Fields: fields,
		},
	}
	m.view = m.Display()
	m.form = m.Edit()
	return &m
}

type Info struct {
	Fields   *Fields
	HideKeys bool
}

func (i *Info) Display() viewport.Model {
	content := i.String()
	height := lipgloss.Height(content)
	vp := viewport.New(util.TermWidth(), height)
	vp.SetContent(content)
	return vp
}

func (i *Info) Edit() prompt.Model {
	items := item.NewItems()
	for _, key := range i.Fields.Keys() {
		f := i.Fields.Get(key)
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
	for _, key := range i.Fields.Keys() {
		var line []string
		field := i.Fields.Get(key)
		if !i.HideKeys {
			k := fieldStyle.KeyStyle.Render(field.Key())
			line = append(line, k, ": ")
		}

		v := fieldStyle.ValueStyle.Render(field.Value())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
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
				i := m.form.Items.Get(cur)
				field := i.Data.(FormField)
				val := m.edit.Value()
				field.Set(val)
				m.form.Items.Set(i.Index(), item.NewItem(field))
				m.edit.Blur()
				m.form = m.Edit()
				cmds = append(cmds, prompt.UpdateVisibleItemsCmd("visible"))
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
					cur := m.form.List.SelectedItem()
					field := m.form.Items.Get(cur).Data.(FormField)
					cmds = append(cmds, EditItemCmd(field))
				case key.Matches(msg, urkey.ExitScreen):
					m.state = view
				}
				m.form, cmd = m.form.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case EditInfoMsg:
		m.form = m.Edit()
		m.state = form
	case EditItemMsg:
		m.edit = textarea.New()
		m.edit.SetValue(msg.Value())
		m.edit.ShowLineNumbers = false
		m.edit.Focus()
	case UpdateContentMsg:
		m.Fields.Set(msg.Key(), msg.Value())
	case tea.WindowSizeMsg:
		m.view = viewport.New(msg.Width-2, msg.Height-2)
		m.form.List.SetSize(msg.Width-2, msg.Height-2)
	case prompt.ReturnSelectionsMsg:
		var field FormField
		if items := m.form.Items.Selections(); len(items) > 0 {
			field = items[0].Data.(FormField)
		}
		cmds = append(cmds, EditItemCmd(field))
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
	var (
		sections []string
		//availHeight = m.form.List.Height()
		availHeight = util.TermHeight()
	)

	var field string
	if m.edit.Focused() {
		field = m.edit.View()
		availHeight -= lipgloss.Height(field)
	}

	switch m.state {
	case view:
		m.view.SetContent(m.String())
		v := m.view.View()
		sections = append(sections, v)
	case form:
		m.form.List.SetSize(m.form.Width, availHeight)
		v := m.form.View()
		sections = append(sections, v)
	}

	if m.edit.Focused() {
		sections = append(sections, field)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
