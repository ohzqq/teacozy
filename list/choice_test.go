package list

import (
	"fmt"
	"testing"
)

func TestNewChoices(t *testing.T) {
	c, err := NewChoices(choiceSlice)
	if err != nil {
		fmt.Errorf("error %v\n", err)
	}
	if len(c) != len(choiceSlice) {
		t.Errorf("len original %v != len new %v", len(choiceSlice), len(c))
	}
	if len(c) > 0 {
		if c[0] == nil {
			t.Errorf("choice %v is nil", choiceSlice[0])
		}
	}
	fmt.Printf("%#v\n", c)

	m, err := NewChoices(choiceMap)
	if err != nil {
		fmt.Errorf("error %v\n", err)
	}
	if len(m) != len(choiceMap) {
		t.Errorf("len original %v != len new %v", len(choiceMap), len(m))
	}
	if len(m) > 0 {
		if m[0] == nil {
			t.Errorf("choice %v is nil", choiceMap[0])
		}
	}
	fmt.Printf("%#v\n", m)
}

var choiceSlice = []any{
	"Artichoke",
	"Baking Flour",
	"Bananas",
	"Barley",
	"Bean Sprouts",
	"Bitter Melon",
	1,
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

var choiceMap = []map[any]any{
	map[any]any{"Artichoke": "Baking "},
	map[any]any{"Bananas": "Flour"},
	map[any]any{"Sprouts": "Barley"},
	map[any]any{"Bean": "four"},
	map[any]any{"Bitter": "Melon"},
	map[any]any{"Cod": "Orange"},
	map[any]any{"Sugar": "Apple"},
	map[any]any{"Cashews": "Cucumber"},
	map[any]any{"Curry": "Currywurst"},
	map[any]any{"Dill": "Dragonfruit"},
	map[any]any{"Eggs": "Furikake"},
	map[any]any{"Garlic": "Gherkinhree"},
	map[any]any{"Ginger": "Grapefruit"},
	map[any]any{"Hazelnuts": "Horseradish"},
	map[any]any{"Jicama": "Kohlrabi"},
	map[any]any{"Leeks": "four"},
	map[any]any{"Milk": "Molasses"},
	map[any]any{"Muesli": "six"},
	map[any]any{"Nopal": "Nectarine"},
	map[any]any{"Nutella": "Milk"},
	map[any]any{"Oatmeal": "Olives"},
	map[any]any{"Papaya": "Gherkin"},
	map[any]any{"Peppers": "Pickle"},
	map[any]any{"Pineapple": "Plantains"},
	map[any]any{"Pocky": "Quince"},
	map[any]any{"Radish": "Ramps"},
	map[any]any{"Tamarind": "Watermelon"},
	map[any]any{"Weißwurst": "Yams"},
	map[any]any{"Yeast": "Yuzu"},
}
