/*
Package mdtable generates markdown tables from string slices with formatting options for alignment and column width.

*/
package mdtable

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Align is a value for markdown and text alignment
type Align uint8

// Align values
const (
	AlignDefault Align = iota // markdown: |--------|  text: | foo     |
	AlignLeft                 // markdown: |:-------|  text: | foo     |
	AlignRight                // markdown: |-------:|  text: |     foo |
	AlignCenter               // markdown: |:------:|  text: |   foo   |
)

// defaultTextAlignment is the alignment used for text alignment on columns set to AlignDefault
const defaultTextAlignment = AlignLeft

func (a Align) headerPrefix() string {
	switch a {
	case AlignLeft, AlignCenter:
		return ":"
	default:
		return "-"
	}
}

func (a Align) headerSuffix() string {
	switch a {
	case AlignRight, AlignCenter:
		return ":"
	default:
		return "-"
	}
}

func (a Align) fillCell(s string, width, padding int) string {
	align := a
	if align == AlignDefault {
		align = defaultTextAlignment
	}
	pad := strings.Repeat(" ", padding)
	var leftFill, rightFill int
	switch align {
	case AlignCenter:
		filled := runewidth.FillLeft(s, width)
		delta := runewidth.StringWidth(filled) - runewidth.StringWidth(s)
		leftFill = delta / 2
		rightFill = delta/2 + delta%2
		s = strings.Repeat(" ", leftFill) + s + strings.Repeat(" ", rightFill)
	case AlignRight:
		s = runewidth.FillLeft(s, width)
	case AlignLeft:
		s = runewidth.FillRight(s, width)
	}
	return pad + s + pad
}

/*
Generate generates a markdown table.

data -- is the table data. The top level slice contains rows and second level slices are cells.
The first row is the header row. You need at least two rows (including the header row) to create
a valid markdown table.

options -- are options for formatting the table. All options are in options.md.

There are three different types of alignment available, Alignment, TextAlignment and HeaderAlignment.
Each of these can be set at either the table or column level.

Alignment -- controls how the markdown will appear when rendered in a browser.

TextAlignment -- controls the text alignment of fields. Default behavior is using the Alignment value
of AlignLeft if that is unset.

HeaderAlignment -- is the same as TextAlignment, but for the header row only. Default behavior is
using the TextAlignment value.

*/
func Generate(data [][]string, options ...Option) []byte {
	t := &table{
		data: data,
	}

	for _, o := range options {
		o(t)
	}

	return t.generate()
}

// Option is an option to control a table's formatting
type Option func(*table)

// Start Options

// HeaderAlignment sets the text alignment for headers
func HeaderAlignment(val Align) Option {
	return func(t *table) {
		t.headerAlignment = val
	}
}

// TextAlignment sets the default text alignment for non-header cells
func TextAlignment(val Align) Option {
	return func(t *table) {
		t.textAlignment = val
	}
}

// Alignment sets alignment for columns.
func Alignment(val Align) Option {
	return func(t *table) {
		t.mdAlignment = val
	}
}

// ColumnAlignment sets the markdown alignment for a column
func ColumnAlignment(column int, alignment Align) Option {
	return func(t *table) {
		t.mdAlignments = setColumnAlignment(t.mdAlignments, column, alignment)
	}
}

// ColumnHeaderAlignment sets the text alignment for a column header
func ColumnHeaderAlignment(column int, alignment Align) Option {
	return func(t *table) {
		t.headerAlignments = setColumnAlignment(t.headerAlignments, column, alignment)
	}
}

// ColumnTextAlignment sets the text alignment for a column
func ColumnTextAlignment(column int, alignment Align) Option {
	return func(t *table) {
		t.textAlignments = setColumnAlignment(t.textAlignments, column, alignment)
	}
}

// ColumnMinWidth sets the minimum width for a column
func ColumnMinWidth(column, width int) Option {
	return func(t *table) {
		if column < 0 {
			return
		}
		delta := (column + 1) - len(t.minWidths)
		if delta > 0 {
			t.minWidths = append(t.minWidths, make([]int, delta)...)
		}
		t.minWidths[column] = width
	}
}

// End Options

// table is a markdown table
type table struct {
	data             [][]string
	mdAlignment      Align
	textAlignment    Align
	headerAlignment  Align
	mdAlignments     []Align
	textAlignments   []Align
	headerAlignments []Align
	minWidths        []int
}

// columnMinWidth returns the minimum width that is set for a column. Returns 0 if non has been set.
func (t *table) columnMinWidth(column int) int {
	if column < 0 {
		return 0
	}
	if len(t.minWidths) < column+1 {
		return 0
	}
	return t.minWidths[column]
}

// columnAlignment returns the markdown alignment for a column
func (t *table) columnAlignment(column int) Align {
	align := getColumnAlignment(t.mdAlignments, column)
	if align != AlignDefault {
		return align
	}
	return t.mdAlignment
}

// columnTextAlignment returns text alignment for a column
//
// Order of preference:
//  1. value set with ColumnTextAlignment
//  2. value set with TextAlignment
//  3. columnAlignment(column)
//  4. defaultTextAlignment (which is AlignLeft)
func (t *table) columnTextAlignment(column int) Align {
	align := getColumnAlignment(t.textAlignments, column)
	if align != AlignDefault {
		return align
	}
	align = t.textAlignment
	if align != AlignDefault {
		return align
	}
	align = t.columnAlignment(column)
	if align != AlignDefault {
		return align
	}
	return defaultTextAlignment
}

// columnHeaderAlignment returns the text alignment for a column header
//
// Order or preference:
//  1. value set with SetColumnHeaderAlignment
//  2. value set with SetHeaderAlignment
//  3. ColumnTextAlignment(column)
func (t *table) columnHeaderAlignment(column int) Align {
	align := getColumnAlignment(t.headerAlignments, column)
	if align != AlignDefault {
		return align
	}
	align = t.headerAlignment
	if align != AlignDefault {
		return align
	}
	return t.columnTextAlignment(column)
}

func getColumnAlignment(alignments []Align, column int) Align {
	if column < 0 {
		return AlignDefault
	}
	if len(alignments) < column+1 {
		return AlignDefault
	}
	return alignments[column]
}

func setColumnAlignment(alignments []Align, column int, align Align) []Align {
	if column < 0 {
		return alignments
	}
	delta := (column + 1) - len(alignments)
	if delta > 0 {
		alignments = append(alignments, make([]Align, delta)...)
	}
	alignments[column] = align
	return alignments
}

// generate returns the markdown representation of table
func (t *table) generate() []byte {
	var buf bytes.Buffer
	if len(t.data) == 0 {
		return buf.Bytes()
	}
	row := 0
	buf.WriteString(t.renderRow(0, t.columnHeaderAlignment) + "\n")
	buf.WriteString(t.renderHeaderRow())
	for row = 1; row < len(t.data); row++ {
		buf.WriteString("\n" + t.renderRow(row, t.columnTextAlignment))
	}
	return buf.Bytes()
}

func (t *table) renderColumnHeader(column int) string {
	width := t.columnWidth(column)
	if width == 0 {
		return "--"
	}
	align := t.columnAlignment(column)
	return align.headerPrefix() + strings.Repeat("-", width) + align.headerSuffix()
}

func cellValue(data [][]string, row, column int) string {
	if row < 0 || column < 0 {
		return ""
	}
	if len(data) < row+1 {
		return ""
	}
	if len(data[row]) < column+1 {
		return ""
	}
	return data[row][column]
}

func (t *table) renderCell(row, column int, alignment Align) string {
	s := cellValue(t.data, row, column)
	width := t.columnWidth(column)
	return alignment.fillCell(s, width, 1)
}

func (t *table) renderRow(row int, alignmentFunc func(int) Align) string {
	cells := make([]string, t.columnCount())
	for i := range cells {
		cells[i] = t.renderCell(row, i, alignmentFunc(i))
	}
	return "|" + strings.Join(cells, "|") + "|"
}

func (t *table) renderHeaderRow() string {
	headers := make([]string, t.columnCount())
	for i := range headers {
		headers[i] = t.renderColumnHeader(i)
	}
	return "|" + strings.Join(headers, "|") + "|"
}

// columnCount returns the number of columns in the table
func (t *table) columnCount() int {
	count := 0
	for _, row := range t.data {
		if len(row) > count {
			count = len(row)
		}
	}
	return count
}

func (t *table) columnWidth(column int) int {
	width := t.columnMinWidth(column)
	for _, row := range t.data {
		if len(row) < column+1 {
			continue
		}
		strLen := len(row[column])
		if strLen > width {
			width = strLen
		}
	}
	return width
}
