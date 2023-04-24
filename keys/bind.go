package keys

import "github.com/londek/reactea"

func Enter() *Binding {
	return New("enter")
}

func Filter() *Binding {
	return New("/").
		WithHelp("filter items")
	//Cmd(ChangeRoute("filter"))
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
		Cmd(StartEditing)
}

func Save() *Binding {
	return New("ctrl+s").
		WithHelp("save edit")
}
