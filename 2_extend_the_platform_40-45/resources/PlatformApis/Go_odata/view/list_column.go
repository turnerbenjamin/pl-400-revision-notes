// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import "errors"

var ErrNilCellStringFunc = errors.New("cellStringGetter cannot be nil")

// ListColumn represents a column in a list view that can display data from
// entities of type T. It provides methods to get the column's header label and
// to render entity data as strings.
type ListColumn[T Entity] interface {
	// Label returns the header text for the column
	Label() string

	// CellString formats and returns the string representation of an entity's
	// data for this column
	CellString(T) string
}

// listColumn is the standard implementation of the ListColumn interface
type listColumn[T Entity] struct {
	// label contains the header text for the column
	label string

	// cellStringGetter is a function that extracts and formats data from an
	// entity
	cellStringGetter func(T) string
}

// NewListColumn creates a new ListColumn with the specified label and cell
// formatting function. Returns an error if the cellStringGetter function is nil.
func NewListColumn[T Entity](label string, cellStringGetter func(T) string) (ListColumn[T], error) {
	if cellStringGetter == nil {
		return nil, ErrNilCellStringFunc
	}

	return &listColumn[T]{
		label:            label,
		cellStringGetter: cellStringGetter,
	}, nil
}

// Label returns the header text for the column
func (lc *listColumn[T]) Label() string {
	return lc.label
}

// CellString applies the column's formatting function to the provided entity
// and returns the resulting string representation
func (lc *listColumn[T]) CellString(entity T) string {
	return lc.cellStringGetter(entity)
}
