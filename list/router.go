package list

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/choose"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/filter"
	"github.com/ohzqq/teacozy/form"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	width  int
	height int
	header string
	footer string
	limit  int
	Props
}

type Props struct {
	*props.Items
	Footer func(string)
}

func (cp Props) Visible(matches ...string) []props.Item {
	if len(matches) != 0 {
		return props.ExactMatches(matches[0], cp.Items.Items)
	}
	return cp.Items.Items
}

func New(choices ...string) *List {
	list := &List{
		mainRouter: router.New(),
		Props: Props{
			Items: props.ItemSlice(choices),
		},
	}
	//items := NewItems(list.ChooseMap)
	//list.Items = items
	//list.Choices = choices
	//list.ChooseMap = MapChoices(choices)

	return list
}

func (c *List) NewProps() *props.Items {
	//p := props.Items{
	//Items: c.Items.Items,
	//Limit: c.Items.Limit,
	//}
	//p.Width = c.width
	//p.Height = c.height
	return c.Refresh()
}

func (c *List) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": choose.ChooseRouteInitializer(choose.ChooseProps{
			Items:      c.NewProps(),
			ToggleItem: c.ToggleSelection,
		}),
		"filter": filter.FilterRouteInitializer(filter.FilterProps{
			Items:      c.Items,
			ToggleItem: c.ToggleSelection,
		}),
		"form": form.FormRouteInitializer(form.FormProps{
			Items: c.Items,
			Save:  c.ChoiceMap,
		}),
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
	if c.footer != "" {
		f := style.Footer.Render(c.footer)
		view = lipgloss.JoinVertical(lipgloss.Left, view, f)
	}
	return view
}

func (m *List) Header(text string) *List {
	m.header = text
	return m
}

func (m *List) Footer(text string) {
	m.footer = text
}

func (m List) Chosen() []map[string]string {
	var chosen []map[string]string
	for _, c := range m.Items.Chosen() {
		chosen = append(chosen, m.Choices[c])
	}
	return chosen
}

func (m *List) Limit(l int) *List {
	m.Items.Limit = l
	return m
}

func (m *List) NoLimit() *List {
	return m.Limit(len(m.Choices))
}

func (m *List) Height(h int) *List {
	m.height = h + 2
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
