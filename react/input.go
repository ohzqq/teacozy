package react

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Component struct {
	reactea.BasicComponent               // It implements AfterUpdate() for us, so we don't have to care!
	reactea.BasicPropfulComponent[Props] // It implements props backend - UpdateProps() and Props()

	textinput textinput.Model // Input for inputting name
}

type Props struct {
	SetText func(string) // SetText function for lifting state up
}

func New() *Component {
	return &Component{textinput: textinput.New()}
}

func (c *Component) Init(props Props) tea.Cmd {
	// Always derive props in Init()! If you are not replacing Init(),
	// reactea.BasicPropfulComponent will take care of it
	c.UpdateProps(props)

	return c.textinput.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			// Lifted state power! Woohooo
			c.Props().SetText(c.textinput.Value())

			// Navigate to displayname, please
			reactea.SetCurrentRoute("displayname")
			return nil
		}
	}

	var cmd tea.Cmd
	c.textinput, cmd = c.textinput.Update(msg)
	return cmd
}

// Here we are not using width and height, but you can!
func (c *Component) Render(int, int) string {
	return fmt.Sprintf("Enter your name: %s\nAnd press [ Enter ]", c.textinput.View())
}
