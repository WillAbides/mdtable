package mdtable

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:generate go test . -write-golden

func TestAlign_fillCell(t *testing.T) {
	runTest := func(align Align, width, padding int, s, want string) {
		t.Helper()
		got := align.fillCell(s, width, padding)
		require.Equal(t, want, got)
	}

	runTest(AlignLeft, 8, 1, "世世世", " 世世世   ")

	runTest(AlignLeft, 8, 1, "foo", " foo      ")
	runTest(AlignLeft, 8, 0, "foo", "foo     ")
	runTest(AlignLeft, 2, 1, "foo", " foo ")

	runTest(AlignDefault, 8, 1, "foo", " foo      ")
	runTest(AlignDefault, 8, 0, "foo", "foo     ")
	runTest(AlignDefault, 2, 1, "foo", " foo ")

	runTest(AlignRight, 8, 1, "foo", "      foo ")
	runTest(AlignRight, 8, 0, "foo", "     foo")
	runTest(AlignRight, 2, 1, "foo", " foo ")

	runTest(AlignCenter, 8, 1, "foo", "   foo    ")
	runTest(AlignCenter, 8, 0, "foo", "  foo   ")
	runTest(AlignCenter, 2, 1, "foo", " foo ")
	runTest(AlignCenter, 9, 1, "foo", "    foo    ")
	runTest(AlignCenter, 9, 0, "foo", "   foo   ")
}

type testTable struct {
	name  string
	table *Table
}

func buildTable(fn func(tbl *Table)) *Table {
	table := new(Table)
	table.SetData(exampleData)
	if fn != nil {
		fn(table)
	}
	return table
}

var exampleData = [][]string{
	{"Date", "Description", "CV2", "Amount"},
	{"1/1/2014", "Domain name", "2233", "$10.98"},
	{"1/1/2014", "January Hosting", "2233", "$54.95"},
	{"1/4/2014", "February Hosting", "2233", "$51.00"},
	{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
}

var testTables = []testTable{
	{
		name:  "defaults",
		table: buildTable(nil),
	},
	{
		name: "alignments-1",
		table: buildTable(func(tbl *Table) {
			tbl.SetAlignment(AlignCenter)
			tbl.SetColumnAlignment(1, AlignLeft)
			tbl.SetColumnAlignment(2, AlignRight)
			tbl.SetColumnTextAlignment(1, AlignCenter)
			tbl.SetColumnHeaderAlignment(1, AlignLeft)
			tbl.SetColumnMinWidth(0, 12)
			tbl.SetColumnMinWidth(2, 12)
			tbl.SetColumnMinWidth(3, 12)
		}),
	},
}

func updateGolden() error {
	err := os.RemoveAll(filepath.Join("testdata", "tables"))
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join("testdata", "tables"), 0o700)
	if err != nil {
		return err
	}
	for i := range testTables {
		name := testTables[i].name
		table := testTables[i].table
		got := table.Render()
		got = append(got, '\n')
		err = ioutil.WriteFile(filepath.Join("testdata", "tables", name+".md"), got, 0o600)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	var err error
	var writeGolden bool
	flag.BoolVar(&writeGolden, "write-golden", false, "write golden files")
	flag.Parse()
	if writeGolden {
		err = updateGolden()
		if err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func TestRender(t *testing.T) {
	for _, td := range testTables {
		t.Run(td.name, func(t *testing.T) {
			want, err := ioutil.ReadFile(filepath.Join("testdata", "tables", td.name+".md"))
			require.NoError(t, err)
			want = bytes.TrimSuffix(want, []byte{'\n'})
			got := td.table.Render()
			require.Equal(t, string(want), string(got))
		})
	}
}
