package keys

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

func TextInput() KeyMap {
	return textInput
}

func TextArea() KeyMap {
	km := textInput
	km.AddBinds(textArea.Keys()...)
	//km = append(km, textArea...)
	return km
}

func TextInputDefault() textinput.KeyMap {
	km := textinput.DefaultKeyMap
	km.DeleteAfterCursor.Unbind()
	km.DeleteBeforeCursor.Unbind()
	km.CharacterForward = CharacterForward
	km.CharacterBackward = CharacterBackward
	km.WordForward = WordForward
	km.WordBackward = WordBackward
	km.DeleteWordBackward = DeleteWordBackward
	km.DeleteWordForward = DeleteWordForward
	km.DeleteCharacterBackward = DeleteCharacterBackward
	km.DeleteCharacterForward = DeleteCharacterForward
	km.LineStart = LineStart
	km.LineEnd = LineEnd
	km.Paste = Paste
	return km
}

func TextAreaDefault() textarea.KeyMap {
	km := textarea.DefaultKeyMap
	km.DeleteAfterCursor.Unbind()
	km.DeleteBeforeCursor.Unbind()
	km.CharacterForward = CharacterForward
	km.CharacterBackward = CharacterBackward
	km.WordForward = WordForward
	km.WordBackward = WordBackward
	km.DeleteWordBackward = DeleteWordBackward
	km.DeleteWordForward = DeleteWordForward
	km.DeleteCharacterBackward = DeleteCharacterBackward
	km.DeleteCharacterForward = DeleteCharacterForward
	km.LineStart = LineStart
	km.LineEnd = LineEnd
	km.Paste = Paste
	km.InsertNewline = InsertNewline
	km.InputBegin = InputBegin
	km.InputEnd = InputEnd
	km.CapitalizeWordForward = CapitalizeWordForward
	km.LowercaseWordForward = LowercaseWordForward
	km.UppercaseWordForward = UppercaseWordForward
	km.TransposeCharacterBackward = TransposeCharacterBackward
	return km
}

var textInput = KeyMap{
	keys: []*Binding{
		NewBind(CharacterBackward),
		NewBind(CharacterForward),
		NewBind(DeleteCharacterBackward),
		NewBind(DeleteCharacterForward),
		NewBind(DeleteWordBackward),
		NewBind(DeleteWordForward),
		NewBind(LineEnd),
		NewBind(LineStart),
		NewBind(Paste),
		NewBind(WordBackward),
		NewBind(WordForward),
	},
}

var textArea = KeyMap{
	keys: []*Binding{
		NewBind(InsertNewline),
		NewBind(LineNext),
		NewBind(LinePrevious),
		NewBind(InputBegin),
		NewBind(InputEnd),
		NewBind(UppercaseWordForward),
		NewBind(LowercaseWordForward),
		NewBind(CapitalizeWordForward),
		NewBind(TransposeCharacterBackward),
	},
}

var (
	CharacterForward = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("right", "char forward"),
	)
	CharacterBackward = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("left", "character backward"),
	)
	WordForward = key.NewBinding(
		key.WithKeys("alt+right", "alt+l"),
		key.WithHelp("alt+right/alt+l", "word forward"),
	)
	WordBackward = key.NewBinding(
		key.WithKeys("alt+left", "alt+h"),
		key.WithHelp("alt+left/alt+h", "word backward"),
	)
	LineNext = key.NewBinding(
		key.WithKeys("down", "alt+j"),
		key.WithHelp("down/alt+j", "line next"),
	)
	LinePrevious = key.NewBinding(
		key.WithKeys("up", "alt+k"),
		key.WithHelp("up/alt+k", "line prev"),
	)
	DeleteWordBackward = key.NewBinding(
		key.WithKeys("alt+backspace"),
		key.WithHelp("alt+backspace", "delete word backward"),
	)
	DeleteWordForward = key.NewBinding(
		key.WithKeys("alt+delete", "alt+d"),
		key.WithHelp("alt+delete/alt+d", "delete word forward"),
	)
	InsertNewline = key.NewBinding(
		key.WithKeys("enter", "alt+m"),
		key.WithHelp("enter/alt+m", "insert new line"),
	)
	DeleteCharacterBackward = key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "delete char backward"),
	)
	DeleteCharacterForward = key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("delete", "delete char forward"),
	)
	LineStart = key.NewBinding(
		key.WithKeys("alt+a"),
		key.WithHelp("alt+a", "line start"),
	)
	LineEnd = key.NewBinding(
		key.WithKeys("alt+e"),
		key.WithHelp("alt+e", "line end"),
	)
	Paste = key.NewBinding(
		key.WithKeys("ctrl+v"),
		key.WithHelp("ctrl+v", "paste"),
	)
	InputBegin = key.NewBinding(
		key.WithKeys("alt+<", "ctrl+home"),
		key.WithHelp("alt+</ctrl+home", "input begin"),
	)
	InputEnd = key.NewBinding(
		key.WithKeys("alt+>", "ctrl+end"),
		key.WithHelp("alt+>/ctrl+end", "input end"),
	)
	CapitalizeWordForward = key.NewBinding(
		key.WithKeys("alt+c"),
		key.WithHelp("alt+c", "captilize character forward"),
	)
	LowercaseWordForward = key.NewBinding(
		key.WithKeys("alt+l"),
		key.WithHelp("alt+l", "lowercase word forward"),
	)
	UppercaseWordForward = key.NewBinding(
		key.WithKeys("alt+u"),
		key.WithHelp("alt+u", "uppercase word forward"),
	)
	TransposeCharacterBackward = key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "transpose charactre backward"),
	)
)
