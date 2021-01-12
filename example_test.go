package mdtable_test

import (
	"fmt"

	"github.com/willabides/mdtable"
)

func Example() {
	data := [][]string{
		// first row is the header
		{"Name", "Favorite Animal", "Lucky Number"},

		// the rest is data
		{"Dave", "Elephant", "7"},
		{"Iris", "Gorilla", "8"},
		{"Ava Gayle", "Sloth", "972.5"},
	}
	table := mdtable.New(data)

	// Table's String() renders the table as markdown
	fmt.Println(table)

	// Output:
	// | Name      | Favorite Animal | Lucky Number |
	// |-----------|-----------------|--------------|
	// | Dave      | Elephant        | 7            |
	// | Iris      | Gorilla         | 8            |
	// | Ava Gayle | Sloth           | 972.5        |
}

func Example_alignment() {
	table := mdtable.New([][]string{
		{"Name", "Favorite Animal", "Lucky Number"},
		{"Dave", "Elephant", "7"},
		{"Iris", "Gorilla", "8"},
		{"Ava Gayle", "Sloth", "972.5"},
	})

	// Right align the whole table
	table.SetAlignment(mdtable.AlignRight)

	// |      Name | Favorite Animal | Lucky Number |
	// |----------:|----------------:|-------------:|
	// |      Dave |        Elephant |            7 |
	// |      Iris |         Gorilla |            8 |
	// | Ava Gayle |           Sloth |        972.5 |

	// Left align header text
	table.SetHeaderAlignment(mdtable.AlignLeft)

	// | Name      | Favorite Animal | Lucky Number |
	// |----------:|----------------:|-------------:|
	// |      Dave |        Elephant |            7 |
	// |      Iris |         Gorilla |            8 |
	// | Ava Gayle |           Sloth |        972.5 |

	// Set Favorite Animal's (offset 1) minimum width to 20
	table.SetColumnMinWidth(1, 20)

	// Center align Favorite Animal's text
	table.SetColumnTextAlignment(1, mdtable.AlignCenter)

	// | Name      | Favorite Animal      | Lucky Number |
	// |----------:|---------------------:|-------------:|
	// |      Dave |       Elephant       |            7 |
	// |      Iris |       Gorilla        |            8 |
	// | Ava Gayle |        Sloth         |        972.5 |

	// Left align Name
	table.SetColumnAlignment(0, mdtable.AlignLeft)

	fmt.Println(table)
	// Output:
	// | Name      | Favorite Animal      | Lucky Number |
	// |:----------|---------------------:|-------------:|
	// | Dave      |       Elephant       |            7 |
	// | Iris      |       Gorilla        |            8 |
	// | Ava Gayle |        Sloth         |        972.5 |
}
