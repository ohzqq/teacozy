package filter

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered FilterState = iota // no filter set
	Filtering                     // user is actively setting a filter
)

type Model struct {
	item.Items
	Choices     []string
	choiceMap   []map[string]string
	Input       textinput.Model
	Viewport    *viewport.Model
	Paginator   paginator.Model
	FilterKeys  func(m *Model) keys.KeyMap
	numSelected int
	limit       int
	filterState FilterState
	aborted     bool
	quitting    bool
	header      string
	Placeholder string
	Prompt      string
	Width       int
	Height      int
	Style       style.List
}

func New(choices ...string) *Model {
	tm := Model{
		Choices:     choices,
		FilterKeys:  FilterKeyMap,
		filterState: Unfiltered,
		Style:       DefaultStyle(),
		limit:       1,
		Prompt:      style.PromptPrefix,
		Height:      10,
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

func (m *Model) CursorUp() {
	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor-1)
	if m.Cursor < m.Viewport.YOffset {
		m.Viewport.SetYOffset(m.Cursor)
	}
}

func (m *Model) CursorDown() {
	m.Cursor = clamp(0, len(m.Matches)-1, m.Cursor+1)
	if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height {
		m.Viewport.LineDown(1)
	}
}

func (m *Model) ToggleSelection() {
	idx := m.Matches[m.Cursor].Index
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.Items.Items[idx].Deselect()
		m.numSelected--
	} else if m.numSelected < m.limit {
		m.Items.Items[idx].Select()
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
	m.CursorDown()
}

func (m *Model) Current() item.Item {
	return m.Matches[m.Cursor]
}

func (m Model) View() string {
	var s strings.Builder

	s.WriteString(m.RenderItems(m.Cursor, m.Matches))

	m.Viewport.SetContent(s.String())

	view := m.Input.View() + "\n" + m.Viewport.View()
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
	tm.Input.Width = tm.Width

	v := viewport.New(tm.Width, tm.Height+4)
	tm.Viewport = &v

	tm.Input.Focus()

	return nil
}
