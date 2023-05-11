package teacozy

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type Page struct {
	reactea.BasicComponent

	header reactea.SomeComponent
	main   reactea.SomeComponent
	footer reactea.SomeComponent

	slug string
}

type Opt func(*Page)

func NewPage(title string, main reactea.SomeComponent, opts ...Opt) *Page {
	return &Page{
		slug: title,
		main: main,
	}
}

func (c Page) Slug() string {
	return c.slug
}

func (c Page) Header() reactea.SomeComponent {
	if head := c.header; head != nil {
		return head
	}
	return nil
}

func (c Page) Main() reactea.SomeComponent {
	return c.main
}

func (c Page) Footer() reactea.SomeComponent {
	if foot := c.footer; foot != nil {
		return foot
	}
	return nil
}

func (c *Page) SetHeader(head reactea.SomeComponent) *Page {
	c.header = head
	return c
}

func (c *Page) SetMain(main reactea.SomeComponent) *Page {
	c.main = main
	return c
}

func (c *Page) SetFooter(footer reactea.SomeComponent) *Page {
	c.footer = footer
	return c
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
