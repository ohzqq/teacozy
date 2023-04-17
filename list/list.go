package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/confirm"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
)

type Option func(*Component)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Cursor int
	KeyMap keys.KeyMap

	Viewport viewport.Model
	start    int
	end      int
}

type Props struct {
	Editable    bool
	Filterable  bool
	Matches     []item.Item
	Selected    map[int]struct{}
	ToggleItems func(...int)
	ShowHelp    func(keys.KeyMap)
}

func New() *Component {
	m := Component{
		Cursor: 0,
	}
	m.DefaultKeyMap()

	return &m
}

func (m *Component) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	m.Viewport = viewport.New(0, 0)
	m.UpdateItems()
	return nil
}

func (m *Component) AfterUpdate() tea.Cmd {
	m.UpdateItems()
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case keys.ReturnSelectionsMsg:
		if reactea.CurrentRoute() == "list" {
			return confirm.GetConfirmation("Accept selected?", AcceptChoices)
		}

	case keys.ShowHelpMsg:
		m.Props().ShowHelp(m.KeyMap)
		cmds = append(cmds, keys.ChangeRoute("help"))

	case keys.ToggleItemMsg:
		cur := m.Props().Matches[m.Cursor].Index
		m.Props().ToggleItems(cur)
		m.MoveDown(1)

	case keys.PageUpMsg:
		m.MoveUp(m.Height())
	case keys.PageDownMsg:
		m.MoveDown(m.Height())
	case keys.HalfPageUpMsg:
		m.MoveUp(m.Height() / 2)
	case keys.HalfPageDownMsg:
		m.MoveDown(m.Height() / 2)
	case keys.LineDownMsg:
		m.MoveDown(1)
	case keys.LineUpMsg:
		m.MoveUp(1)
	case keys.TopMsg:
		m.GotoTop()
	case keys.BottomMsg:
		m.GotoBottom()
	case tea.KeyMsg:
		if reactea.CurrentRoute() == "list" {
			if m.Props().Editable {
				if k := keys.Edit(); key.Matches(msg, k.Binding) {
					return k.TeaCmd
				}
			}
			if m.Props().Filterable {
				if k := keys.Filter(); key.Matches(msg, k.Binding) {
					return k.TeaCmd
				}
			}
		}
		for _, k := range m.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m *Component) Render(w, h int) string {
	m.SetWidth(w)
	m.SetHeight(h)
	m.UpdateItems()
	return m.Viewport.View()
}

// UpdateItems updates the list content based on the previously defined
// columns and rows.
func (m *Component) UpdateItems() {
	items := make([]string, 0, len(m.Props().Matches))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.Cursor >= 0 {
		m.start = clamp(m.Cursor-m.Height(), 0, m.Cursor)
	} else {
		m.start = 0
		m.SetCursor(0)
	}
	m.end = clamp(m.Cursor+m.Height(), m.Cursor, len(m.Props().Matches))

	if m.Cursor > m.end {
		m.SetCursor(clamp(m.Cursor+m.Height(), m.Cursor, len(m.Props().Matches)-1))
	}

	for i := m.start; i < m.end; i++ {
		items = append(items, m.renderItem(i))
	}

	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)
}

func (m *Component) renderItem(rowID int) string {
	item := m.Props().Matches[rowID]

	var s strings.Builder

	switch {
	case rowID == m.Cursor:
		item.Current = true
	case m.isSelected(rowID):
		item.Selected = true
	}

	s.WriteString(item.Render(m.Width(), m.Height()))

	return s.String()
}

func (m Component) isSelected(idx int) bool {
	_, ok := m.Props().Selected[m.Props().Matches[idx].Index]
	return ok
}

// CurrentItem returns the selected row.
// You can cast it to your own implementation.
func (m Component) CurrentItem() int {
	return m.Props().Matches[m.Cursor].Index
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *Component) MoveUp(n int) {
	m.SetCursor(clamp(m.Cursor-n, 0, len(m.Props().Matches)-1))
	m.UpdateItems()
	switch {
	case m.start == 0:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset, 0, m.Cursor))
	case m.start < m.Height():
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+n, 0, m.Cursor))
	case m.Viewport.YOffset >= 1:
		m.Viewport.YOffset = clamp(m.Viewport.YOffset+n, 1, m.Height())
	}
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *Component) MoveDown(n int) {
	m.SetCursor(clamp(m.Cursor+n, 0, len(m.Props().Matches)-1))
	m.UpdateItems()
	switch {
	case m.end == len(m.Props().Matches):
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Height()))
	case m.Cursor > (m.end-m.start)/2:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset-n, 1, m.Cursor))
	case m.Viewport.YOffset > 1:
	case m.Cursor > m.Viewport.YOffset+m.Height()-1:
		m.Viewport.SetYOffset(clamp(m.Viewport.YOffset+1, 0, 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *Component) GotoTop() {
	m.MoveUp(m.Cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Component) GotoBottom() {
	m.MoveDown(len(m.Props().Matches))
}

// SetWidth sets the width of the viewport of the list.
func (m *Component) SetWidth(w int) {
	m.Viewport.Width = w
	m.UpdateItems()
}

// SetHeight sets the height of the viewport of the list.
func (m *Component) SetHeight(h int) {
	m.Viewport.Height = h
	m.UpdateItems()
}

func (m *Component) commonKeys() keys.KeyMap {
	var km = keys.KeyMap{
		keys.PgUp(),
		keys.PgDown(),
		keys.Enter().
			WithHelp("return selections").
			Cmd(keys.ReturnSelections()),
	}
	return km
}

// SetKeyMap sets the keymap for the list.
func (m *Component) SetKeyMap(km keys.KeyMap) {
	m.KeyMap = m.commonKeys()
	m.KeyMap = append(m.KeyMap, km...)
}

func (m *Component) VimKeyMap() *Component {
	m.SetKeyMap(VimKeyMap())

	h := keys.Help().
		AddKeys("h").
		Cmd(keys.ShowHelp(m.KeyMap))
	m.KeyMap = append(m.KeyMap, h)

	return m
}

func (m *Component) DefaultKeyMap() *Component {
	m.SetKeyMap(DefaultKeyMap())

	h := keys.Help().Cmd(keys.ShowHelp(m.KeyMap))
	m.KeyMap = append(m.KeyMap, h)

	return m
}

// Height returns the viewport height of the list.
func (m Component) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the list.
func (m Component) Width() int {
	return m.Viewport.Width
}

// Cursor returns the index of the selected item.
func (m Component) GetCursor() int {
	return m.Cursor
}

// SetCursor sets the cursor position in the list.
func (m *Component) SetCursor(n int) {
	m.Cursor = clamp(n, 0, len(m.Props().Matches)-1)
}

func AcceptChoices(accept bool) tea.Cmd {
	if accept {
		return reactea.Destroy
	}
	return keys.ReturnToList
}

func DefaultKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.ToggleItem(),
		keys.Up(),
		keys.Down(),
		keys.HalfPgUp(),
		keys.HalfPgDown(),
		keys.Home(),
		keys.End(),
		keys.Quit(),
	}
}

func VimKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.ToggleItem().AddKeys(" "),
		keys.Up().AddKeys("k"),
		keys.Down().AddKeys("j"),
		keys.HalfPgUp().AddKeys("K"),
		keys.HalfPgDown().AddKeys("J"),
		keys.Home().AddKeys("g"),
		keys.End().AddKeys("G"),
		keys.Quit().AddKeys("q"),
	}
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
