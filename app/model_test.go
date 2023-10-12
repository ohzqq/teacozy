package app

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
	"golang.org/x/term"
)

func TestNewList(t *testing.T) {
	//var opts []Option
	//opts = append(opts, NoLimit())
	//cs := New(choiceSlice, opts...).Choose()
	//fmt.Printf("%#v\n", cs)
	//items := NewItems(ItemsMap(choiceMap), OrderedList())
	//p.SetSize(0, 10)

	items := list.NewItems(list.ItemsStringSlice(choiceSlice))

	opts := []list.Option{
		list.WithFiltering(true),
		//OrderedList(),
		list.Editable(true),
		//list.WithPager(testPager()),
		//WithLimit(10),
		//WithDescription(true),
	}

	m := list.New(items, opts...)
	//m := Edit(items)
	//m := New(items, WithLimit(1))
	//m := ChooseSome(items, 2)
	//m := New(ItemsMap(choiceMap))
	//m.Editable()

	//m := EditableList(items)
	//m := NewEditableList(noItems)

	a := New().SetList(m)
	//SetInput("poot")
	//SetPager(testPager())
	//SetPager(m.Pager)
	p := tea.NewProgram(a)

	//mod, err := p.Run()
	//if err != nil {
	//return m, err
	//}
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	sel := m.Chosen()
	for _, s := range sel {
		fmt.Printf("%#v\n", s)
	}
	fmt.Printf("%#v\n", a.HasInput())

	//w, h := util.TermSize()
	//println(w)
	//println(h)
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

var noItems = func() []*list.Item { return []*list.Item{} }

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
