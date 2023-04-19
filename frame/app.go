package frame

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/util"
	"github.com/ohzqq/teacozy/view"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	choices      item.Choices
	KeyMap       keys.KeyMap
	Cursor       int
	TotalMatches int
	search       string
	Viewport     viewport.Model
	Start        int
	End          int
}

func New(c []string) *App {
	return &App{
		mainRouter: router.New(),
		choices:    item.SliceToChoices(c),
		Viewport:   viewport.New(util.TermWidth(), 10),
		KeyMap:     view.DefaultKeyMap(),
	}
}

func (m *App) Init(reactea.NoProps) tea.Cmd {
	return m.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := view.NewComponent()
			p := view.NewProps(m.choices)
			p.SetCursor = m.SetCursor
			p.TotalMatches = m.SetTotalMatches
			return component, component.Init(p)
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}
	cmd = c.mainRouter.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
	view := c.mainRouter.Render(util.TermWidth(), 10)

	view += fmt.Sprintf("\n %d", c.Cursor)

	return view
}

func (c *App) SetCursor(n int) {
	c.Cursor = n
}

func (c *App) SetTotalMatches(n int) {
	c.TotalMatches = n
}

func (m *App) UpdateItems() {

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height

	if m.Cursor >= 0 {
		m.Start = clamp(m.Cursor-m.Height(), 0, m.Cursor)
	} else {
		m.Start = 0
		m.SetCursor(0)
	}

	m.End = clamp(m.Cursor+m.Height(), m.Cursor, m.TotalMatches)
	if m.Cursor > m.End {
		m.SetCursor(clamp(m.Cursor+m.Height(), m.Cursor, m.TotalMatches-1))
	}

	l := item.NewList()
	//l.Init(m.ItemProps())
	m.Viewport.SetContent(l.Render(m.Width(), m.Height()))
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *App) MoveUp(n int) {
	m.SetCursor(clamp(m.Cursor-n, 0, m.TotalMatches-1))
	m.UpdateItems()
	switch {
	case m.Start == 0:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Cursor))
	case m.Start < m.Height():
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Height())
	}
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *App) MoveDown(n int) {
	m.SetCursor(clamp(m.Cursor+n, 0, m.TotalMatches-1))
	m.UpdateItems()
	switch {
	case m.End == m.TotalMatches:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Height()))
	case m.Cursor > (m.End-m.Start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
	case m.Viewport.YOffset > 1:
	case m.Cursor > m.Viewport.YOffset+m.Height()-1:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *App) GotoTop() {
	m.MoveUp(m.Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *App) GotoBottom() {
	m.MoveDown(m.TotalMatches)
}

// Height returns the viewport height of the list.
func (m App) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the list.
func (m App) Width() int {
	return m.Viewport.Width
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
