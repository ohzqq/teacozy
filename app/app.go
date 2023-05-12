package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/cmpnt"
	"github.com/ohzqq/teacozy/keys"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	*cmpnt.Pager
}

func New(opts ...cmpnt.Option) *App {
	c := &App{}
	c.Pager = cmpnt.New(opts...)

	return c
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.Pager.Init(reactea.NoProps{})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
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

	cmd = c.Pager.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
	return c.Pager.Render(w, h)
}
