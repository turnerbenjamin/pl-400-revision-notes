package msal

import (
	"context"
	"fmt"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type delegatedClient struct {
	client      *public.Client
	resourceUrl string
}

func GetDelegatedService(c ClientOptions) (DataverseClient, error) {
	client, err := public.New(c.ClientId, public.WithAuthority(c.Authority))
	if err != nil {
		return nil, err
	}

	return &delegatedClient{
		client:      &client,
		resourceUrl: c.ResourceUrl,
	}, nil
}

func (c *delegatedClient) AcquireToken() (string, error) {
	scopes := []string{c.resourceUrl + "user_impersonation"}

	// Looks for a token in the cache
	accounts, err := c.client.Accounts(context.TODO())
	if err == nil && len(accounts) > 0 {
		response, err := c.client.AcquireTokenSilent(context.TODO(), scopes, public.WithSilentAccount(accounts[0]))
		if err == nil {
			return response.AccessToken, nil
		}
	}
	// Opens default browser so client can authenticate
	fmt.Printf("\nPlease authenticate in the browser...\n")
	response, err := c.client.AcquireTokenInteractive(context.TODO(), scopes)

	if err != nil {
		return "", nil
	}
	return response.AccessToken, nil
}
