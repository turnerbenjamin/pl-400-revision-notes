// Package msal provides authentication mechanisms for accessing Microsoft Dataverse
// through the Microsoft Authentication Library (MSAL).
package msal

// ClientOptions contains the configuration parameters needed to establish
// authenticated connections to Microsoft Dataverse services.
//
// It supports both application (client credentials) and delegated
// (user interactive) authentication flows depending on which client
// implementation is used.
type ClientOptions struct {
	// ClientID is the application (client) ID registered in Azure AD
	ClientID string

	// ClientSecret is the client secret for confidential client applications
	// This is only required for application authentication flows
	ClientSecret string

	// ResourceURL is the base URL of the Dataverse environment
	// For example: "https://orgname.crm.dynamics.com/"
	ResourceURL string

	// Authority is the Azure AD authority URL
	// For example: "https://login.microsoftonline.com/tenant-id"
	Authority string
}

// DataverseClient defines the interface for authentication providers that
// can obtain access tokens for Microsoft Dataverse.
//
// This interface is implemented by different client types that support
// various authentication flows, such as client credentials (application
// permissions) and interactive browser authentication (delegated permissions).
type DataverseClient interface {
	// AcquireToken obtains an access token for authenticating with the
	// Dataverse API. The token can be used in Authorization headers for API
	// requests.
	//
	// Returns:
	//   - The access token as a string
	//   - An error if token acquisition fails
	AcquireToken() (string, error)
}
