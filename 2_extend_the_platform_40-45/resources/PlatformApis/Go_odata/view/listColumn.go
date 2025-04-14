// Package view provides UI components for terminal-based applications.
// It includes interactive elements like inputs, lists, and navigation controls.
package view

import "errors"

var ErrNilCellStringFunc = errors.New("getCellString function cannot be nil")

type ListColumn[T Entity] interface {
	Label() string
	CellString(T) string
}

type listColumn[T Entity] struct {
	label            string
	cellStringGetter func(T) string
}

func NewListColumn[T Entity](label string, cellStringGetter func(T) string) (ListColumn[T], error) {
	if cellStringGetter == nil {
		return nil, ErrNilCellStringFunc
	}

	return &listColumn[T]{
		label:            label,
		cellStringGetter: cellStringGetter,
	}, nil
}

func (lc *listColumn[T]) Label() string {
	return lc.label
}

func (lc *listColumn[T]) CellString(entity T) string {
	return lc.cellStringGetter(entity)
}
