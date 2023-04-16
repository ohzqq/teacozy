package style

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	CursorPrefix     = "x"
	PromptPrefix     = "> "
	SelectedPrefix   = "x"
	UnselectedPrefix = " "
)

type prefix struct {
	Text  string
	style lipgloss.Style
}

type itemPrefixes struct {
	prompt     *prefix
	selected   *prefix
	unselected *prefix
	cursor     *prefix
}

var (
	prefixes = itemPrefixes{
		prompt:     prompt,
		selected:   selected,
		unselected: unselected,
		cursor:     cursor,
	}
	prompt = &prefix{
		Text:  PromptPrefix,
		style: Prompt,
	}
	selected = &prefix{
		Text:  SelectedPrefix,
		style: Selected,
	}
	unselected = &prefix{
		Text:  UnselectedPrefix,
		style: Unselected,
	}
	cursor = &prefix{
		Text:  CursorPrefix,
		style: Cursor,
	}
)

func (p prefix) Render() string {
	return p.style.Render(p.Text)
}

func (p *prefix) Set(pre string) {
	p.Text = pre
}

func (p itemPrefixes) Cursor() *prefix {
	return p.cursor
}

func (p itemPrefixes) Selected() *prefix {
	return p.selected
}

func (p itemPrefixes) Unselected() *prefix {
	return p.unselected
}

func (p itemPrefixes) Prompt() *prefix {
	return p.prompt
}

func Prefix() itemPrefixes {
	return prefixes
}
