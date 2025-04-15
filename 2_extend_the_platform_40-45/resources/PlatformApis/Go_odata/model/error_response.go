// Package model provides data structures for working with Dataverse OData
// API responses.
package model

// ErrorResponse represents a standardized error response from the Dataverse
// OData API.
// This structure follows the OData v4 error response format, where error
// details are nested within an "error" object.
type ErrorResponse struct {
	// Error contains the details of the error returned by the API.
	// This matches the OData v4 error format specification.
	Error struct {
		// Code is the error code returned by the Dataverse API.
		// This typically represents the type or category of the error.
		Code string `json:"code"`

		// Message contains a human-readable description of the error.
		// This provides more detailed information about what went wrong.
		Message string `json:"message"`
	} `json:"error"`
}
