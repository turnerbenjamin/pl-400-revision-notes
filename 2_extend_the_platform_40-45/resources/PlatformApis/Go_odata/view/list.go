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
	defaultConsoleWidth          = 80
)

// ErrNoData is returned when attempting to render a list with no data rows.
var ErrNoData = errors.New("entityList must contain at least one row of data")

// ErrEmptyCols is returned when attempting to create a list component with no
// columns.
var ErrEmptyCols = errors.New("list component requires at least one column")

// Entity represents an item that can be displayed in a list component.
// It must provide an ID and a human-readable label.
type Entity interface {
	ID() string
	Label() string
}

// EntityList represents a pageable collection of Entity objects.
// It provides methods to access the current page data and navigate between
// pages.
type EntityList[T Entity] interface {
	// Data returns the entities on the current page.
	Data() []T

	// HasNext returns true if there are more pages after the current one.
	HasNext() bool

	// Next returns the next page of entities.
	Next() (EntityList[T], error)

	// HasPrevious returns true if there are previous pages before the current
	// one.
	HasPrevious() bool

	// Previous returns the previous page of entities.
	Previous() EntityList[T]
}

// ListComponentOptions configures the behaviour and appearance of a list
// component.
type ListComponentOptions[T Entity] struct {
	// Controls define custom keyboard actions available in the list.
	Controls []ListControl

	// Columns define what data is displayed and how it's formatted.
	Columns []ListColumn[T]

	// EntityList provides the data to be displayed.
	EntityList EntityList[T]
}

// listComponent implements an interactive, terminal-based data table with
// navigation controls.
type listComponent[T Entity] struct {
	// columns defines the data columns to display
	columns []ListColumn[T]

	// columnWidths stores calculated display widths for each column
	columnWidths []int

	// entityList provides the underlying data and pagination capabilities
	entityList EntityList[T]

	// data holds the current page of entities being displayed
	data []T

	// selected indicates the currently highlighted row index
	selected int

	// customControls defines keyboard commands available to the user
	customControls []ListControl

	// tableHeaderStrings contains formatted column headers
	tableHeaderStrings []string

	// tableDataStrings contains formatted cell values for all rows
	tableDataStrings [][]string

	// controlsString contains the formatted help text for keyboard controls
	controlsString string
}

// navigationControl represents a keyboard navigation command that allows users
// to move between pages in the list component. It stores the key to press,
// the descriptive label for the command, and whether the control is currently
// enabled based on pagination state.
type navigationControl struct {
	key       string // Character representing the keyboard command
	label     string // Text description shown to the user
	isEnabled bool   // Whether this navigation option is currently available
}

// BuildListComponent creates a new interactive list component that displays
// entity data in a paginated table format with keyboard navigation.
// Returns an error if the options are invalid (no columns or no data).
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

// refreshDataAndCalculateLayout resets the component state and recalculates
// layout.
// It retrieves the current page of data, validates it, and initializes
// formatting.
// Returns an error if data validation fails.
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

// render displays the full list component in the terminal.
// This includes the table header, data rows, and control instructions.
func (lc *listComponent[T]) render() {
	lc.renderTableHeader()
	lc.renderTableRows()
	lc.renderControls()
}

// renderTableHeader displays the column headers with appropriate formatting and
// draws a separator line beneath them.
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

// renderTableRows displays all data rows, highlighting the currently selected
// row with a blue background.
func (lc *listComponent[T]) renderTableRows() {
	for i, rs := range lc.tableDataStrings {
		c := colours.Reset
		if i == lc.selected {
			c = colours.BlueBackground
		}
		fmt.Println(lc.buildRowString(rs, c))
	}
}

// renderControls displays the available keyboard commands and their labels at
// the bottom of the component.
func (lc *listComponent[T]) renderControls() {
	fmt.Print(lc.controlsString)
}

// buildRowString joins an array of cell strings with column dividers and
// applies the specified color formatting.
func (lc *listComponent[T]) buildRowString(cells []string, colour colours.Colour) string {
	formattedDivider := colours.ApplyColour(listColDivider, colours.Reset)
	return colours.ApplyColour(strings.Join(cells, formattedDivider), colour)
}

// handleKeyboardInput processes keyboard events and returns an appropriate
// response.
// Handles arrow keys for navigation and custom control keys.
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

// handleArrowUpPressed moves selection to the previous row if available.
// Returns an error if there's no data to navigate.
func (lc *listComponent[T]) handleArrowUpPressed() (*updateResponse, error) {
	if len(lc.data) == 0 {
		return nil, ErrNoData
	}

	if lc.selected > 0 {
		lc.selected--
	}

	return newUpdateResponse().setContinue(true), nil
}

// handleArrowDownPressed moves selection to the next row if available.
// Returns an error if there's no data to navigate.
func (lc *listComponent[T]) handleArrowDownPressed() (*updateResponse, error) {
	if err := lc.validateData(); err != nil {
		return nil, err
	}

	if lc.selected < len(lc.data)-1 {
		lc.selected++
	}

	return newUpdateResponse().setContinue(true), nil
}

// handleArrowLeftPressed navigates to the previous page of data if available.
// Refreshes the component layout after changing pages.
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

// handleArrowRightPressed navigates to the next page of data if available.
// Refreshes the component layout after changing pages.
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

// handleCustomControlInput processes custom key commands for the currently
// selected item. Returns a response with the command value and target entity ID
// if a valid key is pressed.
func (lc *listComponent[T]) handleCustomControlInput(char rune) (*updateResponse, error) {
	if err := lc.validateData(); err != nil {
		return nil, err
	}

	if lc.selected >= len(lc.data) {
		lc.selected = len(lc.data) - 1
	}

	for _, ci := range lc.customControls {
		if ci.Key() == char {
			target := lc.data[lc.selected]

			return newUpdateResponse().
				setContinue(false).
				setUserInput(ci.Value()).
				setTarget(target.ID()), nil
		}
	}
	return newUpdateResponse().setContinue(true), nil
}

// initialiseFormattedData prepares all display data by calculating column
// widths and formatting header and row content.
func (lc *listComponent[T]) initialiseFormattedData() {
	lc.columnWidths = lc.buildColumnWidths()
	lc.tableHeaderStrings = lc.buildFormattedTableHeader()
	lc.tableDataStrings = lc.buildFormattedTableData()
}

// buildColumnWidths calculates the optimal width for each column based on
// content and available terminal width.
func (lc *listComponent[T]) buildColumnWidths() []int {
	naturalTableWidth, naturalColumnWidths := lc.calculateNaturalTableDimensions()
	columnWidths := make([]int, len(lc.columns))

	dividerCount := len(lc.columns) - 1
	maxWidth := utilities.GetConsoleWidth(defaultConsoleWidth) - dividerCount

	adjMultiplier := min(
		float32(maxWidth)/float32(naturalTableWidth),
		listMaxColumnWidthMultiplier,
	)

	for i := range naturalColumnWidths {
		columnWidths[i] = int(float32(naturalColumnWidths[i]) * adjMultiplier)
	}
	return columnWidths
}

// buildFormattedTableHeader creates formatted strings for each column header.
func (lc *listComponent[T]) buildFormattedTableHeader() []string {
	tableHeaderStrings := make([]string, len(lc.columns))

	for i, c := range lc.columns {
		tableHeaderStrings[i] = lc.formatCellString(c.Label(), lc.columnWidths[i])
	}
	return tableHeaderStrings
}

// buildFormattedTableData creates a 2D array of formatted strings for all data
// cells.
func (lc *listComponent[T]) buildFormattedTableData() [][]string {
	tableRows := make([][]string, len(lc.data))
	for i, r := range lc.data {
		tableRows[i] = lc.buildFormattedTableRow(r)
	}
	return tableRows
}

// buildFormattedTableRow formats a single entity's data into a row of strings.
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

// calculateNaturalTableDimensions determines the natural width needed for each
// column based on its content and returns both individual column widths and
// total table width.
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

// formatCellString prepares a cell's content by truncating if necessary and
// adding padding.
func (lc *listComponent[T]) formatCellString(rawString string, cellWidth int) string {
	contentWidth := cellWidth - listCellPadding
	content := lc.truncatedString(rawString, contentWidth)
	leftPadding := listCellPadding / 2
	return lc.paddedString(content, leftPadding, cellWidth)
}

// truncatedString ensures a string doesn't exceed the specified length by
// truncating and adding an ellipsis if necessary.
func (lc *listComponent[T]) truncatedString(s string, maxLength int) string {
	ts := s
	if len(ts) > maxLength {
		ts = ts[:maxLength-len(truncationMarker)] + truncationMarker
	}
	return ts
}

// paddedString adds left padding to a string and ensures it fills exactly the
// specified width.
func (lc *listComponent[T]) paddedString(s string, leftPadding, totalLength int) string {
	ps := strings.Repeat(" ", leftPadding) + s
	return fmt.Sprintf("%-*s", totalLength, ps)
}

// initialiseControlString creates the formatted string displaying all available
// keyboard commands for the component.
func (lc *listComponent[T]) initialiseControlString() {
	var builder strings.Builder
	builder.WriteString(listCommandsSectionHeader)
	builder.WriteString(lc.buildNavigationControlsString())
	builder.WriteString(lc.buildCustomControlsString())

	lc.controlsString = builder.String()
}

// buildNavigationControlsString formats the pagination control instructions.
func (lc *listComponent[T]) buildNavigationControlsString() string {
	var builder strings.Builder
	controls := lc.getNavigationControls()
	for _, ctl := range controls {
		builder.WriteString(lc.getControlString(ctl.key, ctl.label, ctl.isEnabled))
	}
	return builder.String()
}

// buildCustomControlsString formats the custom control instructions.
func (lc *listComponent[T]) buildCustomControlsString() string {
	var builder strings.Builder
	for _, ctl := range lc.customControls {
		builder.WriteString(lc.getControlString(string(ctl.Key()), ctl.Label(), true))
	}
	return builder.String()
}

// getControlString formats a single control instruction with appropriate
// colours based on whether the control is currently active.
func (lc *listComponent[T]) getControlString(key, label string, isActive bool) string {
	keyString := colours.ApplyColour(key, colours.Orange)
	labelString := fmt.Sprintf(" : %s\n", label)

	if !isActive {
		keyString = colours.ApplyColour(key, colours.Grey)
		labelString = colours.ApplyColour(labelString, colours.Grey)
	}
	return keyString + labelString
}

// getNavigationControls returns the array of navigation controls with their
// current enabled state.
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

// validateData ensures the component has data to display.
// Returns ErrNoData if the data slice is empty.
func (lc *listComponent[T]) validateData() error {
	if len(lc.data) == 0 {
		return ErrNoData
	}
	return nil
}
