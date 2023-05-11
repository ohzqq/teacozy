package pager

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
	Model paginator.Model

	cursor   int
	total    int
	start    int
	end      int
	ReadOnly bool
	Index    int
	KeyMap   keys.KeyMap
	Style    Style
	slug     string
}

type Props struct {
	Items      teacozy.Items
	Selected   map[int]struct{}
	PerPage    int
	Total      int
	ReadOnly   bool
	InputValue string
	SetCurrent func(int)
}

func New() *Component {
	m := &Component{
		KeyMap: keys.DefaultKeyMap(),
		Model:  paginator.New(),
		Style:  DefaultStyle(),
	}
	return m
}

func NewProps(items teacozy.Items) Props {
	p := Props{
		Items:    items,
		Selected: make(map[int]struct{}),
	}
	return p
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.Model.SetTotalPages(c.Props().Items.Len())
	c.SetPerPage(10)
	c.SliceBounds()
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case keys.PageUpMsg:
		m.cursor = clamp(m.cursor-m.Model.PerPage, 0, m.total-1)
		m.Model.PrevPage()
	case keys.PageDownMsg:
		m.cursor = clamp(m.cursor+m.Model.PerPage, 0, m.total-1)
		m.Model.NextPage()
	case keys.HalfPageUpMsg:
		m.HalfUp()
		if m.cursor < m.start {
			m.Model.PrevPage()
		}
	case keys.HalfPageDownMsg:
		m.HalfDown()
		if m.cursor >= m.end {
			m.Model.NextPage()
		}
	case keys.LineDownMsg:
		m.NextItem()
		if m.cursor >= m.end {
			m.Model.NextPage()
		}
	case keys.LineUpMsg:
		m.PrevItem()
		if m.cursor < m.start {
			m.Model.PrevPage()
		}
	case keys.TopMsg:
		m.cursor = 0
		m.Model.Page = 0
	case keys.BottomMsg:
		m.cursor = m.total - 1
		m.Model.Page = m.Model.TotalPages - 1
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return tea.Quit
		}
		for _, k := range m.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.SliceBounds()

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	var s strings.Builder
	h = h - 2

	// get matched items
	items := teacozy.ExactMatches(c.Props().InputValue, c.Props().Items)

	c.SetPerPage(h)

	// update the total items based on the filter results, this prevents from
	// going out of bounds
	c.SetTotal(len(items))

	for i, m := range items[c.Start():c.End()] {
		var cur bool
		if i == c.Highlighted() {
			c.SetCurrent(m.Index)
			cur = true
		}

		var sel bool
		if _, ok := c.Props().Selected[m.Index]; ok {
			sel = true
		}

		label := c.Props().Items.Label(m.Index)
		pre := c.PrefixText(label, sel, cur)
		style := c.PrefixStyle(label, sel, cur)

		// only print the prefix if it's a list or there's a label
		if !c.Props().ReadOnly || label != "" {
			s.WriteString(style.Render(pre))
		}

		// render the rest of the line
		text := lipgloss.StyleRunes(
			m.Str,
			m.MatchedIndexes,
			c.Style.Match,
			c.Style.Normal.Style,
		)

		s.WriteString(lipgloss.NewStyle().Render(text))
		s.WriteString("\n")
	}

	return s.String()
}
func (m Component) Cursor() int {
	m.cursor = clamp(m.cursor, 0, m.end-1)
	return m.cursor
}

func (m *Component) ResetCursor() {
	m.cursor = clamp(m.cursor, 0, m.end-1)
}

func (m Component) Len() int {
	return m.total
}

func (m Component) Current() int {
	return m.Index
}

func (m Component) Start() int {
	return m.start
}

func (m Component) End() int {
	return m.end
}

func (m *Component) SetCursor(n int) *Component {
	m.cursor = n
	return m
}

func (m *Component) SetKeyMap(km keys.KeyMap) {
	m.KeyMap = km
}

func (m *Component) SetCurrent(n int) {
	m.Index = n
}

func (m *Component) SetTotal(n int) *Component {
	m.total = n
	m.Model.SetTotalPages(n)
	m.SliceBounds()
	return m
}

func (m *Component) SetPerPage(n int) *Component {
	m.Model.PerPage = n
	m.Model.SetTotalPages(m.total)
	m.SliceBounds()
	return m
}

func (m Component) Highlighted() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.Model.PerPage {
			return i
		}
	}
	return 0
}

func (m *Component) SliceBounds() (int, int) {
	m.start, m.end = m.Model.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
	return m.start, m.end
}

func (m Component) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m Component) OnFirstItem() bool {
	return m.cursor == 0
}

func (m *Component) NextItem() {
	if !m.OnLastItem() {
		m.cursor++
	}
}

func (m *Component) PrevItem() {
	if !m.OnFirstItem() {
		m.cursor--
	}
}

func (m *Component) HalfDown() {
	if !m.OnLastItem() {
		m.cursor = m.cursor + m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Component) HalfUp() {
	if !m.OnFirstItem() {
		m.cursor = m.cursor - m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Component) DisableKeys() {
	m.KeyMap = keys.NewKeyMap(keys.Quit())
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
