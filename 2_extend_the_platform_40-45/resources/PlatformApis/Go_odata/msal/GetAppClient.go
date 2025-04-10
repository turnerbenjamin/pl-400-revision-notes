package msal

import (
	"context"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

type appClient struct {
	client      *confidential.Client
	resourceUrl string
}

func GetAppService(c ClientConfig) (DataverseService, error) {

	cred, err := confidential.NewCredFromSecret(c.ClientSecret)
	if err != nil {
		return nil, err
	}

	app, err := confidential.New(c.Authority, c.ClientId, cred)
	return &dataverseService{
		client: &appClient{
			client:      &app,
			resourceUrl: c.ResourceUrl,
		},
		baseurl: c.APIBaseUrl,
	}, err
}

func (c *appClient) AcquireToken() (string, error) {

	scopes := []string{c.resourceUrl + ".default"}
	result, err := c.client.AcquireTokenByCredential(context.TODO(), scopes)
	if err != nil {
		return "", err
	}
	accessToken := result.AccessToken
	return accessToken, nil
}
