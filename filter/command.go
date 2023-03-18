package filter

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"github.com/sahilm/fuzzy"
)

// Options is the customization options for the filter command.
type Options struct {
	CursorPrefix          string
	CursorStyle           lipgloss.Style
	Limit                 int
	NoLimit               bool
	Strict                bool
	SelectedPrefix        string
	SelectedPrefixStyle   lipgloss.Style
	UnselectedPrefix      string
	UnselectedPrefixStyle lipgloss.Style
	HeaderStyle           lipgloss.Style
	Header                string
	TextStyle             lipgloss.Style
	MatchStyle            lipgloss.Style
	Placeholder           string
	Prompt                string
	PromptStyle           lipgloss.Style
	Width                 int
	Height                int
	Value                 string
	Reverse               bool
	Fuzzy                 bool
	Choices               []string
}

func New(o Options) *Model {
	o.CursorPrefix = style.Cursor
	o.CursorStyle = style.CursorStyle
	o.Prompt = style.Prompt
	o.PromptStyle = style.PromptStyle
	o.SelectedPrefix = style.SelectedPrefix
	o.SelectedPrefixStyle = style.SelectedStyle
	o.UnselectedPrefix = style.UnselectedPrefix
	o.UnselectedPrefixStyle = style.UnselectedStyle
	o.TextStyle = lipgloss.NewStyle().Foreground(color.Foreground)
	o.MatchStyle = lipgloss.NewStyle().Foreground(color.Pink)
	tm := Model{
		Options:     o,
		selected:    make(map[int]struct{}),
		FilterKeys:  FilterKeyMap,
		ListKeys:    ListKeyMap,
		filterState: Unfiltered,
	}

	tm.textinput = textinput.New()
	//model.textinput.Focus()

	tm.textinput.Prompt = o.Prompt
	tm.textinput.PromptStyle = o.PromptStyle
	tm.textinput.Placeholder = o.Placeholder
	tm.textinput.Width = o.Width

	w, h := util.TermSize()
	if tm.Height == 0 {
		tm.Height = h
	}
	if tm.Width == 0 {
		tm.Width = w
	}

	v := viewport.New(tm.Width, tm.Height)
	tm.viewport = &v

	tm.Items = make([]Item, len(o.Choices))

	for i, thing := range o.Choices {
		//for k, option := range thing {
		tm.Items[i] = Item{
			Index:    i,
			Text:     thing,
			Selected: false,
			Order:    i,
		}
		//}
	}

	if tm.Value != "" {
		tm.textinput.SetValue(tm.Value)
	}
	switch {
	case tm.Value != "" && tm.Fuzzy:
		tm.matches = fuzzy.Find(tm.Value, tm.Choices)
	case tm.Value != "" && !tm.Fuzzy:
		tm.matches = exactMatches(tm.Value, tm.Items)
	default:
		tm.matches = matchAll(tm.Items)
	}

	if tm.NoLimit {
		tm.Limit = len(tm.Choices)
	}

	return &tm
}

func EnterCmd(m *Model) tea.Cmd {
	return ReturnSelectionsCmd(m)
}

func FilterItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Filtering
		m.cursor = 0
		m.textinput.Focus()
		return textinput.Blink()
	}
}

func StopFilteringCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Unfiltered
		m.textinput.Blur()
		return nil
	}
}

type ReturnSelectionsMsg struct {
	choices []string
}

func ReturnSelectionsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.numSelected < 1 {
			m.Items[m.cursor].Selected = true
		}
		var sel ReturnSelectionsMsg
		if len(m.selected) > 0 {
			for k := range m.selected {
				sel.choices = append(sel.choices, m.Choices[k])
			}
		} else if len(m.matches) > m.cursor && m.cursor >= 0 {
			sel.choices = append(sel.choices, m.matches[m.cursor].Str)
		}
		return sel
	}
}

func SelectItemCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Limit == 1 {
			return nil
		}

		if _, ok := m.selected[m.matches[m.cursor].Index]; ok {
			delete(m.selected, m.matches[m.cursor].Index)
			m.numSelected--
		} else if m.numSelected < m.Limit {
			m.currentOrder++
			m.selected[m.matches[m.cursor].Index] = struct{}{}
			m.numSelected++
			m.CursorDown()
		}

		return nil
	}
}

func UpCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.CursorUp()
		return nil
	}
}

func DownCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.CursorDown()
		return nil
	}
}

func TopCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = 0
		//m.paginator.Page = 0
		return nil
	}
}

func BottomCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = len(m.Items) - 1
		//m.paginator.Page = m.paginator.TotalPages - 1
		return nil
	}
}
