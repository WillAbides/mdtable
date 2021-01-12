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

func Test_cellValue(t *testing.T) {
	require.Empty(t, cellValue(nil, 1, 1))
	require.Empty(t, cellValue(exampleData, -1, 1))
	require.Empty(t, cellValue(exampleData, 1, 10))
	require.Empty(t, cellValue(exampleData, 10, 1))
	require.Equal(t, "Domain name", cellValue(exampleData, 1, 1))
}

func TestGenerate(t *testing.T) {
	for _, td := range goldenTests {
		t.Run(td.name, func(t *testing.T) {
			want, err := ioutil.ReadFile(filepath.Join("testdata", "tables", td.name+".md"))
			require.NoError(t, err)
			want = bytes.TrimSuffix(want, []byte{'\n'})
			got := Generate(td.data, td.options...)
			require.Equal(t, string(want), string(got))
		})
	}
}

var exampleData = [][]string{
	{"Date", "Description", "CV2", "Amount"},
	{"1/1/2014", "Domain name", "2233", "$10.98"},
	{"1/1/2014", "January Hosting", "2233", "$54.95"},
	{"1/4/2014", "February Hosting", "2233", "$51.00"},
	{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
}

var goldenTests = []struct {
	name    string
	options []Option
	data    [][]string
}{
	{
		name: "defaults",
		data: exampleData,
	},
	{
		name: "combined-options",
		data: exampleData,
		options: []Option{
			Alignment(AlignCenter),
			ColumnAlignment(1, AlignLeft),
			ColumnAlignment(2, AlignRight),
			ColumnTextAlignment(1, AlignCenter),
			HeaderAlignment(AlignRight),
			ColumnHeaderAlignment(1, AlignLeft),
			ColumnMinWidth(0, 12),
			ColumnMinWidth(2, 12),
			ColumnMinWidth(3, 12),
		},
	},
	{
		name: "empty",
	},
	{
		name:    "Alignment",
		data:    exampleData,
		options: []Option{Alignment(AlignRight)},
	},
	{
		name:    "HeaderAlignment",
		data:    exampleData,
		options: []Option{HeaderAlignment(AlignRight)},
	},
	{
		name:    "TextAlignment",
		data:    exampleData,
		options: []Option{TextAlignment(AlignRight)},
	},
	{
		name:    "ColumnAlignment",
		data:    exampleData,
		options: []Option{ColumnAlignment(1, AlignRight)},
	},
	{
		name:    "ColumnTextAlignment",
		data:    exampleData,
		options: []Option{ColumnTextAlignment(1, AlignRight)},
	},
	{
		name:    "ColumnHeaderAlignment",
		data:    exampleData,
		options: []Option{ColumnHeaderAlignment(1, AlignRight)},
	},
	{
		name:    "ColumnMinWidth",
		data:    exampleData,
		options: []Option{ColumnMinWidth(0, 40)},
	},
	{
		name:    "ColumnMinWidth-small",
		data:    exampleData,
		options: []Option{ColumnMinWidth(0, 2)},
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
	for _, gt := range goldenTests {
		got := Generate(gt.data, gt.options...)
		got = append(got, '\n')
		err = ioutil.WriteFile(filepath.Join("testdata", "tables", gt.name+".md"), got, 0o600)
		if err != nil {
			return err
		}
	}
	return nil
}
