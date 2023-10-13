package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type StatusMsg struct {
	// How long status messages should stay visible. By default this is
	// 1 second.
	StatusMessageLifetime time.Duration

	statusMessage      string
	statusMessageTimer *time.Timer
}

// Status message code from https://github.com/charmbracelet/bubbles/blob/v0.16.1/list/list.go

type statusMessageTimeoutMsg struct{}

func (m *Model) hideStatusMessage() {
	m.statusMessage = ""
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
}

func (m *Model) NewStatusMessage(s string) tea.Cmd {
	m.statusMessage = s
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
