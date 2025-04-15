// Package model provides data structures for working with Dataverse OData
// API responses.
package model

import "github.com/turnerbenjamin/go_odata/view"

// entityList implements the view.EntityList interface and provides
// functionality for navigating paginated collections of entities from
// Dataverse OData responses.
type entityList[T view.Entity] struct {
	// fetchNext is a function that retrieves the next page of results using the
	// provided URL
	fetchNext func(string) (*GetManyResponse[T], error)

	// data contains the current page of entity records
	data []T

	// next contains the URL to fetch the next page of results, or empty
	// string if there are no more pages
	next string

	// previous references the previous page in the collection, or nil if this
	// is the first page
	previous view.EntityList[T]
}

// CreateEntityList creates a new EntityList from an initial GetManyResponse.
// It implements the view.EntityList interface for navigating paginated
// collections.
//
// Parameters:
//   - getManyResponse: The initial page of results from the OData API
//   - fetchNext: A function that will be called to retrieve subsequent pages
//
// Returns:
//   - An EntityList implementation that provides access to the current page
//     and navigation to other pages
func CreateEntityList[T view.Entity](getManyResponse GetManyResponse[T], fetchNext func(string) (*GetManyResponse[T], error)) view.EntityList[T] {
	return &entityList[T]{
		fetchNext: fetchNext,
		data:      getManyResponse.Data,
		next:      getManyResponse.Next,
		previous:  nil,
	}
}

// Data returns the current page of entity records
func (el *entityList[T]) Data() []T {
	return el.data
}

// HasNext returns true if there is a next page of results available.
func (el *entityList[T]) HasNext() bool {
	return el.next != ""
}

// Next retrieves the next page of results and returns it as a new EntityList.
// The current EntityList is set as the previous page in the returned list.
//
// Returns:
//   - A new EntityList containing the next page of results
//   - An error if the next page could not be retrieved
func (el *entityList[T]) Next() (view.EntityList[T], error) {
	r, err := el.fetchNext(el.next)
	if err != nil {
		return nil, err
	}

	return &entityList[T]{
		fetchNext: el.fetchNext,
		data:      r.Data,
		next:      r.Next,
		previous:  el,
	}, nil
}

// HasPrevious returns true if there is a previous page available.
func (el *entityList[T]) HasPrevious() bool {
	return el.previous != nil
}

// Previous returns the previous page in the result set.
// If there is no previous page (this is the first page),
// it returns nil.
func (el *entityList[T]) Previous() view.EntityList[T] {
	return el.previous
}
