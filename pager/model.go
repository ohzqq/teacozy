package pager

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/ohzqq/bubbles/viewport"
)

type Model struct {
	view     viewport.Model
	render   Render
	text     string
	rendered string
	width    int
	height   int
}

type Render func(text string, width int) string

func New(text string, fn Render) *Model {
	return &Model{
		render: fn,
		text:   text,
	}
}

func RenderHTML(html string, width int) string {
	conv := md.NewConverter("", true, nil)
	mark, err := conv.ConvertString(html)
	if err != nil {
		return html
	}
	return RenderMarkdown(mark, width)
}

func RenderMarkdown(mark string, width int) string {
	gr, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return mark
	}

	r, err := gr.Render(mark)
	if err != nil {
		return mark
	}
	return r
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) Render(w, h int) {
	m.width = w
	m.height = h
	m.rendered = m.render(m.text, w)
}

func (m Model) View() string {
	m.view = viewport.New(m.width, m.height)
	m.view.SetContent(m.rendered)
	return m.view.View()
}

func (m Model) Init() tea.Cmd {
	return nil
}
