package filter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/files"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
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
}

func New(o Options) *Model {
	model := Model{
		choices:               choices,
		Indicator:             o.CursorPrefix,
		matches:               matches,
		Header:                o.Header,
		textinput:             i,
		viewport:              &v,
		IndicatorStyle:        o.CursorStyle.ToLipgloss(),
		UnselectedPrefixStyle: o.SelectedPrefixStyle.ToLipgloss(),
		SelectedPrefix:        o.SelectedPrefix,
		UnselectedPrefixStyle: o.UnselectedPrefixStyle.ToLipgloss(),
		UnselectedPrefix:      o.UnselectedPrefix,
		MatchStyle:            o.MatchStyle.ToLipgloss(),
		HeaderStyle:           o.HeaderStyle.ToLipgloss(),
		MatchStyle:            o.TextStyle.ToLipgloss(),
		Height:                o.Height,
		selected:              make(map[string]struct{}),
		Limit:                 o.Limit,
		Reverse:               o.Reverse,
		Fuzzy:                 o.Fuzzy,
	}
	return &model
}

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run() error {
	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.PromptStyle = o.PromptStyle.ToLipgloss()
	i.Placeholder = o.Placeholder
	i.Width = o.Width

	v := viewport.New(o.Width, o.Height)

	var choices []string
	if input, _ := stdin.Read(); input != "" {
		input = strings.TrimSuffix(input, "\n")
		if input != "" {
			choices = strings.Split(input, "\n")
		}
	} else {
		choices = files.List()
	}

	if len(choices) == 0 {
		return errors.New("no options provided, see `gum filter --help`")
	}

	//options := []tea.ProgramOption{tea.WithOutput(os.Stderr)}
	options := []tea.ProgramOption{}
	if o.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	var matches []fuzzy.Match
	if o.Value != "" {
		i.SetValue(o.Value)
	}
	switch {
	case o.Value != "" && o.Fuzzy:
		matches = fuzzy.Find(o.Value, choices)
	case o.Value != "" && !o.Fuzzy:
		matches = exactMatches(o.Value, choices)
	default:
		matches = matchAll(choices)
	}

	if o.NoLimit {
		o.Limit = len(choices)
	}

	p := tea.NewProgram(Model{
		choices:               choices,
		Indicator:             o.CursorPrefix,
		matches:               matches,
		Header:                o.Header,
		textinput:             i,
		viewport:              &v,
		IndicatorStyle:        o.CursorStyle.ToLipgloss(),
		UnselectedPrefixStyle: o.SelectedPrefixStyle.ToLipgloss(),
		SelectedPrefix:        o.SelectedPrefix,
		UnselectedPrefixStyle: o.UnselectedPrefixStyle.ToLipgloss(),
		UnselectedPrefix:      o.UnselectedPrefix,
		MatchStyle:            o.MatchStyle.ToLipgloss(),
		HeaderStyle:           o.HeaderStyle.ToLipgloss(),
		MatchStyle:            o.TextStyle.ToLipgloss(),
		Height:                o.Height,
		selected:              make(map[string]struct{}),
		Limit:                 o.Limit,
		Reverse:               o.Reverse,
		Fuzzy:                 o.Fuzzy,
	}, options...)

	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("unable to run filter: %w", err)
	}
	m := tm.(model)
	if m.aborted {
		return exit.ErrAborted
	}

	// allSelections contains values only if limit is greater
	// than 1 or if flag --no-limit is passed, hence there is
	// no need to further checks
	if len(m.selected) > 0 {
		for k := range m.selected {
			fmt.Println(k)
		}
	} else if len(m.matches) > m.cursor && m.cursor >= 0 {
		fmt.Println(m.matches[m.cursor].Str)
	}

	if !o.Strict && len(m.textinput.Value()) != 0 && len(m.matches) == 0 {
		fmt.Println(m.textinput.Value())
	}
	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}
