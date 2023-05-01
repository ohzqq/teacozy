package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	headers map[string]ComponentInitializer
	mains   map[string]ComponentInitializer
	footers map[string]ComponentInitializer

	Header reactea.SomeComponent
	Main   reactea.SomeComponent
	Footer reactea.SomeComponent
}

type ComponentInitializer func() (reactea.SomeComponent, tea.Cmd)
type Opt func(c *Component)

const RoutePlaceholder = ":header/:main/:footer"

func New(main ComponentInitializer, opts ...Opt) *Component {
	c := &Component{
		headers: make(map[string]ComponentInitializer),
		mains:   make(map[string]ComponentInitializer),
		footers: make(map[string]ComponentInitializer),
	}
	c.headers["default"] = DefaultInitializer()
	c.mains["default"] = main
	c.footers["default"] = DefaultInitializer()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.initializeRoute())

	cmds = append(cmds, keys.ChangeRoute("default/default/default"))

	return tea.Batch(cmds...)
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case keys.ChangeRouteMsg:
		reactea.SetCurrentRoute(msg.Name)
		return nil

	case ChangeHeaderMsg:
		if header, ok := c.headers[msg.Name]; ok {
			c.Header, cmd = header()
			cmds = append(cmds, cmd)
		}
	case NewHeaderMsg:
		c.Header, cmd = msg.Init()
		cmds = append(cmds, cmd)

	case ChangeMainMsg:
		if main, ok := c.mains[msg.Name]; ok {
			c.Main, cmd = main()
			cmds = append(cmds, cmd)
		}
	case NewMainMsg:
		c.Main, cmd = msg.Init()
		cmds = append(cmds, cmd)

	case ChangeFooterMsg:
		if footer, ok := c.footers[msg.Name]; ok {
			c.Footer, cmd = footer()
			cmds = append(cmds, cmd)
		}
	case NewFooterMsg:
		c.Footer, cmd = msg.Init()
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}

		if msg.String() == "m" {
			return NewMain(func() (reactea.SomeComponent, tea.Cmd) {
				footer := reactea.Componentify[string](Renderer)
				return footer, footer.Init("new main")
			})
		}
	}

	if c.Header != nil {
		cmds = append(cmds, c.Header.Update(msg))
	}

	if c.Main != nil {
		cmds = append(cmds, c.Main.Update(msg))
	}

	if c.Footer != nil {
		cmds = append(cmds, c.Footer.Update(msg))
	}

	return tea.Batch(cmds...)
}

func (c *Component) AfterUpdate() tea.Cmd {
	if !reactea.WasRouteChanged() {
		return nil
	}

	if c.Header != nil {
		c.Header.Destroy()
	}
	c.Header = nil

	if c.Main != nil {
		c.Main.Destroy()
	}
	c.Main = nil

	if c.Footer != nil {
		c.Footer.Destroy()
	}
	c.Footer = nil

	return c.initializeRoute()
}

func (c *Component) initializeRoute() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if params, ok := reactea.RouteMatchesPlaceholder(reactea.CurrentRoute(), RoutePlaceholder); ok {
		if header, ok := c.headers[params["header"]]; ok {
			c.Header, cmd = header()
		}
		cmds = append(cmds, cmd)

		if main, ok := c.mains[params["main"]]; ok {
			c.Main, cmd = main()
		}
		cmds = append(cmds, cmd)

		if footer, ok := c.footers[params["footer"]]; ok {
			c.Footer, cmd = footer()
		}
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	var views []string

	if c.Header != nil {
		if header := c.Header.Render(w, h); header != "" {
			views = append(views, header)
		}
	}

	if c.Main != nil {
		if main := c.Main.Render(w, h); main != "" {
			views = append(views, main)
		}
	}

	if c.Footer != nil {
		if footer := c.Footer.Render(w, h); footer != "" {
			views = append(views, footer)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func DefaultInitializer() ComponentInitializer {
	return func() (reactea.SomeComponent, tea.Cmd) {
		return invisibleComponent(), nil
	}
}

func invisibleComponent() reactea.SomeComponent {
	return &struct {
		reactea.BasicComponent
		reactea.InvisibleComponent
	}{}
}

func (c *Component) AddHeader(name string, init ComponentInitializer) *Component {
	c.headers[name] = init
	return c
}

func (c *Component) AddMain(name string, init ComponentInitializer) *Component {
	c.mains[name] = init
	return c
}

func (c *Component) AddFooter(name string, init ComponentInitializer) *Component {
	c.footers[name] = init
	return c
}

func Header(name string, init ComponentInitializer) Opt {
	return func(c *Component) {
		c.headers[name] = init
	}
}

func Main(name string, init ComponentInitializer) Opt {
	return func(c *Component) {
		c.mains[name] = init
	}
}

func Footer(name string, init ComponentInitializer) Opt {
	return func(c *Component) {
		c.footers[name] = init
	}
}

func Headers(headers map[string]ComponentInitializer) Opt {
	return func(c *Component) {
		c.headers = headers
	}
}

func Mains(mains map[string]ComponentInitializer) Opt {
	return func(c *Component) {
		c.mains = mains
	}
}

func Footers(footers map[string]ComponentInitializer) Opt {
	return func(c *Component) {
		c.footers = footers
	}
}

type changeComponentMsg struct {
	Name string
}

type ChangeHeaderMsg changeComponentMsg

func ChangeHeader(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeHeaderMsg{Name: name}
	}
}

type ChangeMainMsg changeComponentMsg

func ChangeMain(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeMainMsg{Name: name}
	}
}

type ChangeFooterMsg changeComponentMsg

func ChangeFooter(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeFooterMsg{Name: name}
	}
}

type newComponentMsg struct {
	Init ComponentInitializer
}

type NewHeaderMsg newComponentMsg

func NewHeader(route ComponentInitializer) tea.Cmd {
	return func() tea.Msg {
		return NewHeaderMsg{Init: route}
	}
}

type NewMainMsg newComponentMsg

func NewMain(route ComponentInitializer) tea.Cmd {
	return func() tea.Msg {
		return NewMainMsg{Init: route}
	}
}

type NewFooterMsg newComponentMsg

func NewFooter(route ComponentInitializer) tea.Cmd {
	return func() tea.Msg {
		return NewFooterMsg{Init: route}
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}
