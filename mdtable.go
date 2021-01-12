package mdtable

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Align determines column alignment
type Align uint8

// Align values
const (
	AlignDefault Align = iota
	AlignLeft
	AlignRight
	AlignCenter
)

func (a Align) String() string {
	switch a {
	case AlignDefault:
		return "AlignDefault"
	case AlignLeft:
		return "AlignLeft"
	case AlignRight:
		return "AlignRight"
	case AlignCenter:
		return "AlignCenter"
	default:
		return "Invalid"
	}
}

// DefaultTextAlignment is the alignment used for text alignment on columns set to AlignDefault
const DefaultTextAlignment = AlignLeft

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
		align = DefaultTextAlignment
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

// Table is a markdown table
type Table struct {
	data             [][]string
	mdAlignment      Align
	textAlignment    Align
	headerAlignment  Align
	mdAlignments     []Align
	textAlignments   []Align
	headerAlignments []Align
	minWidths        []int
}

// New returns a new *Table
func New(data [][]string) *Table {
	return &Table{
		data: data,
	}
}

// Data returns the cell values for this table.
//
// The top level slice are rows. The first row is the header data.
// Second level slices are values for each column.
func (t *Table) Data() [][]string {
	return t.data
}

// SetData sets the table data. See Data for a description.
func (t *Table) SetData(data [][]string) {
	t.data = data
}

// HeaderAlignment returns the default text alignment for headers in this table.
func (t *Table) HeaderAlignment() Align {
	return t.headerAlignment
}

// SetHeaderAlignment sets the default text alignment for headers in this table.
func (t *Table) SetHeaderAlignment(headerAlignment Align) {
	t.headerAlignment = headerAlignment
}

// TextAlignment returns the default text alignment for non-header cells in this table.
func (t *Table) TextAlignment() Align {
	return t.textAlignment
}

// SetTextAlignment sets the default text alignment for non-header cells in this table.
func (t *Table) SetTextAlignment(align Align) {
	t.textAlignment = align
}

// Alignment returns the markdown alignment for columns in this table
func (t *Table) Alignment() Align {
	return t.mdAlignment
}

// SetAlignment sets alignment for columns in this table
func (t *Table) SetAlignment(align Align) {
	t.mdAlignment = align
}

// ColumnMinWidth returns the minimum width that is set for a column. Returns 0 if non has been set.
func (t *Table) ColumnMinWidth(column int) int {
	if column < 0 {
		return 0
	}
	if len(t.minWidths) < column+1 {
		return 0
	}
	return t.minWidths[column]
}

// SetColumnMinWidth sets the minimum width for a column
func (t *Table) SetColumnMinWidth(column, width int) {
	if column < 0 {
		return
	}
	delta := (column + 1) - len(t.minWidths)
	if delta > 0 {
		t.minWidths = append(t.minWidths, make([]int, delta)...)
	}
	t.minWidths[column] = width
}

// ColumnAlignment returns the markdown alignment for a column
func (t *Table) ColumnAlignment(column int) Align {
	align := getColumnAlignment(t.mdAlignments, column)
	if align != AlignDefault {
		return align
	}
	return t.Alignment()
}

// SetColumnAlignment sets the markdown alignment for a column
func (t *Table) SetColumnAlignment(column int, align Align) {
	t.mdAlignments = setColumnAlignment(t.mdAlignments, column, align)
}

// ColumnTextAlignment returns text alignment for a column
//
// Order of preference:
//  1. value set with SetColumnTextAlignment
//  2. value set with SetTextAlignment
//  3. ColumnAlignment(column)
//  4. DefaultTextAlignment (which is AlignLeft)
func (t *Table) ColumnTextAlignment(column int) Align {
	align := getColumnAlignment(t.textAlignments, column)
	if align != AlignDefault {
		return align
	}
	align = t.TextAlignment()
	if align != AlignDefault {
		return align
	}
	align = t.ColumnAlignment(column)
	if align != AlignDefault {
		return align
	}
	return DefaultTextAlignment
}

// SetColumnTextAlignment sets the text alignment for a column
func (t *Table) SetColumnTextAlignment(column int, align Align) {
	t.textAlignments = setColumnAlignment(t.textAlignments, column, align)
}

// ColumnHeaderAlignment returns the text alignment for a column header
//
// Order or preference:
//  1. value set with SetColumnHeaderAlignment
//  2. value set with SetHeaderAlignment
//  3. ColumnTextAlignment(column)
func (t *Table) ColumnHeaderAlignment(column int) Align {
	align := getColumnAlignment(t.headerAlignments, column)
	if align != AlignDefault {
		return align
	}
	align = t.HeaderAlignment()
	if align != AlignDefault {
		return align
	}
	return t.ColumnTextAlignment(column)
}

// SetColumnHeaderAlignment sets the text alignment for a column header
func (t *Table) SetColumnHeaderAlignment(column int, align Align) {
	t.headerAlignments = setColumnAlignment(t.headerAlignments, column, align)
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

// Render returns the markdown representation of Table
func (t *Table) Render() []byte {
	var buf bytes.Buffer
	if len(t.data) == 0 {
		return buf.Bytes()
	}
	row := 0
	buf.WriteString(t.renderRow(0, t.ColumnHeaderAlignment) + "\n")
	buf.WriteString(t.renderHeaderRow())
	for row = 1; row < len(t.data); row++ {
		buf.WriteString("\n" + t.renderRow(row, t.ColumnTextAlignment))
	}
	return buf.Bytes()
}

func (t *Table) String() string {
	return string(t.Render())
}

func (t *Table) renderColumnHeader(column int) string {
	width := t.columnWidth(column)
	if width == 0 {
		return "--"
	}
	align := t.ColumnAlignment(column)
	return align.headerPrefix() + strings.Repeat("-", width) + align.headerSuffix()
}

func (t *Table) cellValue(row, column int) string {
	if len(t.data) < row+1 {
		return ""
	}
	if len(t.data[row]) < column+1 {
		return ""
	}
	return t.data[row][column]
}

func (t *Table) renderCell(row, column int, alignment Align) string {
	s := t.cellValue(row, column)
	width := t.columnWidth(column)
	return alignment.fillCell(s, width, 1)
}

func (t *Table) renderRow(row int, alignmentFunc func(int) Align) string {
	cells := make([]string, t.ColumnCount())
	for i := range cells {
		cells[i] = t.renderCell(row, i, alignmentFunc(i))
	}
	return "|" + strings.Join(cells, "|") + "|"
}

func (t *Table) renderHeaderRow() string {
	headers := make([]string, t.ColumnCount())
	for i := range headers {
		headers[i] = t.renderColumnHeader(i)
	}
	return "|" + strings.Join(headers, "|") + "|"
}

// ColumnCount returns the number of columns in the table
func (t *Table) ColumnCount() int {
	count := 0
	for _, row := range t.data {
		if len(row) > count {
			count = len(row)
		}
	}
	return count
}

func (t *Table) columnWidth(column int) int {
	width := t.ColumnMinWidth(column)
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
