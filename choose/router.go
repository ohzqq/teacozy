package choose

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	Choices []string
	//Selected    map[int]struct{}
	numSelected int
	width       int
	height      int
	header      string
	item.Items
}

type ChooseProps struct {
	item.Items
	//Selected   map[int]struct{}
	ToggleItem func(int)
	Height     int
	Width      int
}

func NewRouter(choices ...string) *List {
	list := &List{
		Choices:    choices,
		mainRouter: router.New(),
		height:     4,
		header:     "poot",
		//Selected:   make(map[int]struct{}),
	}
	list.Items = item.New(choices)

	w, h := util.TermSize()
	if list.height == 0 {
		list.height = h - 4
	}
	if list.width == 0 {
		list.width = w
	}

	return list
}

func (c *List) NewProps() ChooseProps {
	return ChooseProps{
		Width:  c.width,
		Height: c.height,
		Items:  c.Items,
		//Selected:   c.Selected,
		ToggleItem: c.ToggleSelection,
	}
}

func (c *List) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := New()

			return component, component.Init(c.NewProps())
		},
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewFilter()

			return component, component.Init(c.NewProps())
		},
	})
}

func (c *List) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ReturnSelectionsMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *List) Render(width, height int) string {
	view := c.mainRouter.Render(c.width, c.height)
	if c.header != "" {
		header := c.header + strings.Repeat(" ", c.width)
		view = lipgloss.JoinVertical(lipgloss.Left, header, view)
	}
	return view
}

func (m *List) ToggleSelection(idx int) {
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
	} else if m.numSelected < m.Items.Limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
}

func (cp ChooseProps) Visible(str ...string) []item.Item {
	if len(str) != 0 {
		return item.ExactMatches(str[0], cp.Items.Items)
	}
	return cp.Items.Items
}

func (m *List) Header(text string) *List {
	m.header = text
	return m
}

//func (m *Component) ChoiceMap(choices []map[string]string) *Component {
//  m.choiceMap = choices
//  return m
//}

func (m *List) Limit(l int) *List {
	m.Items.Limit = l
	return m
}

func (m *List) NoLimit() *List {
	return m.Limit(len(m.Choices))
}

func (m *List) Height(h int) *List {
	m.height = h
	return m
}

func (m *List) Width(w int) *List {
	m.width = w
	return m
}

func DefaultStyle() style.List {
	var s style.List
	s.Cursor = style.Cursor
	s.SelectedPrefix = style.Selected
	s.UnselectedPrefix = style.Unselected
	s.Text = style.Foreground
	s.Match = lipgloss.NewStyle().Foreground(color.Cyan())
	s.Header = lipgloss.NewStyle().Foreground(color.Purple())
	s.Prompt = style.Prompt
	return s
}
