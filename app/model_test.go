package app

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
	"golang.org/x/term"
)

//func TestNewList(t *testing.T) {
//  //var opts []Option
//  //opts = append(opts, NoLimit())
//  //cs := New(choiceSlice, opts...).Choose()
//  //fmt.Printf("%#v\n", cs)
//  //items := NewItems(ItemsMap(choiceMap), OrderedList())
//  //p.SetSize(0, 10)

//  //items := list.NewItems(list.ItemsStringSlice(choiceSlice))

//  opts := []Option{
//    WithList(testItemParser(), testListOpts()...),
//    WithDescription(),
//  }

//  //m := list.New(items, opts...)

//  //m := EditableList(items)
//  //m := NewEditableList(noItems)

//  a := New(opts...)
//  //SetList(m).
//  //SetPager(testPager())
//  a.AddCommands(testCommand())

//  err := a.Run()
//  if err != nil {
//    log.Fatal(err)
//  }
//  //sel := m.Chosen()
//  //for _, s := range sel {
//  //fmt.Printf("%#v\n", s)
//  //}
//  //fmt.Printf("%#v\n", a.ShowCommand())

//  println(10 / 1 * (1 - 0))
//}

var (
	oneSub       = "0x0"
	oneMain      = "66x66"
	halfVMain    = "66x33"
	halfVSub     = "66x33"
	halfHMain    = "33x66"
	halfHSub     = "33x66"
	thirdVMain   = "66x44"
	thirdVSub    = "66x22"
	thirdHMain   = "44x66"
	thirdHSub    = "22x66"
	quarterVMain = "66x48"
	quarterVSub  = "66x16"
	quarterHMain = "48x66"
	quarterHSub  = "16x66"
)

func TestOneHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(1).Split(Horizontal)
	main := printSize(l.Main())
	sub := printSize(l.Sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	switch l.split {
	case 1:
		if main != oneMain {
			t.Errorf("got %s expect %s\n", main, oneMain)
		}
		if sub != oneSub {
			t.Errorf("got %s expect %s\n", sub, oneSub)
		}
	case Half:
		if main != halfHMain {
			t.Errorf("got %s expect %s\n", main, halfHMain)
		}
		if sub != halfHSub {
			t.Errorf("got %s expect %s\n", sub, halfHSub)
		}
	case Third:
		if main != thirdHMain {
			t.Errorf("got %s expect %s\n", main, thirdHMain)
		}
		if sub != thirdHSub {
			t.Errorf("got %s expect %s\n", sub, thirdHSub)
		}
	case Quarter:
		if main != quarterHMain {
			t.Errorf("got %s expect %s\n", main, quarterHMain)
		}
		if sub != quarterHSub {
			t.Errorf("got %s expect %s\n", sub, quarterHSub)
		}
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}

func TestHalfHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(1).Split(Horizontal)
	main := printSize(l.Main())
	sub := printSize(l.Sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	switch l.split {
	case 1:
		if main != oneMain {
			t.Errorf("got %s expect %s\n", main, oneMain)
		}
		if sub != oneSub {
			t.Errorf("got %s expect %s\n", sub, oneSub)
		}
	case Half:
		if main != halfHMain {
			t.Errorf("got %s expect %s\n", main, halfHMain)
		}
		if sub != halfHSub {
			t.Errorf("got %s expect %s\n", sub, halfHSub)
		}
	case Third:
		if main != thirdHMain {
			t.Errorf("got %s expect %s\n", main, thirdHMain)
		}
		if sub != thirdHSub {
			t.Errorf("got %s expect %s\n", sub, thirdHSub)
		}
	case Quarter:
		if main != quarterHMain {
			t.Errorf("got %s expect %s\n", main, quarterHMain)
		}
		if sub != quarterHSub {
			t.Errorf("got %s expect %s\n", sub, quarterHSub)
		}
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}

func TestThirdHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(1).Split(Horizontal)
	main := printSize(l.Main())
	sub := printSize(l.Sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	switch l.split {
	case 1:
		if main != oneMain {
			t.Errorf("got %s expect %s\n", main, oneMain)
		}
		if sub != oneSub {
			t.Errorf("got %s expect %s\n", sub, oneSub)
		}
	case Half:
		if main != halfHMain {
			t.Errorf("got %s expect %s\n", main, halfHMain)
		}
		if sub != halfHSub {
			t.Errorf("got %s expect %s\n", sub, halfHSub)
		}
	case Third:
		if main != thirdHMain {
			t.Errorf("got %s expect %s\n", main, thirdHMain)
		}
		if sub != thirdHSub {
			t.Errorf("got %s expect %s\n", sub, thirdHSub)
		}
	case Quarter:
		if main != quarterHMain {
			t.Errorf("got %s expect %s\n", main, quarterHMain)
		}
		if sub != quarterHSub {
			t.Errorf("got %s expect %s\n", sub, quarterHSub)
		}
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}

func TestQuarterHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(1).Split(Horizontal)
	main := printSize(l.Main())
	sub := printSize(l.Sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	switch l.split {
	case 1:
		if main != oneMain {
			t.Errorf("got %s expect %s\n", main, oneMain)
		}
		if sub != oneSub {
			t.Errorf("got %s expect %s\n", sub, oneSub)
		}
	case Half:
		if main != halfHMain {
			t.Errorf("got %s expect %s\n", main, halfHMain)
		}
		if sub != halfHSub {
			t.Errorf("got %s expect %s\n", sub, halfHSub)
		}
	case Third:
		if main != thirdHMain {
			t.Errorf("got %s expect %s\n", main, thirdHMain)
		}
		if sub != thirdHSub {
			t.Errorf("got %s expect %s\n", sub, thirdHSub)
		}
	case Quarter:
		if main != quarterHMain {
			t.Errorf("got %s expect %s\n", main, quarterHMain)
		}
		if sub != quarterHSub {
			t.Errorf("got %s expect %s\n", sub, quarterHSub)
		}
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}
func TestOneVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(1)
	main := printSize(l.Main())
	sub := printSize(l.Sub())

	if main != oneMain {
		t.Errorf("got %s expect %s\n", main, oneMain)
	}
	if sub != oneSub {
		t.Errorf("got %s expect %s\n", sub, oneSub)
	}

	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func TestHalfVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Half)
	main := printSize(l.Main())
	sub := printSize(l.Sub())

	if main != halfVMain {
		t.Errorf("got %s expect %s\n", main, halfVMain)
	}
	if sub != halfVSub {
		t.Errorf("got %s expect %s\n", sub, halfVSub)
	}

	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func TestThirdVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Third)
	main := printSize(l.Main())
	sub := printSize(l.Sub())
	if main != thirdVMain {
		t.Errorf("got %s expect %s\n", main, thirdVMain)
	}
	if sub != thirdVSub {
		t.Errorf("got %s expect %s\n", sub, thirdVSub)
	}
	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func TestQuarterVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Quarter)
	main := printSize(l.Main())
	sub := printSize(l.Sub())
	if main != quarterVMain {
		t.Errorf("got %s expect %s\n", main, quarterVMain)
	}
	if sub != quarterVSub {
		t.Errorf("got %s expect %s\n", sub, quarterVSub)
	}
	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func printSize(w, h int) string {
	size := fmt.Sprintf("%dx%d", w, h)
	return size
}

func testCommand() Command {
	return Command{
		Name: "poot",
		//Key:  inKey,
		Cmd: NewStatusMessage,
	}
}

var inKey = key.NewBinding(
	key.WithKeys("a"),
)

func testItemParser() list.ParseItems {
	return list.ItemsStringSlice(choiceSlice)
}

func testListOpts() []list.Option {
	opts := []list.Option{
		list.WithFiltering(true),
		//OrderedList(),
		list.Editable(true),
		list.WithLimit(10),
		//WithDescription(true),
	}
	return opts
}

func testPager() *pager.Model {
	txt := []string{
		"AArtichokeArtichokeArtichokeArtichokeArtichokeArtichokeArtichokertichoke",
	}
	txt = append(txt, choiceSlice...)
	text := strings.Join(txt, "\n- ")
	p := pager.New(pager.RenderText).SetText(text)
	return p
}

func testTermSize(t *testing.T) {
	oldState, err := term.MakeRaw(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdout.Fd()), oldState)
	//if term.IsTerminal(oldState) {
	//  println("in a term")
	//} else {
	//  println("not in a term")
	//}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		//return
		log.Fatal(err)
	}
	println("width:", width, "height:", height)

}

//func TestNewBookList(t *testing.T) {
//  d, err := os.ReadFile("../testdata/search-results.json")
//  if err != nil {
//    fmt.Printf("%v\n", err)
//  }

//  var books []cdb.Book
//  err = json.Unmarshal(d, &books)
//  if err != nil {
//    fmt.Printf("%v\n", err)
//  }

//  var choices []string
//  for _, book := range books {
//    choices = append(choices, book.Title)
//  }

//  var opts []Option
//  opts = append(opts, NoLimit())
//  cs := New(choices, opts...).Choose()

//  fmt.Printf("%#v\n", cs)

//}

var choiceSlice = []string{
	"Artichoke",
	"Baking Flour",
	"Bananas",
	"Barley",
	"Bean Sprouts",
	"Bitter Melon",
	"Blood Orange",
	"Brown Sugar",
	"Cashew Apple",
	"Cashews",
	"Cat Food",
	"Coconut Milk",
	"Cucumber",
	"Curry Paste",
	"Currywurst",
	"Dill",
	"Dragonfruit",
	"Dried Shrimp",
	"Eggs",
	"Fish Cake",
	"Furikake",
	"Garlic",
	"Gherkin",
	"Ginger",
	"Granulated Sugar",
	"Grapefruit",
	"Green Onion",
	"Hazelnuts",
	"Heavy whipping cream",
	"Honey Dew",
	"Horseradish",
	"Jicama",
	"Kohlrabi",
	"Leeks",
	"Lentils",
	"Licorice Root",
	"Meyer Lemons",
	"Milk",
	"Molasses",
	"Muesli",
	"Nectarine",
	"Niagamo Root",
	"Nopal",
	"Nutella",
	"Oat Milk",
	"Oatmeal",
	"Olives",
	"Papaya",
	"Party Gherkin",
	"Peppers",
	"Persian Lemons",
	"Pickle",
	"Pineapple",
	"Plantains",
	"Pocky",
	"Powdered Sugar",
	"Quince",
	"Radish",
	"Ramps",
	"Star Anise",
	"Sweet Potato",
	"Tamarind",
	"Unsalted Butter",
	"Watermelon",
	"Weißwurst",
	"Yams",
	"Yeast",
	"Yuzu",
}

var choiceMap = map[string]string{
	"Artichoke": "Baking ",
	"Bananas":   "Flour",
	"Sprouts":   "Barley",
	"Bean":      "four",
	"Bitter":    "Melon",
	"Cod":       "Orange",
	"Sugar":     "Apple",
	"Cashews":   "Cucumber",
	"Curry":     "Currywurst",
	"Dill":      "Dragonfruit",
	"Eggs":      "Furikake",
	"Garlic":    "Gherkinhree",
	"Ginger":    "Grapefruit",
	"Hazelnuts": "Horseradish",
	"Jicama":    "Kohlrabi",
	"Leeks":     "four",
	"Milk":      "Molasses",
	"Muesli":    "six",
	"Nopal":     "Nectarine",
	"Nutella":   "Milk",
	"Oatmeal":   "Olives",
	"Papaya":    "Gherkin",
	"Peppers":   "Pickle",
	"Pineapple": "Plantains",
	"Pocky":     "Quince",
	"Radish":    "Ramps",
	"Tamarind":  "Watermelon",
	"Weißwurst": "Yams",
	"Yeast":     "Yuzu",
}

var choiceSliceMap = []map[string]string{
	map[string]string{"Artichoke": "Baking "},
	map[string]string{"Bananas": "Flour"},
	map[string]string{"Sprouts": "Barley"},
	map[string]string{"Bean": "four"},
	map[string]string{"Bitter": "Melon"},
	map[string]string{"Cod": "Orange"},
	map[string]string{"Sugar": "Apple"},
	map[string]string{"Cashews": "Cucumber"},
	map[string]string{"Curry": "Currywurst"},
	map[string]string{"Dill": "Dragonfruit"},
	map[string]string{"Eggs": "Furikake"},
	map[string]string{"Garlic": "Gherkinhree"},
	map[string]string{"Ginger": "Grapefruit"},
	map[string]string{"Hazelnuts": "Horseradish"},
	map[string]string{"Jicama": "Kohlrabi"},
	map[string]string{"Leeks": "four"},
	map[string]string{"Milk": "Molasses"},
	map[string]string{"Muesli": "six"},
	map[string]string{"Nopal": "Nectarine"},
	map[string]string{"Nutella": "Milk"},
	map[string]string{"Oatmeal": "Olives"},
	map[string]string{"Papaya": "Gherkin"},
	map[string]string{"Peppers": "Pickle"},
	map[string]string{"Pineapple": "Plantains"},
	map[string]string{"Pocky": "Quince"},
	map[string]string{"Radish": "Ramps"},
	map[string]string{"Tamarind": "Watermelon"},
	map[string]string{"Weißwurst": "Yams"},
	map[string]string{"Yeast": "Yuzu"},
}
