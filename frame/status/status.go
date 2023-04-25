package frame

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Status struct {
	// How long status messages should stay visible. By default this is
	// 1 second.
	StatusMessageLifetime time.Duration
	statusMessage         string
	statusMessageTimer    *time.Timer
	status                string
}

type StatusProps struct {
	SetStatus func(string)
}

// from: https://github.com/charmbracelet/bubbles/blob/v0.15.0/list/list.go#L290

type statusMessageTimeoutMsg struct{}

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time. Note that this also returns a command.
func (m *Status) NewStatusMessage(s string) tea.Cmd {
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

func (m *Status) hideStatusMessage() {
	m.status = ""
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
}
