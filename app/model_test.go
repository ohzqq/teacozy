package app

import (
	"log"
	"testing"

	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/list"
)

func TestModel(t *testing.T) {
	items := list.NewItems(list.ItemsStringSlice(choiceSlice))

	opts := []list.Option{
		list.WithFiltering(true),
		//OrderedList(),
		//list.Editable(true),
		list.WithLimit(10),
		//WithDescription(true),
	}

	app := New(items, opts...)
	p := reactea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	println(app.items.Len())
}

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
