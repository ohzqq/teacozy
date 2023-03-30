package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	Choices     []string
	choiceMap   []map[string]string
	numSelected int
	width       int
	height      int
	header      string
	Input       textarea.Model
	inputValue  string
	itemIndex   int
	limit       int
	Items
}

type Props struct {
	Items
	Height int
	Width  int
}

func (cp Props) Visible(str ...string) []Item {
	if len(str) != 0 {
		return ExactMatches(str[0], cp.Items.Items)
	}
	return cp.Items.Items
}

type ChooseProps struct {
	Props
	ToggleItem func(int)
}

func New(choices ...string) *List {
	list := &List{
		Choices:    choices,
		choiceMap:  mapChoices(choices),
		mainRouter: router.New(),
	}
	list.Items = NewItems(list.choiceMap)

	w, h := util.TermSize()
	if list.height == 0 {
		list.height = h - 4
	}
	if list.width == 0 {
		list.width = w
	}

	list.Input = textarea.New()
	list.Input.ShowLineNumbers = false
	list.Input.SetWidth(list.width)

	return list
}

func (c *List) NewProps() Props {
	//items := NewChoiceMap(c.choiceMap)
	//items.Limit = c.limit
	return Props{
		Width:  c.width,
		Height: c.height,
		Items:  c.Items,
	}
}

func (c List) GatherItems() Items {
	matches := make([]Item, len(c.choiceMap))
	for i, option := range c.choiceMap {
		for label, val := range option {
			item := NewItem(val, i)
			item.Label = label
			matches[i] = item
		}
	}
	return Items{
		Items:    matches,
		Selected: c.Selected,
	}
}

func (c *List) NewChooseProps() ChooseProps {
	return ChooseProps{
		Props:      c.NewProps(),
		ToggleItem: c.ToggleSelection,
	}
}

func (c *List) NewFormProps() FormProps {
	return FormProps{
		Props:    c.NewProps(),
		EditItem: c.EditItem,
	}
}

func (c *List) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewChoice()

			return component, component.Init(c.NewChooseProps())
		},
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewFilter()

			return component, component.Init(c.NewChooseProps())
		},
		"form": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewForm()

			return component, component.Init(c.NewFormProps())
		},
		"field": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewField()

			println("field")
			props := FieldProps{
				//Item: c.Items.Items[c.itemIndex],
			}

			return component, component.Init(props)
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

func (m *List) Header(text string) *List {
	m.header = text
	return m
}

func (m *List) ChoiceMap(choices []map[string]string) *List {
	m.choiceMap = choices
	m.Items = m.GatherItems()
	return m
}

func mapChoices(c []string) []map[string]string {
	choices := make([]map[string]string, len(c))
	for i, val := range c {
		choices[i] = map[string]string{"": val}
	}
	return choices
}

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

func (m *List) SetInputValue(val string) {
	m.Input.SetValue(val)
}

func (m *List) EditItem(idx int) {
	m.itemIndex = idx
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
