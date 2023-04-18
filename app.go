package teacozy

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/confirm"
	"github.com/ohzqq/teacozy/edit"
	"github.com/ohzqq/teacozy/help"
	"github.com/ohzqq/teacozy/input"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/util"
	"github.com/ohzqq/teacozy/view"
	"golang.org/x/exp/maps"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]
	prevRoute  string

	Style
	width          int
	height         int
	footer         string
	title          string
	keyMap         keys.KeyMap
	editable       bool
	filterable     bool
	confirmChoices bool

	inputValue string
	filter     string

	confirm confirm.Props

	// How long status messages should stay visible. By default this is
	// 1 second.
	StatusMessageLifetime time.Duration
	statusMessage         string
	statusMessageTimer    *time.Timer
	status                string

	list        *list.Component
	Choices     item.Choices
	selected    map[int]struct{}
	numSelected int
	limit       int

	helpKeyMap item.Choices
}

type Style struct {
	Confirm lipgloss.Style
	Footer  lipgloss.Style
	Header  lipgloss.Style
	Status  lipgloss.Style
}

func New(opts ...Option) (*App, error) {
	a := &App{
		mainRouter:            router.New(),
		width:                 util.TermWidth(),
		height:                10,
		Style:                 DefaultStyle(),
		selected:              make(map[int]struct{}),
		limit:                 1,
		StatusMessageLifetime: time.Second,
		keyMap:                DefaultKeyMap(),
		list:                  list.New(),
	}

	for _, opt := range opts {
		opt(a)
	}

	if a.Choices.Len() == 0 {
		return a, fmt.Errorf("at least one choice is needed")
	}

	if a.limit == -1 {
		a.limit = a.Choices.Len()
	}

	return a, nil
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.list.Init(c.listProps())

	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := new(struct {
				reactea.BasicComponent
				reactea.InvisibleComponent
			})
			return component, nil
		},
		"form": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := list.New()
			return component, component.Init(c.listProps())
		},
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c.ResetInput()
			component := input.New()
			c.list.DefaultKeyMap()
			p := input.Props{
				Filter:   c.SetInput,
				ShowHelp: c.setHelp,
			}
			return component, component.Init(p)
		},
		"edit": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c.ResetInput()
			component := edit.New()
			c.list.SetKeyMap(keys.Global)
			c.SetInput(c.CurrentItem().Value())
			p := edit.Props{
				Save:     c.SetInput,
				Value:    c.inputValue,
				ShowHelp: c.setHelp,
			}
			return component, component.Init(p)
		},
		"confirm": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := confirm.New()
			p := c.confirm
			return component, component.Init(p)
		},
		"help": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := help.New()
			p := view.Props{
				Matches: item.ChoicesToItems(c.helpKeyMap),
			}
			return component, component.Init(help.Props{Props: p})
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	switch reactea.CurrentRoute() {
	case "":
		reactea.SetCurrentRoute("list")
		fallthrough
	case "list":
		c.list.VimKeyMap()
	}

	reactea.AfterUpdate(c)

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case keys.ReturnToListMsg:
		c.ResetInput()
		c.ResetFilter()
		cmds = append(cmds, keys.ChangeRoute("list"))

	case statusMessageTimeoutMsg:
		c.SetStatus("")
		c.hideStatusMessage()
		cmds = append(cmds, keys.ReturnToList)

	case confirm.GetConfirmationMsg:
		switch reactea.CurrentRoute() {
		case "list":
			if !c.confirmChoices {
				return reactea.Destroy
			}
			fallthrough
		default:
			c.confirm = msg.Props
			cmds = append(cmds, keys.ChangeRoute("confirm"))
		}

	case edit.SaveEditMsg:
		idx := c.list.CurrentItem()
		c.Choices.Set(idx, c.inputValue)
		cmds = append(cmds, keys.ReturnToList)

	case keys.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case "prev":
			return keys.ChangeRoute(c.prevRoute)
		}
		c.prevRoute = reactea.CurrentRoute()
		reactea.SetCurrentRoute(route)

	case tea.KeyMsg:
		for _, k := range c.keyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	switch reactea.CurrentRoute() {
	case "filter":
		c.filter = c.inputValue
		fallthrough
	case "list":
		cmds = append(cmds, c.list.Update(msg))
	}

	cmds = append(cmds, c.mainRouter.Update(msg))

	return tea.Batch(cmds...)
}

func (m *App) AfterUpdate() tea.Cmd {
	m.list.UpdateProps(m.listProps())
	return nil
}

func (c *App) Render(width, height int) string {
	w := c.Width()
	h := c.Height()

	var view []string

	if head := c.renderHeader(w, h); head != "" {
		h -= lipgloss.Height(head)
		view = append(view, head)
	}

	footer := c.renderFooter(w, h)
	if footer != "" {
		h -= lipgloss.Height(footer)
	}

	view = append(view, c.renderBody(w, h))

	if footer != "" {
		view = append(view, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, view...)
}

func (c App) renderBody(w, h int) string {
	var body string

	switch reactea.CurrentRoute() {
	case "help":
		fallthrough
	case "view":
		body = c.mainRouter.Render(w, h)
	default:
		body = c.list.Render(w, h)
	}

	return body
}

func (c App) renderHeader(w, h int) string {
	var header string
	if c.title != "" {
		header = c.Style.Header.Render(c.title)
	}

	if c.status != "" {
		header = c.Style.Status.Render(c.status)
	}

	switch reactea.CurrentRoute() {
	case "filter":
		header = c.mainRouter.Render(w, h)
	}
	return header
}

func (c App) renderFooter(w, h int) string {
	var footer string
	if c.footer != "" {
		footer = c.Style.Header.Render(c.footer)
	}
	switch reactea.CurrentRoute() {
	case "edit":
		fallthrough
	case "confirm":
		footer = c.mainRouter.Render(w, h)
	}
	return footer
}

func (c *App) listProps() list.Props {
	p := list.Props{
		Matches:     c.Choices.Filter(c.filter),
		Selected:    c.selected,
		ToggleItems: c.ToggleItems,
		Filterable:  c.filterable,
		Editable:    c.editable,
		ShowHelp:    c.setHelp,
	}
	return p
}

func Filter(search string, choices item.Choices) []item.Item {
	return choices.Filter(search)
}

func (m App) CurrentItem() item.Choice {
	return m.Choices[m.list.CurrentItem()]
}

func (m *App) ToggleItems(items ...int) {
	for _, idx := range items {
		if _, ok := m.selected[idx]; ok {
			delete(m.selected, idx)
			m.numSelected--
		} else if m.numSelected < m.limit {
			m.selected[idx] = struct{}{}
			m.numSelected++
		}
	}
}

func (m App) Chosen() []map[string]string {
	var chosen []map[string]string
	if len(m.selected) > 0 {
		for k := range m.selected {
			chosen = append(chosen, m.Choices[k])
		}
	}
	return chosen
}

func (m App) Selections() []int {
	return maps.Keys(m.selected)
}

func (c *App) setHelp(km []map[string]string) {
	c.helpKeyMap = item.MapToChoices(km)
}

func (c *App) ClearSelections() tea.Cmd {
	c.selected = make(map[int]struct{})
	return nil
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.Global
	return km
}

func (c *App) SetInput(value string) {
	c.inputValue = value
}

func (c *App) ResetInput() {
	c.inputValue = ""
}

func (c App) InputValue() string {
	return c.inputValue
}

func (c *App) SetFilter(s string) {
	c.filter = s
}

func (c *App) ResetFilter() {
	c.filter = ""
}

func (c App) FilterValue() string {
	return c.filter
}

func (c App) Height() int {
	return c.height
}

func (c App) Width() int {
	return c.width
}

func (c App) Limit() int {
	return c.limit
}

func (c *App) SetFooter(h string) *App {
	c.footer = h
	return c
}

func (c *App) SetStatus(h string) *App {
	c.status = h
	return c
}

func (c *App) SetTitle(h string) *App {
	c.title = h
	return c
}

func (c *App) SetLimit(l int) *App {
	c.limit = l
	return c
}

func DefaultStyle() Style {
	return Style{
		Confirm: lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()),
		Footer:  lipgloss.NewStyle().Foreground(color.Green()),
		Header:  lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
		Status:  lipgloss.NewStyle().Foreground(color.Green()),
	}
}

// from: https://github.com/charmbracelet/bubbles/blob/v0.15.0/list/list.go#L290

type statusMessageTimeoutMsg struct{}

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time. Note that this also returns a command.
func (m *App) NewStatusMessage(s string) tea.Cmd {
	m.status = s
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}

	m.statusMessageTimer = time.NewTimer(m.StatusMessageLifetime)

	// Wait for timeout
	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}

func (m *App) hideStatusMessage() {
	m.status = ""
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
}
