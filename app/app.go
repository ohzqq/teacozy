package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/app/confirm"
	"github.com/ohzqq/teacozy/app/edit"
	"github.com/ohzqq/teacozy/app/input"
	"github.com/ohzqq/teacozy/app/item"
	"github.com/ohzqq/teacozy/app/list"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/maps"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]
	PrevRoute  string

	Style          style.App
	width          int
	height         int
	footer         string
	title          string
	keyMap         keys.KeyMap
	editable       bool
	filterable     bool
	confirmChoices bool

	inputValue string
	search     string
	edit       string

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
}

type Option func(*App)

func New(opts ...Option) (*App, error) {
	a := &App{
		mainRouter:            router.New(),
		width:                 util.TermWidth(),
		height:                10,
		Style:                 style.DefaultAppStyle(),
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

func (c *App) listProps() list.Props {
	p := list.Props{
		Matches:     c.Choices.Filter(c.search),
		Selected:    c.selected,
		ToggleItems: c.ToggleItems,
		Filterable:  c.filterable,
		Editable:    c.editable,
	}
	return p
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
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c.ResetInput()
			component := input.New()
			c.list.DefaultKeyMap()
			p := input.Props{Filter: c.SetInput}
			return component, component.Init(p)
		},
		"edit": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c.ResetInput()
			component := edit.New()
			c.list.SetKeyMap(keys.Global)
			c.SetInput(c.CurrentItem().Value())
			p := edit.Props{
				Save:  c.SetInput,
				Value: c.inputValue,
			}
			return component, component.Init(p)
		},
		"confirm": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := confirm.New()
			p := c.confirm
			return component, component.Init(p)
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
		reactea.SetCurrentRoute("list")
		return nil

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

	case edit.ConfirmEditMsg:
		if c.inputValue != c.CurrentItem().Value() {
			cmd := confirm.Action("save edit?", edit.Save)
			cmds = append(cmds, cmd)
		}
	case edit.SaveEditMsg:
		idx := c.list.CurrentItem()
		c.Choices.Set(idx, c.inputValue)
		cmds = append(cmds, keys.ReturnToList)
	case edit.StartEditingMsg:
		cmds = append(cmds, message.ChangeRoute("edit"))

	case keys.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case "prev":
			route = c.PrevRoute
		case "help":
			//p := c.NewProps(KeymapToProps(c.help))
			//p.Height = c.Items.Height
			//p.Width = c.Items.Width
			//c.Routes["help"] = view.New().Initializer(p)
		}
		c.ChangeRoute(route)

	case message.QuitMsg:
		return reactea.Destroy

	case tea.KeyMsg:
		for _, k := range c.keyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	switch reactea.CurrentRoute() {
	case "help":
	case "filter":
		c.search = c.inputValue
		fallthrough
	case "list":
		cmds = append(cmds, c.list.Update(msg))
	}

	cmds = append(cmds, c.mainRouter.Update(msg))

	return tea.Batch(cmds...)
}

func (m App) CurrentItem() item.Choice {
	return m.Choices[m.list.CurrentItem()]
}

func (m *App) AfterUpdate() tea.Cmd {
	m.list.UpdateProps(m.listProps())
	return nil
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

	list := c.list.Render(w, h)
	view = append(view, list)

	if footer != "" {
		view = append(view, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, view...)
}

func (c App) renderHeader(w, h int) string {
	var header string
	if c.title != "" {
		header = c.Style.Header.Render(c.title)
	}

	if c.status != "" {
		header = lipgloss.NewStyle().Foreground(color.Green()).Render(c.status)
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

func (c *App) ChangeRoute(r string) {
	if p := reactea.CurrentRoute(); p == "" {
		c.PrevRoute = "default"
	} else {
		c.PrevRoute = p
	}
	reactea.SetCurrentRoute(r)
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

func WithSlice[E any](c []E) Option {
	return func(a *App) {
		a.Choices = item.ChoiceSliceToMap(c)
	}
}

func WithMap(c []map[string]string) Option {
	return func(a *App) {
		a.Choices = item.ChoiceMap(c)
	}
}

func NoLimit() Option {
	return func(a *App) {
		a.limit = -1
	}
}

func WithLimit(l int) Option {
	return func(a *App) {
		a.limit = l
	}
}

func WithTitle(t string) Option {
	return func(a *App) {
		a.title = t
	}
}
func ConfirmChoices() Option {
	return func(a *App) {
		a.confirmChoices = true
	}
}

func Editable() Option {
	return func(a *App) {
		a.editable = true
	}
}

func WithFilter() Option {
	return func(a *App) {
		a.filterable = true
	}
}

func Filter(search string, choices item.Choices) []item.Item {
	return choices.Filter(search)
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
	c.search = s
}

func (c *App) ResetFilter() {
	c.search = ""
}

func (c App) FilterValue() string {
	return c.search
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

func (c *App) ClearSelections() tea.Cmd {
	c.selected = make(map[int]struct{})
	return nil
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
