package form

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

type Form struct {
	*Fields
	style.Frame
	Model        list.Model
	Input        textarea.Model
	Hash         map[string]string
	Changed      bool
	view         bool
	Info         *info.Info
	section      *info.Section
	SaveFormFunc SaveFormFunc
}

func New(fields *Fields) Form {
	m := Form{
		Fields:       fields,
		Frame:        style.DefaultFrameStyle(),
		SaveFormFunc: SaveFormAsHash,
	}

	m.Model = list.New(
		m.Fields.Items(),
		itemDelegate(),
		m.Width(),
		m.Height(),
	)
	m.Model.SetShowStatusBar(false)
	m.Model.SetShowHelp(false)
	m.Model.Styles = style.ListStyles()
	m.section = info.NewSection().SetFields(fields)
	m.Info = info.New().SetSections(m.section)
	m.Info.Editable = true

	return m
}

func (m *Form) SetFields(fields *Fields) {
	m.Fields = fields
}

func (m *Form) Start() {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.Hash)
}

func (m *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		if m.Input.Focused() {
			if key.Matches(msg, key.SaveAndExit) {
				cur := m.Model.SelectedItem().(*Field)
				val := m.Input.Value()
				if original := cur.Content(); original != val {
					cur.Set(val)
					m.Changed = true
					m.Info.SetSections(m.Fields.PreviewChanges())
					cmds = append(cmds, m.FieldChanged(cur))
				}
				m.Input.Blur()
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case m.view:
				switch {
				case key.Matches(msg, key.ToggleAllItems):
					m.view = !m.view
					m.Info.ToggleVisible()
				case key.Matches(msg, key.SaveAndExit):
					m.Info.SetSections(m.Fields.ConfirmChanges())
					cmds = append(cmds, ViewFormCmd())
				case key.Matches(msg, Yes):
					s := m.Info.Sections[0]
					s.SetTitle("")
					m.section = s
					m.SaveChanges()
					cmds = append(cmds, HideFormCmd())
					cmds = append(cmds, SaveAndExitFormCmd())
				case key.Matches(msg, No):
					m.Info.SetSections(m.section)
					m.UndoChanges()
					cmds = append(cmds, HideFormCmd())
				case key.Matches(msg, key.EditField):
					cmds = append(cmds, HideFormCmd())
				case key.Matches(msg, key.PrevScreen):
					cmds = append(cmds, HideFormCmd())
				}
				var mod tea.Model
				mod, cmd = m.Info.Update(msg)
				m.Info = mod.(*info.Info)
				//m.Info, cmd = m.Info.Update(msg)
				cmds = append(cmds, cmd)
			default:
				switch {
				case key.Matches(msg, key.ToggleAllItems):
					m.view = !m.view
					m.Info.ToggleVisible()
				case key.Matches(msg, key.SaveAndExit):
					switch {
					case m.Changed:
						cmds = append(cmds, ViewFormCmd())
					default:
						//cmds = append(cmds, ExitFormCmd())
					}
				}
				m.Model, cmd = m.Model.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case teacozy.UpdateStatusMsg:
		cmds = append(cmds, m.Model.NewStatusMessage(msg.Msg))
	case tea.WindowSizeMsg:
		m.Frame.SetSize(msg.Width-1, msg.Height-2)
		m.Model.SetSize(msg.Width-1, msg.Height-2)
		m.Info.SetSize(msg.Width-1, msg.Height-2)
	case EditFormItemMsg:
		cur := m.Model.SelectedItem().(*Field)
		m.Input = textarea.New()
		m.Input.SetValue(cur.Content())
		m.Input.ShowLineNumbers = false
		m.Input.Focus()
	case teacozy.SetListItemMsg:
		idx := m.Model.Index()
		m.Model.SetItem(idx, msg.Item)
	case ViewFormMsg:
		m.view = true
		m.Info.Show()
	case HideFormMsg:
		m.Info.Hide()
		m.view = false
	case SaveAndExitFormMsg:
		cmds = append(cmds, m.SaveFormFunc(m))
		cmds = append(cmds, FormChangedCmd())
	}

	return m, tea.Batch(cmds...)
}

func (m *Form) FieldChanged(item *Field) tea.Cmd {
	return func() tea.Msg {
		item.Update()
		return teacozy.SetListItemMsg{Item: item}
	}
}

func (m *Form) Init() tea.Cmd {
	return nil
}

func (m Form) View() string {
	var (
		sections    []string
		availHeight = m.Height()
		field       string
	)

	if m.view {
		info := m.Info.View()
		//info := m.Confirm.View()
		return info
		//availHeight -= m.Info.Frame.Height()
		//sections = append(sections, info)
	} else {
		if m.Input.Focused() {
			iHeight := availHeight / 3
			m.Input.SetHeight(iHeight)
			field = m.Input.View()
			availHeight -= iHeight
		}

		//m.Frame.SetSize(m.Frame.Width(), availHeight)
		m.Model.SetSize(m.Width(), availHeight)
		content := m.Model.View()
		sections = append(sections, content)

		if m.Input.Focused() {
			sections = append(sections, field)
		}
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
