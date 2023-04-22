package keys

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type KeyMap []*Binding

type Binding struct {
	key.Binding
	help   string
	TeaCmd tea.Cmd
}

func New(keys ...string) *Binding {
	k := Binding{
		Binding: key.NewBinding(),
	}
	k.WithKeys(keys...)
	return &k
}

func NewBind(k key.Binding) *Binding {
	return &Binding{
		Binding: k,
	}
}

func (k *Binding) Cmd(cmd tea.Cmd) *Binding {
	k.TeaCmd = cmd
	return k
}

func (k *Binding) WithKeys(keys ...string) *Binding {
	k.Binding.SetKeys(keys...)
	k.WithHelp(k.help)
	return k
}

func (k *Binding) AddKeys(keys ...string) *Binding {
	keys = append(keys, k.Binding.Keys()...)
	k.Binding.SetKeys(keys...)
	k.WithHelp(k.help)
	return k
}

func (k *Binding) WithHelp(h string) *Binding {
	k.help = h
	k.Binding.SetHelp(strings.Join(k.Keys(), "/"), h)
	return k
}

func (km KeyMap) Map() []map[string]string {
	c := make([]map[string]string, len(km))
	for i, k := range km {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}

func (km KeyMap) Get(name string) *Binding {
	for _, bind := range km {
		for _, k := range bind.Keys() {
			if k == name {
				return bind
			}
		}
	}
	return km.New(name)
}

func (km KeyMap) New(keys ...string) *Binding {
	b := New(keys...)
	km.AddBind(b)
	return b
}

func (km KeyMap) AddBind(b *Binding) {
	km = append(km, b)
}

func MapKeys(keys ...key.Binding) []map[string]string {
	c := make([]map[string]string, len(keys))
	for i, k := range keys {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}

var Global = KeyMap{
	Quit(),
}

func Enter() *Binding {
	return New("enter")
}

func Filter() *Binding {
	return New("/").
		WithHelp("filter items").
		Cmd(ChangeRoute("filter"))
}

func Save() *Binding {
	return New("ctrl+s").
		WithHelp("save edit")
}

func HalfPgUp() *Binding {
	return New("ctrl+u").
		WithHelp("½ page up").
		Cmd(HalfPageUp)
}

func HalfPgDown() *Binding {
	return New("ctrl+d").
		WithHelp("½ page down").
		Cmd(HalfPageDown)
}

func PgUp() *Binding {
	return New("pgup").
		WithHelp("page up").
		Cmd(PageUp)
}

func PgDown() *Binding {
	return New("pgdown").
		WithHelp("page down").
		Cmd(PageDown)
}

func End() *Binding {
	return New("end").
		WithHelp("list bottom").
		Cmd(Bottom)
}

func Home() *Binding {
	return New("home").
		WithHelp("list top").
		Cmd(Top)
}

func Up() *Binding {
	return New("up").
		WithHelp("move up").
		Cmd(LineUp)
}

func Down() *Binding {
	return New("down").
		WithHelp("move down").
		Cmd(LineDown)
}

func Next() *Binding {
	return New("right").
		WithHelp("next page").
		Cmd(NextPage)
}

func Prev() *Binding {
	return New("left").
		WithHelp("prev page").
		Cmd(PrevPage)
}

func Toggle() *Binding {
	return New("tab").
		WithHelp("select item").
		Cmd(ToggleItem)
}

func Quit() *Binding {
	return New("ctrl+c").
		WithHelp("quit program").
		Cmd(reactea.Destroy)
}

func Help() *Binding {
	return New("f1").
		WithHelp("show help").
		Cmd(ShowHelp)
}

func Yes() *Binding {
	return New("y").
		WithHelp("confirm action")
}

func No() *Binding {
	return New("n").
		WithHelp("reject action")
}

func Esc() *Binding {
	return New("esc").
		WithHelp("exit screen").
		Cmd(ReturnToList)
}

func Edit() *Binding {
	return New("e").
		WithHelp("edit field").
		Cmd(ChangeRoute("edit"))
}

func TextInput() KeyMap {
	return textInput
}

func TextArea() KeyMap {
	km := textInput
	km = append(km, textArea...)
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

var textInput = KeyMap{
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
}

var textArea = KeyMap{
	NewBind(InsertNewline),
	NewBind(LineNext),
	NewBind(LinePrevious),
	NewBind(InputBegin),
	NewBind(InputEnd),
	NewBind(UppercaseWordForward),
	NewBind(LowercaseWordForward),
	NewBind(CapitalizeWordForward),
	NewBind(TransposeCharacterBackward),
}
