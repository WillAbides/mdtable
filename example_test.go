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

	b := mdtable.Generate(data)
	fmt.Println(string(b))

	// Output:
	// | Name      | Favorite Animal | Lucky Number |
	// |-----------|-----------------|--------------|
	// | Dave      | Elephant        | 7            |
	// | Iris      | Gorilla         | 8            |
	// | Ava Gayle | Sloth           | 972.5        |
}

func Example_options() {
	// This adds options one at a time and shows what the output of
	// mdtable.Generate would be after each option is added.

	data := [][]string{
		{"Name", "Favorite Animal", "Lucky Number"},
		{"Dave", "Elephant", "7"},
		{"Iris", "Gorilla", "8"},
		{"Ava Gayle", "Sloth", "972.5"},
	}

	var options []mdtable.Option

	// Right align the whole table
	options = append(options,
		mdtable.Alignment(mdtable.AlignRight),
	)

	/*
		|      Name | Favorite Animal | Lucky Number |
		|----------:|----------------:|-------------:|
		|      Dave |        Elephant |            7 |
		|      Iris |         Gorilla |            8 |
		| Ava Gayle |           Sloth |        972.5 |
	*/

	// Left align header text
	options = append(options,
		mdtable.HeaderAlignment(mdtable.AlignLeft),
	)

	/*
		| Name      | Favorite Animal | Lucky Number |
		|----------:|----------------:|-------------:|
		|      Dave |        Elephant |            7 |
		|      Iris |         Gorilla |            8 |
		| Ava Gayle |           Sloth |        972.5 |
	*/

	// Set Favorite Animal's (offset 1) minimum width to 20 and center its text
	options = append(options,
		mdtable.ColumnMinWidth(1, 20),
		mdtable.ColumnTextAlignment(1, mdtable.AlignCenter),
	)

	/*
		| Name      | Favorite Animal      | Lucky Number |
		|----------:|---------------------:|-------------:|
		|      Dave |       Elephant       |            7 |
		|      Iris |       Gorilla        |            8 |
		| Ava Gayle |        Sloth         |        972.5 |
	*/

	options = append(options,
		mdtable.ColumnAlignment(0, mdtable.AlignLeft), // Left align Name
	)

	b := mdtable.Generate(data, options...)
	fmt.Println(string(b))

	// Output:
	// | Name      | Favorite Animal      | Lucky Number |
	// |:----------|---------------------:|-------------:|
	// | Dave      |       Elephant       |            7 |
	// | Iris      |       Gorilla        |            8 |
	// | Ava Gayle |        Sloth         |        972.5 |
}
