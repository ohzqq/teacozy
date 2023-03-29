package choose

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/util"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	Choices     []string
	Selected    map[int]struct{}
	numSelected int
	Width       int
	Height      int
	header      string
	item.Items
}

type ChooseProps struct {
	item.Items
	Selected   map[int]struct{}
	ToggleItem func(int)
	Height     int
	Width      int
}

func (cp ChooseProps) Visible(str ...string) []item.Item {
	if len(str) != 0 {
		return ExactMatches(str[0], cp.Items.Items)
	}
	return cp.Items.Items
}

func NewRouter(choices ...string) *Component {
	list := &Component{
		Choices:    choices,
		mainRouter: router.New(),
		Height:     4,
		header:     "poot",
		Selected:   make(map[int]struct{}),
	}
	list.Items = item.New(choices)
	list.Limit = 2

	w, h := util.TermSize()
	if list.Height == 0 {
		list.Height = h - 4
	}
	if list.Width == 0 {
		list.Width = w
	}

	return list
}

func (c *Component) NewProps() ChooseProps {
	return ChooseProps{
		Width:      c.Width,
		Height:     c.Height,
		Items:      c.Items,
		Selected:   c.Selected,
		ToggleItem: c.ToggleSelection,
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
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

func (c *Component) Update(msg tea.Msg) tea.Cmd {
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

func (c *Component) Render(width, height int) string {
	view := c.mainRouter.Render(c.Width, c.Height)
	if c.header != "" {
		header := c.header + strings.Repeat(" ", c.Width)
		view = lipgloss.JoinVertical(lipgloss.Left, header, view)
	}
	return view
}

func (m *Component) ToggleSelection(idx int) {
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
	} else if m.numSelected < m.Limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
}

func (m *Component) SetCursor(cur int) {
	m.Cursor = cur
}

func ExactMatches(search string, choices []item.Item) []item.Item {
	matches := []item.Item{}
	for _, choice := range choices {
		search = strings.ToLower(search)
		matchedString := strings.ToLower(choice.Str)

		index := strings.Index(matchedString, search)
		if index >= 0 {
			for s := range search {
				choice.MatchedIndexes = append(choice.MatchedIndexes, index+s)
			}
			matches = append(matches, choice)
		}
	}

	return matches
}
