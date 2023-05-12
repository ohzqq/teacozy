package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/header"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/util"
)

type Page struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
	Model paginator.Model

	confirmChoices bool
	readOnly       bool

	width  int
	height int

	numSelected int
	limit       int
	CurrentItem int
	noLimit     bool

	cursor int
	total  int
	start  int
	end    int

	footer string

	choices teacozy.Items
	keyMap  keys.KeyMap
	Style   teacozy.Style

	Header *header.Component
	title  string
	header string

	help keys.KeyMap

	teacozy.State
}

type Props struct {
	teacozy.State
}

type AppStyle struct {
	Footer lipgloss.Style
}

func New(opts ...Option) *Page {
	c := &Page{
		width:  util.TermWidth(),
		height: util.TermHeight() - 2,
		limit:  10,
		Model:  paginator.New(),
		Style:  teacozy.DefaultStyle(),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.State = teacozy.NewProps(c.choices)
	c.State.SetCurrent = c.SetCurrent
	c.State.SetHelp = c.SetHelp
	c.State.ReadOnly = c.readOnly

	return c
}

func (c *Page) Init(reactea.NoProps) tea.Cmd {
	if c.noLimit {
		c.limit = c.choices.Len()
	}

	if !c.ReadOnly {
		c.AddKey(keys.Toggle().AddKeys(" "))
	}

	c.SetPerPage(c.height)
	c.SetTotal(c.choices.Len())
	c.SliceBounds()

	//c.Paginator = pagy.New(c.height, c.choices.Len())
	c.SetKeyMap(keys.VimKeyMap())

	c.Header = header.New()
	c.Header.Init(
		header.Props{
			Title: c.title,
		},
	)

	c.AddKey(keys.Help())

	return nil
}

func (c *Page) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case keys.PageUpMsg:
		c.cursor = clamp(c.cursor-c.Model.PerPage, 0, c.total-1)
		c.Model.PrevPage()
	case keys.PageDownMsg:
		c.cursor = clamp(c.cursor+c.Model.PerPage, 0, c.total-1)
		c.Model.NextPage()
	case keys.HalfPageUpMsg:
		c.HalfUp()
		if c.cursor < c.start {
			c.Model.PrevPage()
		}
	case keys.HalfPageDownMsg:
		c.HalfDown()
		if c.cursor >= c.end {
			c.Model.NextPage()
		}
	case keys.LineDownMsg:
		c.NextItem()
		if c.cursor >= c.end {
			c.Model.NextPage()
		}
	case keys.LineUpMsg:
		c.PrevItem()
		if c.cursor < c.start {
			c.Model.PrevPage()
		}
	case keys.TopMsg:
		c.cursor = 0
		c.Model.Page = 0
	case keys.BottomMsg:
		c.cursor = c.total - 1
		c.Model.Page = c.Model.TotalPages - 1
	case keys.ShowHelpMsg:
		cmds = append(cmds, keys.ChangeRoute("help"))
		//help := NewProps(c.help)
		//help.SetName("help")
		//return ChangeRoute(&help)

	case keys.UpdateItemMsg:
		return msg.Cmd(c.Current())

	case keys.ToggleItemsMsg, keys.ToggleItemMsg:
		c.ToggleItems(c.Current())
		cmds = append(cmds, keys.LineDown)

	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}

		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	c.SliceBounds()
	//c.Paginator, cmd = c.Paginator.Update(msg)
	//cmds = append(cmds, cmd)

	cmd = c.Header.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Page) Render(w, h int) string {
	var s strings.Builder
	//h = h - 2
	h = c.height - 2

	// get matched items
	items := c.ExactMatches(c.Search)

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
		if _, ok := c.Selected[m.Index]; ok {
			sel = true
		}

		label := c.Items.Label(m.Index)
		pre := c.State.PrefixText(label, sel, cur)
		style := c.State.PrefixStyle(label, sel, cur)

		// only print the prefix if it's a list or there's a label
		if !c.ReadOnly || label != "" {
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

func (c Page) renderHeader(w, h int) string {
	return c.Header.Render(w, h)
}

func (c Page) renderFooter(w, h int) string {
	var footer string

	//footer = fmt.Sprintf(
	//"cur route %v, per %v",
	//reactea.CurrentRoute(),
	//c.router.PrevRoute,
	//)

	//if c.footer != "" {
	//  footer = c.Style.Footer.Render(c.footer)
	//}

	return footer
}

func (c Page) ToggleItem() {
	c.ToggleItems(c.Current())
}

func (c *Page) AddKey(k *keys.Binding) *Page {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (c *Page) ToggleItems(items ...int) {
	for _, idx := range items {
		c.CurrentItem = idx
		if _, ok := c.Selected[idx]; ok {
			delete(c.Selected, idx)
			c.numSelected--
		} else if c.numSelected < c.limit {
			c.Selected[idx] = struct{}{}
			c.numSelected++
		}
	}
}

func (m Page) Chosen() []map[string]string {
	var chosen []map[string]string
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			l := m.choices.Label(k)
			v := m.choices.String(k)
			chosen = append(chosen, map[string]string{l: v})
		}
	}
	return chosen
}

func (m *Page) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *Page) Current() int {
	return m.CurrentItem
}

func (c *Page) SetWidth(n int) *Page {
	c.width = n
	return c
}

func (c *Page) SetHelp(km keys.KeyMap) {
	c.help = km
}

func (c *Page) SetHeight(n int) *Page {
	c.height = n
	return c
}

func (c *Page) SetSize(w, h int) *Page {
	c.width = w
	c.height = h
	return c
}

func (c *Page) SetHeader(h string) {
	c.header = h
}

func (m Page) Cursor() int {
	m.cursor = clamp(m.cursor, 0, m.end-1)
	return m.cursor
}

func (m *Page) ResetCursor() {
	m.cursor = clamp(m.cursor, 0, m.end-1)
}

func (m Page) Len() int {
	return m.total
}

func (m Page) Start() int {
	return m.start
}

func (m Page) End() int {
	return m.end
}

func (m *Page) SetCursor(n int) *Page {
	m.cursor = n
	return m
}

func (m *Page) SetKeyMap(km keys.KeyMap) {
	m.keyMap = km
}

func (m *Page) SetTotal(n int) *Page {
	m.total = n
	m.Model.SetTotalPages(n)
	m.SliceBounds()
	return m
}

func (m *Page) SetPerPage(n int) *Page {
	m.Model.PerPage = n
	m.Model.SetTotalPages(m.total)
	m.SliceBounds()
	return m
}

func (m Page) Highlighted() int {
	for i := 0; i < m.end; i++ {
		if i == m.cursor%m.Model.PerPage {
			return i
		}
	}
	return 0
}

func (m *Page) SliceBounds() (int, int) {
	m.start, m.end = m.Model.GetSliceBounds(m.total)
	m.start = clamp(m.start, 0, m.total-1)
	return m.start, m.end
}

func (m Page) OnLastItem() bool {
	return m.cursor == m.total-1
}

func (m Page) OnFirstItem() bool {
	return m.cursor == 0
}

func (m *Page) NextItem() {
	if !m.OnLastItem() {
		m.cursor++
	}
}

func (m *Page) PrevItem() {
	if !m.OnFirstItem() {
		m.cursor--
	}
}

func (m *Page) HalfDown() {
	if !m.OnLastItem() {
		m.cursor = m.cursor + m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Page) HalfUp() {
	if !m.OnFirstItem() {
		m.cursor = m.cursor - m.Model.PerPage/2 - 1
		m.cursor = clamp(m.cursor, 0, m.total-1)
	}
}

func (m *Page) DisableKeys() {
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
