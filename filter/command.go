package filter

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
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
	model := Model{
		Options:  o,
		selected: make(map[string]struct{}),
	}
	return &model
}

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run() error {
	m := New(o)
	m.textinput = textinput.New()
	m.textinput.Focus()

	m.textinput.Prompt = o.Prompt
	m.textinput.PromptStyle = style.PromptStyle
	m.textinput.Placeholder = o.Placeholder
	m.textinput.Width = o.Width

	v := viewport.New(o.Width, o.Height)
	m.viewport = &v

	if len(m.Choices) == 0 {
		return errors.New("no options provided, see `gum filter --help`")
	}

	//options := []tea.ProgramOption{tea.WithOutput(os.Stderr)}
	options := []tea.ProgramOption{}
	if m.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	if m.Value != "" {
		m.textinput.SetValue(m.Value)
	}
	switch {
	case m.Value != "" && m.Fuzzy:
		m.matches = fuzzy.Find(m.Value, m.Choices)
	case m.Value != "" && !m.Fuzzy:
		m.matches = exactMatches(m.Value, m.Choices)
	default:
		m.matches = matchAll(m.Choices)
	}

	if m.NoLimit {
		m.Limit = len(m.Choices)
	}

	p := tea.NewProgram(m, options...)

	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("unable to run filter: %w", err)
	}
	model := tm.(Model)
	if model.aborted {
		return fmt.Errorf("aborted")
	}

	// allSelections contains values only if limit is greater
	// than 1 or if flag --no-limit is passed, hence there is
	// no need to further checks
	if len(model.selected) > 0 {
		for k := range model.selected {
			fmt.Println(k)
		}
	} else if len(model.matches) > model.cursor && model.cursor >= 0 {
		fmt.Println(model.matches[model.cursor].Str)
	}

	if !o.Strict && len(model.textinput.Value()) != 0 && len(model.matches) == 0 {
		fmt.Println(model.textinput.Value())
	}
	return nil
}
