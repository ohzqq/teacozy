package filter

import (
	"errors"
	"fmt"

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

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run() error {
	m := New(o)
	if len(m.Choices) == 0 {
		return errors.New("no options provided")
	}

	//options := []tea.ProgramOption{tea.WithOutput(os.Stderr)}
	options := []tea.ProgramOption{}
	if m.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	p := tea.NewProgram(m, options...)

	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("unable to run filter: %w", err)
	}
	model := tm.(*Model)
	if model.aborted {
		return fmt.Errorf("aborted")
	}

	// allSelections contains values only if limit is greater
	// than 1
	if len(model.selected) > 0 {
		for k := range model.selected {
			fmt.Println(model.Choices[k])
		}
	} else if len(model.matches) > model.cursor && model.cursor >= 0 {
		fmt.Println(model.matches[model.cursor].Str)
	}

	if !o.Strict && len(model.textinput.Value()) != 0 && len(model.matches) == 0 {
		fmt.Println(model.textinput.Value())
	}
	return nil
}
