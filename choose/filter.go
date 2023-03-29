package choose

import (
	"strings"

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

func NewFilter(choices ...string) *Filter {
	tm := Filter{
		Choices:    choices,
		FilterKeys: FilterKeyMap,
		Style:      DefaultStyle(),
		limit:      2,
		Prompt:     style.PromptPrefix,
		Height:     10,
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
	case StopFilteringMsg:
		reactea.SetCurrentRoute("default")
	case tea.KeyMsg:
		for _, k := range GlobalsKeyMap(m) {
			if k.Matches(msg) {
				cmd = k.Cmd
				cmds = append(cmds, cmd)
			}
		}
		for _, k := range m.FilterKeys(m) {
			if k.Matches(msg) {
				cmd = k.Cmd
				cmds = append(cmds, cmd)
			}
		}
		m.Input, cmd = m.Input.Update(msg)
		m.Matches = item.ExactMatches(m.Input.Value(), m.Items.Items)
		if m.Input.Value() == "" {
			m.Items.Matches = m.Items.Items
		}
		cmds = append(cmds, cmd)
	}

	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor)
	return tea.Batch(cmds...)
}

func (m *Filter) CursorUp() {
	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor-1)
	if m.Cursor < m.Viewport.YOffset {
		m.Viewport.SetYOffset(m.Cursor)
	}
}

func (m *Filter) CursorDown() {
	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor+1)
	if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height {
		m.Viewport.LineDown(1)
	}
}

func (m *Filter) ToggleSelection() {
	idx := m.Matches[m.Cursor].Index
	m.Props().ToggleItem(idx)
	m.CursorDown()
}

func (m *Filter) Render(w, h int) string {
	m.Viewport.Height = h
	m.Viewport.Width = w

	return m.View()
}

func (m *Filter) View() string {
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
	tm.Items = props.Items
	tm.UpdateProps(props)
	return tm.init()
}

func (tm *Filter) init() tea.Cmd {
	tm.Input.Width = tm.Width

	v := viewport.New(tm.Width, tm.Height)
	tm.Viewport = &v
	tm.Input.Focus()

	return nil
}
