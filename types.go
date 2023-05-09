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
	Initialize(map[string]string) (*Page, tea.Cmd)
}

type Page struct {
	reactea.BasicComponent
	header reactea.SomeComponent
	main   reactea.SomeComponent
	footer reactea.SomeComponent
	width  int
	height int
}

type PageProps struct {
	Width  int
	Height int
}

func NewPage(w, h int) *Page {
	return &Page{
		width:  w,
		height: h,
	}
}

func (c *Page) Initialize(params map[string]string) (*Page, tea.Cmd) {
	return c, nil
}

func (c *Page) Header() reactea.SomeComponent {
	return c.header
}

func (c *Page) Main() reactea.SomeComponent {
	return c.main
}

func (c *Page) Footer() reactea.SomeComponent {
	return c.footer
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

func (c *Page) Update(msg tea.Msg) (*Page, tea.Cmd) {
	reactea.AfterUpdate(c)

	//var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return c, reactea.Destroy
		}
	}

	if c.header != nil {
		cmds = append(cmds, c.header.Update(msg))
	}

	if c.main != nil {
		cmds = append(cmds, c.main.Update(msg))
	}

	if c.footer != nil {
		cmds = append(cmds, c.footer.Update(msg))
	}

	return c, tea.Batch(cmds...)
}

func (c *Page) View() string {
	return c.Render(c.width, c.height)
}

func (c *Page) Render(w, h int) string {
	var views []string

	if c.header != nil {
		if head := c.header.Render(w, h); head != "" {
			views = append(views, head)
			h -= lipgloss.Height(head)
		}
	}

	var footer string
	if c.footer != nil {
		if f := c.footer.Render(w, h); f != "" {
			footer = f
			h -= lipgloss.Height(footer)
		}
	}

	if c.main != nil {
		if m := c.main.Render(w, h); m != "" {
			views = append(views, m)
		}
	}

	if footer != "" {
		views = append(views, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (c Page) Height() int {
	return c.height
}

func (c Page) Width() int {
	return c.width
}

func (c Page) Init() tea.Cmd {
	return nil
}
