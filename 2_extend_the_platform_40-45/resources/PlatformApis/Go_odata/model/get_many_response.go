// Package model provides data structures for working with Dataverse OData
// API responses.
package model

// GetManyResponse is a generic data structure that represents a standard OData
// response when retrieving multiple entities.
//
// The generic type parameter T represents the entity type contained in the
// response.
// This allows for strongly-typed handling of different entity collections while
// maintaining the common OData response structure.
//
// This structure handles the standard OData collection response format which
// includes the array of entities in the "value" property and an optional next
// link for pagination.
type GetManyResponse[T any] struct {
	// Next contains the URL to retrieve the next page of results.
	// It corresponds to the "@odata.nextLink" property in the OData response.
	// This field will be empty when there are no more pages to retrieve.
	Next string `json:"@odata.nextLink"`

	// Data contains the collection of entities returned by the API.
	// It maps to the "value" property in the OData response.
	Data []T `json:"value"`
}
