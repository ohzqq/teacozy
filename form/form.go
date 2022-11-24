package form

import (
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
	Model    list.Model
	Input    textarea.Model
	Hash     map[string]string
	Changed  bool
	SaveForm SaveForm
	confirm  bool
	Info     *info.Info
	//Confirm  *teacozy.Menu
}

func New() Form {
	form := Form{
		Fields: NewFields(),
		Frame:  style.DefaultFrameStyle(),
	}

	return form
}

func (m *Form) SetFields(fields *Fields) {
	m.Fields = fields
}

func (m *Form) SetFormData(fd teacozy.FormData) {
	m.Fields = NewFields()
	m.Fields.SetData(fd)
}

func (m *Form) InitModel() {
	m.Model = list.New(
		m.Fields.Items(),
		itemDelegate(),
		m.Width(),
		m.Height(),
	)
	m.Model.SetShowStatusBar(false)
	m.Model.SetShowHelp(false)
	m.Model.Styles = style.ListStyles()

	m.Info = info.New(m.Fields)
}

func (m *Form) Start() {
	m.InitModel()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+V\n", m.Hash)
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
				if original := cur.Value(); original != val {
					cur.Set(val)
					m.Changed = true
					cmds = append(cmds, m.FieldChanged(cur))
				}
				m.Input.Blur()
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case m.confirm:
				switch {
				case key.Matches(msg, key.SaveAndExit):
					m.Info.SetTitle(`save and exit? y\n`)
					m.Info.Show()
					m.confirm = true
				case key.Matches(msg, Yes):
					m.Info.SetTitle("")
					m.SaveChanges()
					m.Info.Hide()
					m.confirm = false
				case key.Matches(msg, No):
					m.UndoChanges()
					m.Info.SetTitle("")
					m.Info.Hide()
					m.confirm = false
				case key.Matches(msg, key.EditField):
					m.confirm = false
				case key.Matches(msg, key.ToggleAllItems):
					m.confirm = !m.confirm
					m.Info.Toggle()
				case key.Matches(msg, key.PrevScreen):
					m.Info.Hide()
					m.confirm = false
				}
				var mod tea.Model
				mod, cmd = m.Info.Update(msg)
				m.Info = mod.(*info.Info)
				cmds = append(cmds, cmd)
			default:
				switch {
				case key.Matches(msg, key.ToggleAllItems):
					m.confirm = !m.confirm
					m.Info.Toggle()
				case key.Matches(msg, key.SaveAndExit):
					switch {
					case m.Changed:
						m.confirm = true
						m.Info.Show()
					default:
						//cmds = append(cmds, ExitFormCmd())
					}
				}
				m.Model, cmd = m.Model.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	//case UpdateStatusMsg:
	//cmds = append(cmds, m.Model.NewStatusMessage(msg.Msg))
	case tea.WindowSizeMsg:
		//m.width = msg.Width - 1
		//m.height = msg.Height - 1
		m.Frame.SetSize(msg.Width-1, msg.Height-2)
		m.Model.SetSize(msg.Width-1, msg.Height-2)
		m.Info.SetSize(msg.Width-1, msg.Height-2)
	case EditFormItemMsg:
		cur := m.Model.SelectedItem().(*Field)
		m.Input = textarea.New()
		m.Input.SetValue(cur.Value())
		m.Input.ShowLineNumbers = false
		m.Input.Focus()
	case teacozy.SetListItemMsg:
		idx := m.Model.Index()
		m.Model.SetItem(idx, msg.Item)
	//case ConfirmMenuMsg:
	//  if msg == true {
	//    m.SaveChanges()
	//    cmds = append(cmds, SaveFormCmd(m.SaveFormFunc))
	//  } else {
	//    m.UndoChanges()
	//    m.Changed = false
	//    cmds = append(cmds, ExitFormCmd())
	//  }
	case SaveAndExitFormMsg:
		cmds = append(cmds, msg.Save(m))
		cmds = append(cmds, FormChangedCmd())
		//cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m *Form) FieldChanged(item *Field) tea.Cmd {
	return func() tea.Msg {
		item.Changed()
		m.Changed = true
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

	if m.confirm {
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
