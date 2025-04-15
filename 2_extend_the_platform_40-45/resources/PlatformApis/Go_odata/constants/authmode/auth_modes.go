// Package authmode defines authentication mode constants used when
// connecting to Microsoft Dataverse.
package authmode

// AuthenticationMode represents the type of authentication flow to use
// when acquiring tokens for API access. It's implemented as a string type
// for type safety when working with authentication options.
type AuthenticationMode string

// Authentication mode constants define the available authentication flows.
const (
	// Application represents authentication using client credentials
	// (app permissions)
	Application AuthenticationMode = "Application"

	// User represents interactive authentication with user delegation
	User AuthenticationMode = "User"

	// Invalid represents an invalid authentication mode selection
	Invalid AuthenticationMode = "Invalid"
)
