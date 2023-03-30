package filter

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Filter struct {
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
		limit:      1,
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
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	if m.quitting {
		return []int{}
	}
	return m.Chosen()
}

func (m *Filter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.Height == 0 || m.Height > msg.Height {
			m.Viewport.Height = msg.Height - lipgloss.Height(m.Input.View())
		}

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
	return m, tea.Batch(cmds...)
}

func (m *Filter) CursorUp() {
	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor-1)
	if m.Cursor < m.Viewport.YOffset {
		m.Viewport.SetYOffset(m.Cursor)
	} else {
		m.Viewport.GotoBottom()
	}
}

func (m *Filter) CursorDown() {
	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor+1)
	if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height {
		m.Viewport.LineDown(1)
	} else {
		m.Viewport.GotoTop()
	}
}

func (m *Filter) ToggleSelection() {
	idx := m.Matches[m.Cursor].Index
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
	} else if m.numSelected < m.limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
	m.CursorDown()
}

func (m Filter) View() string {
	var s strings.Builder

	s.WriteString(m.RenderItems(m.Cursor, m.Matches))

	m.Viewport.SetContent(s.String())

	view := m.Input.View() + "\n" + m.Viewport.View()
	return view
}

//nolint:unparam
func clamp(min, max, val int) int {
	if val < min {
		return max
	}
	if val > max {
		return min
	}
	return val
}

func (tm *Filter) Init() tea.Cmd {
	tm.Items = item.New(tm.Choices)
	tm.Input.Width = tm.Width

	v := viewport.New(tm.Width, tm.Height+4)
	tm.Viewport = &v

	tm.Input.Focus()

	return nil
}
