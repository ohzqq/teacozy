package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Choose struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ChooseProps]
	item.Items
	Choices     []string
	choiceMap   []map[string]string
	Viewport    *viewport.Model
	Paginator   paginator.Model
	ListKeys    func(m *Choose) keys.KeyMap
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

type ChooseProps struct {
	item.Items
	ToggleItem func(int)
}

func (c *Component) NewProps() ChooseProps {
	return ChooseProps{
		Items:      item.New(c.Choices),
		ToggleItem: c.ToggleSelection,
	}
}

func New(choices ...string) *Choose {
	tm := Choose{
		Choices:  choices,
		ListKeys: ListKeyMap,
		Style:    DefaultStyle(),
		limit:    2,
		Prompt:   style.PromptPrefix,
		Height:   4,
	}

	w, h := util.TermSize()
	if tm.Height == 0 {
		tm.Height = h - 4
	}
	if tm.Width == 0 {
		tm.Width = w
	}

	tm.header = "poot"

	return &tm
}

func (m *Choose) Run() []int {
	//p := tea.NewProgram(m)
	//if err := p.Start(); err != nil {
	//log.Fatal(err)
	//}

	//if m.quitting {
	//return []int{}
	//}
	return m.Chosen()
}

func (m *Choose) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Make place in the view port if header is set
		if m.header != "" {
			m.Viewport.Height = m.Viewport.Height - lipgloss.Height(m.Style.Header.Render(m.header))
		}
		m.Viewport.Width = msg.Width
	case ReturnSelectionsMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		for _, k := range GlobalKeyMap(m) {
			if k.Matches(msg) {
				cmd = k.Cmd
				cmds = append(cmds, cmd)
			}
		}
		for _, k := range m.ListKeys(m) {
			if k.Matches(msg) {
				cmd = k.Cmd
				cmds = append(cmds, cmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m *Choose) CursorUp() {
	start, _ := m.Paginator.GetSliceBounds(len(m.Items.Matches))
	m.Cursor--
	if m.Cursor < 0 {
		m.Cursor = len(m.Items.Matches) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
	}
	if m.Cursor < start {
		m.Paginator.PrevPage()
	}
}

func (m *Choose) CursorDown() {
	_, end := m.Paginator.GetSliceBounds(len(m.Items.Matches))
	m.Cursor++
	if m.Cursor >= len(m.Items.Matches) {
		m.Cursor = 0
		m.Paginator.Page = 0
	}
	if m.Cursor >= end {
		m.Paginator.NextPage()
	}
}

func (m *Choose) ToggleSelection() {
	idx := m.Matches[m.Cursor].Index
	m.Props().ToggleItem(idx)
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
	} else if m.numSelected < m.limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
	m.CursorDown()
}

func (m *Choose) Render(w, h int) string {
	v := viewport.New(w, h+4)
	m.Viewport = &v
	return m.View()
}

func (m *Choose) View() string {
	var s strings.Builder

	start, end := m.Paginator.GetSliceBounds(len(m.Items.Matches))

	for i, match := range m.Items.Matches[start:end] {
		pre := "x"

		if match.Label != "" {
			pre = match.Label
		}

		switch {
		case i == m.Cursor%m.Height:
			pre = match.Style.Cursor.Render(pre)
		default:
			if _, ok := m.Selected[match.Index]; ok {
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

	var view string
	if m.Paginator.TotalPages <= 1 {
		view = s.String()
	} else if m.Paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", m.Height-m.Paginator.ItemsOnPage(len(m.Items.Matches))+1))
		s.WriteString("  " + m.Paginator.View())
	}

	view = s.String()
	if m.header != "" {
		header := m.Style.Header.Render(m.header + strings.Repeat(" ", m.Width))
		view = lipgloss.JoinVertical(lipgloss.Left, header, view)
	}

	m.Viewport.SetContent(view)
	view = m.Viewport.View()

	return view
}

//nolint:unparam
func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (tm *Choose) Init(props ChooseProps) tea.Cmd {
	tm.Items = props.Items
	tm.UpdateProps(props)
	return tm.init()
}

func (tm *Choose) init() tea.Cmd {

	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Dots
	tm.Paginator.ActiveDot = style.Subdued.Render(style.Bullet)
	tm.Paginator.InactiveDot = style.VerySubdued.Render(style.Bullet)
	tm.Paginator.SetTotalPages((len(tm.Items.Matches) + tm.Height - 1) / tm.Height)
	tm.Paginator.PerPage = tm.Height
	return nil
}
