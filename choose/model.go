package choose

import (
	"log"
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

type Model struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[item.Items]
	item.Items
	Choices     []string
	choiceMap   []map[string]string
	Viewport    *viewport.Model
	Paginator   paginator.Model
	ListKeys    func(m *Model) keys.KeyMap
	numSelected int
	limit       int
	aborted     bool
	quitting    bool
	cursor      int
	header      string
	Placeholder string
	Prompt      string
	Width       int
	Height      int
	Style       style.List
}

func New(choices ...string) *Model {
	tm := Model{
		Choices:  choices,
		ListKeys: ListKeyMap,
		Style:    DefaultStyle(),
		limit:    1,
		Prompt:   style.PromptPrefix,
		Height:   10,
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

func (m *Model) Run() []int {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	if m.quitting {
		return []int{}
	}
	return m.Chosen()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	return m, tea.Batch(cmds...)
}

func (m *Model) CursorUp() {
	start, _ := m.Paginator.GetSliceBounds(len(m.Items.Matches))
	m.Cursor--
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.Items.Matches) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
	}
	if m.cursor < start {
		m.Paginator.PrevPage()
	}
}

func (m *Model) CursorDown() {
	_, end := m.Paginator.GetSliceBounds(len(m.Items.Matches))
	m.Cursor++
	m.cursor++
	if m.cursor >= len(m.Items.Matches) {
		m.cursor = 0
		m.Paginator.Page = 0
	}
	if m.cursor >= end {
		m.Paginator.NextPage()
	}
}

func (m *Model) ToggleSelection() {
	idx := m.Matches[m.cursor].Index
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
	} else if m.numSelected < m.limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
	m.CursorDown()
}

func (m *Model) View() string {
	var s strings.Builder

	start, end := m.Paginator.GetSliceBounds(len(m.Items.Matches))

	for i, match := range m.Items.Matches[start:end] {
		pre := "x"

		if match.Label != "" {
			pre = match.Label
		}

		switch {
		case i == m.cursor%m.Height:
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

func (tm *Model) Init() tea.Cmd {
	tm.Items = item.New(tm.Choices)

	v := viewport.New(tm.Width, tm.Height+4)
	tm.Viewport = &v

	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Dots
	tm.Paginator.ActiveDot = style.Subdued.Render(style.Bullet)
	tm.Paginator.InactiveDot = style.VerySubdued.Render(style.Bullet)
	tm.Paginator.SetTotalPages((len(tm.Items.Matches) + tm.Height - 1) / tm.Height)
	tm.Paginator.PerPage = tm.Height
	return nil
}
