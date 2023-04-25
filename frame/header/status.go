package status

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	// How long status messages should stay visible. By default this is
	// 1 second.
	StatusMessageLifetime time.Duration
	statusMessage         string
	statusMessageTimer    *time.Timer
	Style                 lipgloss.Style
}

type Props struct {
	SetStatus func(string)
	Status    string
}

func New() *Component {
	return &Component{
		Style: lipgloss.NewStyle().Foreground(color.Pink()),
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *Component) Render(w, h int) string {
	return c.Style.Render(c.Props().Status)
}

// from: https://github.com/charmbracelet/bubbles/blob/v0.15.0/list/list.go#L290

type StatusMessageTimeoutMsg struct{}

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time. Note that this also returns a command.
func (m *Component) NewStatusMessage(s string) tea.Cmd {
	m.SetStatus(s)
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

func (m *Component) HideStatusMessage() {
	m.SetStatus("")
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
}
