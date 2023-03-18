package choose

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

var (
	subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
	verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
)

// Options is the customization options for the choose command.
type Options struct {
	Options               []map[int]string
	Limit                 int
	NoLimit               bool
	Ordered               bool
	Height                int
	Width                 int
	Cursor                string
	CursorPrefix          string
	SelectedPrefix        string
	UnselectedPrefix      string
	CursorStyle           lipgloss.Style
	ItemStyle             lipgloss.Style
	SelectedItemStyle     lipgloss.Style
	TextStyle             lipgloss.Style
	MatchStyle            lipgloss.Style
	SelectedPrefixStyle   lipgloss.Style
	UnselectedPrefixStyle lipgloss.Style
	HeaderStyle           lipgloss.Style
	Placeholder           string
	Header                string
	Prompt                string
	PromptStyle           lipgloss.Style
	Value                 string
	Reverse               bool
	Fuzzy                 bool
	Strict                bool
}

func New(o Options) *Model {
	filterIn := textinput.New()
	filterIn.Focus()
	tm := Model{
		Options:   o,
		KeyMap:    ListKeyMap,
		textinput: filterIn,
	}
	tm.Cursor = style.Cursor
	tm.SelectedPrefix = style.SelectedPrefix
	tm.UnselectedPrefix = style.UnselectedPrefix
	tm.CursorPrefix = style.CursorPrefix

	tm.CursorStyle = style.CursorStyle
	tm.ItemStyle = style.UnselectedStyle
	tm.SelectedItemStyle = style.SelectedStyle

	w, h := util.TermSize()

	if tm.Height == 0 {
		tm.Height = h
	}
	if tm.Width == 0 {
		tm.Width = w
	}
	vp := viewport.New(o.Width, o.Height)
	tm.viewport = &vp

	tm.Items = make([]Item, len(o.Options))

	for i, thing := range o.Options {
		for k, option := range thing {
			tm.Items[i] = Item{
				Id:       k,
				Text:     option,
				Selected: false,
				Order:    i,
			}
		}
	}

	if len(tm.Items) == 1 {
		tm.Limit = 1
	}

	// We don't need to display prefixes if we are only picking one option.
	// Simply displaying the cursor is enough.
	if tm.Limit == 1 && !o.NoLimit {
		tm.SelectedPrefix = ""
		tm.UnselectedPrefix = ""
		tm.CursorPrefix = ""
	}

	// If we've set no limit then we can simply select as many options as there
	// are so let's set the limit to the number of options.
	if o.NoLimit {
		tm.Limit = len(o.Options)
	}

	// Use the pagination model to display the current and total number of
	// pages.
	pager := paginator.New()
	pager.SetTotalPages((len(tm.Items) + tm.Height - 1) / tm.Height)
	pager.PerPage = tm.Height
	pager.Type = paginator.Dots
	pager.ActiveDot = subduedStyle.Render("•")
	pager.InactiveDot = verySubduedStyle.Render("•")

	// Disable Keybindings since we will control it ourselves.
	pager.UseHLKeys = false
	pager.UseLeftRightKeys = false
	pager.UseJKKeys = false
	pager.UsePgUpPgDownKeys = false

	tm.paginator = pager

	return &tm
}
