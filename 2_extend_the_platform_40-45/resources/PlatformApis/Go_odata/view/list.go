package view

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/view/colours"
	"golang.org/x/term"
)

type listComponent[T any] struct {
	controls       []ListControl
	controlsString string
	columns        listColumns[T]
	columnWidths   []int
	columnStrings  []string
	rawData        []T
	formattedData  [][]string
	selected       int
	hasNext        bool
	hasPrevious    bool
}

type ListControl interface {
	GetLabel() string
	GetValue() string
	GetKey() rune
}

type listColumn[T any] struct {
	label   string
	getData func(T) string
}

type listColumns[T any] []listColumn[T]

func CreateListColumns[T any]() listColumns[T] {
	return listColumns[T]{}
}

func (lcs listColumns[T]) WithColumn(label string, getData func(T) string) listColumns[T] {
	return append(lcs, listColumn[T]{
		label:   label,
		getData: getData,
	})
}

func BuildListComponent[T any](controls []ListControl, columns listColumns[T], data []T, hasNext, hasPrevious bool) Component {
	lc := listComponent[T]{
		controls:      controls,
		columns:       columns,
		columnWidths:  make([]int, len(columns)),
		columnStrings: make([]string, len(columns)),
		rawData:       data,
		hasNext:       hasNext,
		hasPrevious:   hasPrevious,
		formattedData: make([][]string, len(data)),
		selected:      0,
	}
	lc.initialiseFormattedData()
	lc.InitialiseControlString()
	return &lc
}

func (li *listComponent[T]) render() {
	div := "|"
	println(getRowString(li.columnStrings, div, colours.ORANGE))
	println(strings.Repeat("-", len(strings.Join(li.columnStrings, div))))
	for i, rs := range li.formattedData {
		c := colours.RESET
		if i == li.selected {
			c = colours.SELECTED_ROW
		}
		fmt.Println(getRowString(rs, div, c))
	}
	fmt.Println(li.controlsString)
}

func getRowString(cells []string, div string, col colours.Color) string {
	fd := fmt.Sprintf("%s%s%s", colours.RESET, div, col)
	fs := fmt.Sprintf("%s%s%s", col, strings.Join(cells, fd), colours.RESET)
	return fs
}

func (li *listComponent[T]) isInteractive() bool {
	return true
}

func (li *listComponent[T]) handleKeyboardInput(c rune, k keyboard.Key) *updateResponse {
	ur := updateResponse{
		doContinue: true,
	}

	if k == keyboard.KeyArrowUp && li.selected > 0 {
		li.selected--
	}
	if k == keyboard.KeyArrowDown && li.selected < len(li.rawData)-1 {
		li.selected++
	}

	return &ur
}

func (li *listComponent[T]) initialiseFormattedData() {

	cellPadding := 2

	totalColWidth, colWidths := li.calculateNaturalColumnWidths(cellPadding)
	li.columnWidths = colWidths

	dividerCount := len(li.columns) - 1
	maxWidth := GetConsoleSize() - dividerCount
	adjMultiplier := min(float32(maxWidth)/float32(totalColWidth), 1.2)

	for i, c := range li.columns {
		aw := int(float32(colWidths[i]) * adjMultiplier)
		li.columnWidths[i] = aw
		li.columnStrings[i] = formatCellString(c.label, aw, cellPadding)
	}

	for ri, r := range li.rawData {
		li.formattedData[ri] = make([]string, len(li.columns))

		for ci, c := range li.columns {
			w := li.columnWidths[ci]
			rs := c.getData(r)
			fs := formatCellString(rs, w, cellPadding)
			li.formattedData[ri][ci] = fs
		}

	}
}

func (li *listComponent[T]) calculateNaturalColumnWidths(cellPadding int) (int, []int) {

	columnWidths := make([]int, len(li.columns))
	totalColumnWidths := 0

	for i, c := range li.columns {
		max := 0
		for _, e := range li.rawData {
			l := len(c.getData(e)) + cellPadding
			if l > max {
				max = l
			}
		}
		columnWidths[i] = max
		totalColumnWidths += max
	}
	return totalColumnWidths, columnWidths
}

func GetConsoleSize() (width int) {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err.Error())
	}
	return w
}

func formatCellString(rawString string, length, cellPadding int) string {
	contentSpace := length - cellPadding

	fs := rawString
	if len(rawString) > contentSpace {
		fs = rawString[:contentSpace-1] + "…"
	}
	leftPadding := cellPadding / 2
	fs = strings.Repeat(" ", leftPadding) + fs
	fs = fmt.Sprintf("%-*s", length, fs)
	return fs
}

type defaultControl struct {
	char     rune
	label    string
	isActive bool
}

func (li *listComponent[T]) getDefaultControls() []defaultControl {
	return []defaultControl{
		{
			char:     '🡒',
			label:    "Next page",
			isActive: li.hasNext,
		},
		{
			char:     '🡐',
			label:    "Previous page",
			isActive: li.hasPrevious,
		},
		{
			char:     '🡓',
			label:    "Next row",
			isActive: true,
		},
		{
			char:     '🡑',
			label:    "Previous row",
			isActive: true,
		},
	}
}

func (li *listComponent[T]) InitialiseControlString() {
	ctlStr := "\n\nCommands\n\n"
	for _, dc := range li.getDefaultControls() {
		c := colours.ORANGE
		cr := colours.RESET
		if !dc.isActive {
			c = colours.GREY
			cr = colours.GREY
		}
		ctlStr += fmt.Sprintf(" %s%s%s : %s\n", c, string(dc.char), cr, dc.label)
	}

	for _, ctl := range li.controls {
		c := colours.ORANGE
		ctlStr += fmt.Sprintf(" %s%s%s : %s\n", c, string(rune(ctl.GetKey())), colours.RESET, ctl.GetLabel())
	}
	li.controlsString = ctlStr
}
