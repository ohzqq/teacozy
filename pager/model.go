package pager

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/ohzqq/bubbles/viewport"
	"github.com/ohzqq/teacozy/util"
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
	w, h := util.TermSize()
	m := &Model{
		render: fn,
		text:   text,
		width:  w,
		height: h,
		view:   viewport.New(w, h),
	}

	return m
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

	//println(r)
	return r
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.width == 0 {
			m.width = msg.Width
		}
		if m.height == 0 {
			m.height = msg.Height
		}
		if msg.Width < m.width {
			m.view.Width = msg.Width
		}
		if msg.Height < m.height {
			m.view.Height = msg.Height
		}
		if msg.Width > m.width {
			m.view.Width = m.width
		}
		if msg.Height > m.height {
			m.view.Height = m.height
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		}
	}

	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	m.view.SetContent(m.render(m.text, m.view.Width))

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.view.View()
}

func (m *Model) SetSize(w, h int) *Model {
	m.width = w
	m.height = h
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
