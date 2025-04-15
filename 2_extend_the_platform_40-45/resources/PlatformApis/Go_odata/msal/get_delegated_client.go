// Package msal provides authentication mechanisms for Microsoft identity platform.
// It contains implementations of the DataverseClient interface for different
// authentication flows.
package msal

import (
	"context"
	"fmt"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

// delegatedClient implements the DataverseClient interface using a delegated
// authentication flow with MSAL. It manages token acquisition for
// user-delegated permissions to access Dataverse.
type delegatedClient struct {
	client      *public.Client
	resourceURL string
}

// GetDelegatedService creates a new DataverseClient that uses delegated
// authentication. It requires user authentication through a browser. The
// returned client implements the DataverseClient interface defined in
// dataverse_service.go.
func GetDelegatedService(c ClientOptions) (DataverseClient, error) {
	client, err := public.New(c.ClientID, public.WithAuthority(c.Authority))
	if err != nil {
		return nil, err
	}

	return &delegatedClient{
		client:      &client,
		resourceURL: c.ResourceURL,
	}, nil
}

// AcquireToken obtains an access token for Dataverse API access.
// It first attempts to acquire a token silently from the cache.
// If that fails, it initiates an interactive authentication flow that opens a
// browser for user login.
// Returns the access token as a string or an error if authentication fails.
func (c *delegatedClient) AcquireToken() (string, error) {

	resourceURL := c.resourceURL
	if len(resourceURL) > 0 && !strings.HasSuffix(resourceURL, "/") {
		resourceURL += "/"
	}
	scopes := []string{resourceURL + "user_impersonation"}

	// Looks for a token in the cache
	accounts, err := c.client.Accounts(context.TODO())
	if err == nil && len(accounts) > 0 {
		response, err := c.client.AcquireTokenSilent(
			context.TODO(), scopes,
			public.WithSilentAccount(accounts[0]))
		if err == nil {
			return response.AccessToken, nil
		}
	}
	// Opens default browser so client can authenticate
	fmt.Printf("\nPlease authenticate in the browser...\n")
	response, err := c.client.AcquireTokenInteractive(context.TODO(), scopes)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}
