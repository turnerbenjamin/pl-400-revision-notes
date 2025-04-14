// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/turnerbenjamin/go_odata/utilities"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

const (
	listColDivider               = "|"
	listRowDivider               = "-"
	listCellPadding              = 2
	truncationMarker             = "…"
	listMaxColumnWidthMultiplier = 1.2
	listCommandsSectionHeader    = "\n\nCommands\n\n"
	listNextPageLabel            = "Next page"
	listPreviousPageLabel        = "Previous page"
	rightArrowChar               = "🡒"
	leftArrowChar                = "🡐"
)

var ErrNoData = errors.New("entityList must contain at least one row of data")
var ErrEmptyCols = errors.New("List component requires at least one column")

type Entity interface {
	ID() string
	Label() string
}

type EntityList[T Entity] interface {
	Data() []T
	HasNext() bool
	Next() (EntityList[T], error)
	HasPrevious() bool
	Previous() EntityList[T]
}

type ListComponentOptions[T Entity] struct {
	Controls   []ListControl
	Columns    []ListColumn[T]
	EntityList EntityList[T]
}

type listComponent[T Entity] struct {
	columns            []ListColumn[T]
	columnWidths       []int
	entityList         EntityList[T]
	data               []T
	selected           int
	customControls     []ListControl
	tableHeaderStrings []string
	tableDataStrings   [][]string
	controlsString     string
}

type navigationControl struct {
	key       string
	label     string
	isEnabled bool
}

func BuildListComponent[T Entity](options ListComponentOptions[T]) (InteractiveComponent, error) {
	if len(options.Columns) == 0 {
		return nil, ErrEmptyCols
	}
	lc := listComponent[T]{
		columns:        options.Columns,
		entityList:     options.EntityList,
		customControls: options.Controls,
		selected:       0,
	}
	err := lc.refreshDataAndCalculateLayout()
	return &lc, err
}

func (lc *listComponent[T]) refreshDataAndCalculateLayout() error {
	lc.selected = 0
	lc.data = lc.entityList.Data()

	if err := lc.validateData(); err != nil {
		return err
	}

	lc.initialiseFormattedData()
	lc.initialiseControlString()
	return nil
}

func (lc *listComponent[T]) render() {
	lc.renderTableHeader()
	lc.renderTableRows()
	lc.renderControls()
}

func (lc *listComponent[T]) renderTableHeader() {
	headerRow := lc.buildRowString(lc.tableHeaderStrings, colours.Orange)
	fmt.Println(headerRow)

	separatorLength := 0
	for _, width := range lc.columnWidths {
		separatorLength += width
	}
	separatorLength += len(listColDivider) * (len(lc.columnWidths) - 1)

	fmt.Println(strings.Repeat(listRowDivider, separatorLength))
}

func (lc *listComponent[T]) renderTableRows() {
	for i, rs := range lc.tableDataStrings {
		c := colours.Reset
		if i == lc.selected {
			c = colours.BlueBackground
		}
		fmt.Println(lc.buildRowString(rs, c))
	}
}

func (lc *listComponent[T]) renderControls() {
	fmt.Print(lc.controlsString)
}

func (lc *listComponent[T]) buildRowString(cells []string, colour colours.Colour) string {
	formattedDivider := colours.ApplyColour(listColDivider, colours.Reset)
	return colours.ApplyColour(strings.Join(cells, formattedDivider), colour)
}

func (lc *listComponent[T]) handleKeyboardInput(char rune, key keyboard.Key) (*updateResponse, error) {

	switch key {
	case keyboard.KeyArrowUp:
		return lc.handleArrowUpPressed()
	case keyboard.KeyArrowDown:
		return lc.handleArrowDownPressed()
	case keyboard.KeyArrowLeft:
		return lc.handleArrowLeftPressed()
	case keyboard.KeyArrowRight:
		return lc.handleArrowRightPressed()
	default:
		return lc.handleCustomControlInput(char)
	}
}

func (lc *listComponent[T]) handleArrowUpPressed() (*updateResponse, error) {
	if len(lc.data) == 0 {
		return nil, ErrNoData
	}

	if lc.selected > 0 {
		lc.selected--
	}

	return newUpdateResponse().setContinue(true), nil
}

func (lc *listComponent[T]) handleArrowDownPressed() (*updateResponse, error) {
	if err := lc.validateData(); err != nil {
		return nil, err
	}

	if lc.selected < len(lc.data)-1 {
		lc.selected++
	}

	return newUpdateResponse().setContinue(true), nil
}

func (lc *listComponent[T]) handleArrowLeftPressed() (*updateResponse, error) {
	if !lc.entityList.HasPrevious() {
		return newUpdateResponse().setContinue(true), nil
	}
	lc.entityList = lc.entityList.Previous()
	err := lc.refreshDataAndCalculateLayout()
	if err != nil {
		return nil, err
	}
	return newUpdateResponse().setContinue(true).setFullRefresh(), nil
}

func (lc *listComponent[T]) handleArrowRightPressed() (*updateResponse, error) {
	if !lc.entityList.HasNext() {
		return newUpdateResponse().setContinue(true), nil
	}

	n, err := lc.entityList.Next()
	if err != nil {
		return nil, err
	}

	lc.entityList = n
	err = lc.refreshDataAndCalculateLayout()
	if err != nil {
		return nil, err
	}

	return newUpdateResponse().setContinue(true).setFullRefresh(), nil
}

func (lc *listComponent[T]) handleCustomControlInput(char rune) (*updateResponse, error) {
	if err := lc.validateData(); err != nil {
		return nil, err
	}

	if lc.selected >= len(lc.data) {
		lc.selected = len(lc.data) - 1
	}

	for _, ci := range lc.customControls {
		if ci.GetKey() == char {
			target := lc.data[lc.selected]

			return newUpdateResponse().
				setContinue(false).
				setUserInput(ci.GetValue()).
				setTarget(target.ID()), nil
		}
	}
	return newUpdateResponse().setContinue(true), nil
}

func (lc *listComponent[T]) initialiseFormattedData() {
	lc.columnWidths = lc.buildColumnWidths()
	lc.tableHeaderStrings = lc.buildFormattedTableHeader()
	lc.tableDataStrings = lc.buildFormattedTableData()
}

func (lc *listComponent[T]) buildColumnWidths() []int {
	naturalTableWidth, naturalColumnWidths := lc.calculateNaturalTableDimensions()
	columnWidths := make([]int, len(lc.columns))

	dividerCount := len(lc.columns) - 1
	maxWidth := utilities.GetConsoleWidth() - dividerCount

	adjMultiplier := min(
		float32(maxWidth)/float32(naturalTableWidth),
		listMaxColumnWidthMultiplier,
	)

	for i := range naturalColumnWidths {
		columnWidths[i] = int(float32(naturalColumnWidths[i]) * adjMultiplier)
	}
	return columnWidths
}

func (lc *listComponent[T]) buildFormattedTableHeader() []string {
	tableHeaderStrings := make([]string, len(lc.columns))

	for i, c := range lc.columns {
		tableHeaderStrings[i] = lc.formatCellString(c.Label(), lc.columnWidths[i])
	}
	return tableHeaderStrings
}

func (lc *listComponent[T]) buildFormattedTableData() [][]string {
	tableRows := make([][]string, len(lc.data))
	for i, r := range lc.data {
		tableRows[i] = lc.buildFormattedTableRow(r)
	}
	return tableRows
}

func (lc *listComponent[T]) buildFormattedTableRow(entity T) []string {
	tableRow := make([]string, len(lc.columns))
	for i, col := range lc.columns {
		colWidth := lc.columnWidths[i]
		cellData := col.CellString(entity)
		formattedData := lc.formatCellString(cellData, colWidth)
		tableRow[i] = formattedData
	}
	return tableRow
}

func (lc *listComponent[T]) calculateNaturalTableDimensions() (
	naturalTableWidth int,
	naturalColumnWidths []int,
) {
	naturalColumnWidths = make([]int, len(lc.columns))
	naturalTableWidth = 0

	for i, column := range lc.columns {
		maxCellStringLength := len(column.Label()) + listCellPadding

		for _, rowData := range lc.data {
			cellStringLength := len(column.CellString(rowData)) + listCellPadding

			if cellStringLength > maxCellStringLength {
				maxCellStringLength = cellStringLength
			}
		}
		naturalColumnWidths[i] = maxCellStringLength
		naturalTableWidth += maxCellStringLength
	}
	return naturalTableWidth, naturalColumnWidths
}

func (lc *listComponent[T]) formatCellString(rawString string, cellWidth int) string {
	contentWidth := cellWidth - listCellPadding
	content := lc.truncatedString(rawString, contentWidth)
	leftPadding := listCellPadding / 2
	return lc.paddedString(content, leftPadding, cellWidth)
}

func (lc *listComponent[T]) truncatedString(s string, maxLength int) string {
	ts := s
	if len(ts) > maxLength {
		ts = ts[:maxLength-len(truncationMarker)] + truncationMarker
	}
	return ts
}

func (lc *listComponent[T]) paddedString(s string, leftPadding, totalLength int) string {
	ps := strings.Repeat(" ", leftPadding) + s
	return fmt.Sprintf("%-*s", totalLength, ps)
}

func (lc *listComponent[T]) initialiseControlString() {
	var builder strings.Builder
	builder.WriteString(listCommandsSectionHeader)
	builder.WriteString(lc.buildNavigationControlsString())
	builder.WriteString(lc.buildCustomControlsString())

	lc.controlsString = builder.String()
}

func (lc *listComponent[T]) buildNavigationControlsString() string {
	var builder strings.Builder
	controls := lc.getNavigationControls()
	for _, ctl := range controls {
		builder.WriteString(lc.getControlString(ctl.key, ctl.label, ctl.isEnabled))
	}
	return builder.String()
}

func (lc *listComponent[T]) buildCustomControlsString() string {
	var builder strings.Builder
	for _, ctl := range lc.customControls {
		builder.WriteString(lc.getControlString(string(ctl.GetKey()), ctl.GetLabel(), true))
	}
	return builder.String()
}

func (lc *listComponent[T]) getControlString(key, label string, isActive bool) string {
	keyString := colours.ApplyColour(key, colours.Orange)
	labelString := fmt.Sprintf(" : %s\n", label)

	if !isActive {
		keyString = colours.ApplyColour(key, colours.Grey)
		labelString = colours.ApplyColour(labelString, colours.Grey)
	}
	return keyString + labelString
}

func (lc *listComponent[T]) getNavigationControls() []navigationControl {
	return []navigationControl{
		{
			key:       rightArrowChar,
			label:     listNextPageLabel,
			isEnabled: lc.entityList.HasNext(),
		},
		{
			key:       leftArrowChar,
			label:     listPreviousPageLabel,
			isEnabled: lc.entityList.HasPrevious(),
		},
	}
}

func (lc *listComponent[T]) validateData() error {
	if len(lc.data) == 0 {
		return ErrNoData
	}
	return nil
}
