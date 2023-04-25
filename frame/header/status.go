package header

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	// How long status messages should stay visible. By default this is
	// 1 second.
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
	status                string
	Style                 Style
}

type Style struct {
	Header lipgloss.Style
	Status lipgloss.Style
}

type Props struct {
	Title string
}

func New() *Component {
	return &Component{
		Style:                 DefaultStyle(),
		StatusMessageLifetime: time.Second,
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case StatusMessageTimeoutMsg:
		c.HideStatusMessage()
	case keys.UpdateStatusMsg:
		return c.NewStatusMessage(msg.Status)
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	if s := c.status; s != "" {
		return c.Style.Status.Render(s)
	}
	if t := c.Props().Title; t != "" {
		return c.Style.Header.Render(t)
	}
	return ""
}

func DefaultStyle() Style {
	return Style{
		Header: lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
		Status: lipgloss.NewStyle().Foreground(color.Green()),
	}
}

// from: https://github.com/charmbracelet/bubbles/blob/v0.15.0/list/list.go#L290

type StatusMessageTimeoutMsg struct{}

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time. Note that this also returns a command.
func (c *Component) NewStatusMessage(s string) tea.Cmd {
	c.status = s
	if c.statusMessageTimer != nil {
		c.statusMessageTimer.Stop()
	}

	c.statusMessageTimer = time.NewTimer(c.StatusMessageLifetime)

	// Wait for timeout
	return func() tea.Msg {
		<-c.statusMessageTimer.C
		return StatusMessageTimeoutMsg{}
	}
}

func (m *Component) HideStatusMessage() {
	m.status = ""
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
}
