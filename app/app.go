package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/cmpnt"
	"github.com/ohzqq/teacozy/keys"
)

type Page struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	*cmpnt.Page
}

type AppStyle struct {
	Footer lipgloss.Style
}

func New(opts ...cmpnt.Option) *Page {
	c := &Page{
		//width:  util.TermWidth(),
		//height: util.TermHeight() - 2,
		//limit:  10,
		//Page: cmpnt
		//Model: paginator.New(),
		//Style: pager.DefaultStyle(),
	}
	c.Page = cmpnt.New(opts...)

	//c.State = teacozy.NewProps(c.choices)
	//c.State.SetCurrent = c.SetCurrent
	//c.State.SetHelp = c.SetHelp
	//c.State.ReadOnly = c.readOnly

	return c
}

func (c *Page) Init(reactea.NoProps) tea.Cmd {
	return c.Page.Init(cmpnt.PageProps{Items: c.Choices})
}

func (c *Page) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case keys.ShowHelpMsg:
		cmds = append(cmds, keys.ChangeRoute("help"))
		//help := NewProps(c.help)
		//help.SetName("help")
		//return ChangeRoute(&help)

	case keys.UpdateItemMsg:
		return msg.Cmd(c.Current())

	case keys.ToggleItemsMsg, keys.ToggleItemMsg:
		c.ToggleItems(c.Current())
		cmds = append(cmds, keys.LineDown)

	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}

		//for _, k := range c.KeyMap.Keys() {
		//  if key.Matches(msg, k.Binding) {
		//    cmds = append(cmds, k.TeaCmd)
		//  }
		//}

	}

	cmd = c.Page.Update(msg)
	cmds = append(cmds, cmd)

	cmd = c.Header.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Page) Render(w, h int) string {
	return c.Page.Render(w, h)
}

func (c Page) renderHeader(w, h int) string {
	return c.Header.Render(w, h)
}

func (c Page) renderFooter(w, h int) string {
	var footer string

	//footer = fmt.Sprintf(
	//"cur route %v, per %v",
	//reactea.CurrentRoute(),
	//c.router.PrevRoute,
	//)

	//if c.footer != "" {
	//  footer = c.Style.Footer.Render(c.footer)
	//}

	return footer
}
