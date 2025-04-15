// Package msal provides authentication mechanisms for Microsoft identity
// platform. It contains implementations of the DataverseClient interface for
// different authentication flows.
package msal

import (
	"context"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

// appClient implements the DataverseClient interface using a confidential
// application client from the MSAL library.
type appClient struct {
	client      *confidential.Client
	resourceURL string
}

// GetAppService creates and returns a DataverseClient implementation using
// the provided ClientOptions. This client can be used to authenticate and
// communicate with Microsoft Dataverse services.
//
// The returned client uses confidential application authentication flow.
// It acquires tokens using the client credentials grant type.
//
// Parameters:
//   - c: ClientOptions containing required authentication parameters
//
// Returns:
//   - A DataverseClient implementation
//   - An error if client creation fails
func GetAppService(c ClientOptions) (DataverseClient, error) {
	cred, err := confidential.NewCredFromSecret(c.ClientSecret)
	if err != nil {
		return nil, err
	}

	app, err := confidential.New(c.Authority, c.ClientID, cred)
	if err != nil {
		return nil, err
	}

	return &appClient{
		client:      &app,
		resourceURL: c.ResourceURL,
	}, nil
}

// AcquireToken obtains an access token for accessing the Dataverse API.
// It uses the client credentials flow to acquire a token with the default
// scope for the configured resource URL.
//
// Returns:
//   - The access token as a string
//   - An error if token acquisition fails
func (c *appClient) AcquireToken() (string, error) {

	resourceURL := c.resourceURL
	if len(resourceURL) > 0 && !strings.HasSuffix(resourceURL, "/") {
		resourceURL += "/"
	}
	scopes := []string{resourceURL + ".default"}

	result, err := c.client.AcquireTokenByCredential(context.TODO(), scopes)
	if err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
