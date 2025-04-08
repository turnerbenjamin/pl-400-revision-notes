package msal

import (
	"context"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type delegatedClient struct {
	client      *public.Client
	resourceUrl string
}

func GetDelegatedService(c DataverseServiceConfig) (DataverseService, error) {
	client, err := public.New(c.ClientId, public.WithAuthority(c.Authority))
	return &dataverseService{
		client: &delegatedClient{
			client:      &client,
			resourceUrl: c.ResourceUrl,
		},
		baseurl: c.APIBaseUrl,
	}, err
}

func (c *delegatedClient) AcquireToken() (string, error) {
	// Looks for a token in the cache
	scopes := []string{c.resourceUrl + "user_impersonation"}
	response, err := c.client.AcquireTokenSilent(context.TODO(), scopes)
	if err == nil {
		return response.AccessToken, nil
	}
	// Opens default browser so client can authenticate
	response, err = c.client.AcquireTokenInteractive(context.TODO(), scopes)
	if err != nil {
		return "", nil
	}
	return response.AccessToken, nil
}
