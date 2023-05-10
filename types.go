package teacozy

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type PlaceHolder string

func (ph PlaceHolder) Matches(route string) (map[string]string, bool) {
	return reactea.RouteMatchesPlaceholder(route, string(ph))
}

type PageComponent interface {
	Header() reactea.SomeComponent
	Main() reactea.SomeComponent
	Footer() reactea.SomeComponent
}

type Page struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[PageProps]
	header reactea.SomeComponent
	main   reactea.SomeComponent
	footer reactea.SomeComponent
	width  int
	height int
}

type PageProps struct {
	Page   PageComponent
	Width  int
	Height int
}

func NewPage(w, h int) *Page {
	return &Page{
		width:  w,
		height: h,
	}
}

func (c *Page) Init(props PageProps) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *Page) Header() reactea.SomeComponent {
	return c.Props().Page.Header()
}

func (c *Page) Main() reactea.SomeComponent {
	return c.Props().Page.Main()
}

func (c *Page) Footer() reactea.SomeComponent {
	return c.Props().Page.Footer()
}

func (c *Page) SetHeader(comp reactea.SomeComponent) {
	c.header = comp
}

func (c *Page) SetMain(comp reactea.SomeComponent) {
	c.main = comp
}

func (c *Page) SetFooter(comp reactea.SomeComponent) {
	c.footer = comp
}

func (c *Page) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	//var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	if c.Header() != nil {
		cmds = append(cmds, c.Header().Update(msg))
	}

	if c.Main() != nil {
		cmds = append(cmds, c.Main().Update(msg))
	}

	if c.Footer() != nil {
		cmds = append(cmds, c.Footer().Update(msg))
	}

	return tea.Batch(cmds...)
}

func (c *Page) View() string {
	return c.Render(c.Props().Width, c.Props().Height)
}

func (c *Page) Render(w, h int) string {
	var views []string

	if c.Header() != nil {
		if head := c.Header().Render(w, h); head != "" {
			views = append(views, head)
			h -= lipgloss.Height(head)
		}
	}

	var footer string
	if c.Footer() != nil {
		if f := c.Footer().Render(w, h); f != "" {
			footer = f
			h -= lipgloss.Height(footer)
		}
	}

	if c.Main() != nil {
		if m := c.Main().Render(w, h); m != "" {
			views = append(views, m)
		}
	}

	if footer != "" {
		views = append(views, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}
