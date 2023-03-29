package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Filter struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ChooseProps]
	item.Items
	Cursor      int
	Choices     []string
	choiceMap   []map[string]string
	Input       textinput.Model
	Viewport    *viewport.Model
	FilterKeys  func(m *Filter) keys.KeyMap
	numSelected int
	limit       int
	aborted     bool
	quitting    bool
	header      string
	Placeholder string
	Prompt      string
	Width       int
	Height      int
	Style       style.List
}

type FilterKeys struct {
	Up               key.Binding
	Down             key.Binding
	ToggleItem       key.Binding
	Quit             key.Binding
	ReturnSelections key.Binding
	StopFiltering    key.Binding
}

func NewFilter(choices ...string) *Filter {
	tm := Filter{
		Choices: choices,
		Style:   DefaultStyle(),
		limit:   2,
		Prompt:  style.PromptPrefix,
		Height:  10,
	}

	w, h := util.TermSize()
	if tm.Height == 0 {
		tm.Height = h - 4
	}
	if tm.Width == 0 {
		tm.Width = w
	}
	tm.Input = textinput.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.PromptStyle = tm.Style.Prompt
	tm.Input.Placeholder = tm.Placeholder

	tm.header = "poot"

	return &tm
}

func (m *Filter) Run() []int {
	//p := tea.NewProgram(m)
	//if err := p.Start(); err != nil {
	//log.Fatal(err)
	//}
	//if m.quitting {
	//return []int{}
	//}
	return m.Chosen()
}

func (m *Filter) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.Height == 0 || m.Height > msg.Height {
			m.Viewport.Height = msg.Height - lipgloss.Height(m.Input.View())
		}

		m.Viewport.Width = msg.Width
	case ReturnSelectionsMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case UpMsg:
		m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor-1)
		if m.Cursor < m.Viewport.YOffset {
			m.Viewport.SetYOffset(m.Cursor)
		}
	case DownMsg:
		m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor+1)
		if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height {
			m.Viewport.LineDown(1)
		}
	case ToggleItemMsg:
		if m.Props().Limit == 1 {
			return nil
		}
		idx := m.Props().Visible()[m.Cursor].Index
		m.Props().ToggleItem(idx)
		cmds = append(cmds, DownCmd())
	case StopFilteringMsg:
		if m.Props().Limit == 1 {
			cmds = append(cmds, ToggleItemCmd())
		}

		m.Input.Reset()
		m.Input.Blur()
		reactea.SetCurrentRoute("default")
		return nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, filterKey.StopFiltering):
			cmds = append(cmds, StopFilteringCmd())
		case key.Matches(msg, filterKey.Up):
			cmds = append(cmds, UpCmd())
		case key.Matches(msg, filterKey.Down):
			cmds = append(cmds, DownCmd())
		case key.Matches(msg, filterKey.ToggleItem):
			cmds = append(cmds, ToggleItemCmd())
		case key.Matches(msg, filterKey.Quit):
			m.quitting = true
			cmds = append(cmds, ReturnSelectionsCmd())
		case key.Matches(msg, filterKey.ReturnSelections):
			cmds = append(cmds, ReturnSelectionsCmd())
		}
		m.Input, cmd = m.Input.Update(msg)
		m.Matches = m.Props().Visible(m.Input.Value())
		cmds = append(cmds, cmd)
	}

	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor)
	return tea.Batch(cmds...)
}

func (m *Filter) Render(w, h int) string {
	m.Viewport.Height = h
	m.Viewport.Width = w

	var s strings.Builder

	for i, match := range m.Matches {
		pre := "x"

		if match.Label != "" {
			pre = match.Label
		}

		switch {
		case i == m.Cursor:
			pre = match.Style.Cursor.Render(pre)
		default:
			if _, ok := m.Props().Selected[match.Index]; ok {
				pre = match.Style.Selected.Render(pre)
			} else if match.Label == "" {
				pre = strings.Repeat(" ", lipgloss.Width(pre))
			} else {
				pre = match.Style.Label.Render(pre)
			}
		}

		s.WriteString("[")
		s.WriteString(pre)
		s.WriteString("]")

		s.WriteString(match.RenderText())
		s.WriteRune('\n')
	}

	m.Viewport.SetContent(s.String())

	view := m.Input.View() + "\n" + m.Viewport.View()
	return view
}

func (tm *Filter) Init(props ChooseProps) tea.Cmd {
	tm.UpdateProps(props)
	tm.Matches = tm.Props().Visible()
	tm.Input.Width = tm.Width

	v := viewport.New(0, 0)
	tm.Viewport = &v
	tm.Input.Focus()

	return nil
}
