package msal

type ClientOptions struct {
	ClientId     string
	ClientSecret string
	ResourceUrl  string
	Authority    string
}

type DataverseClient interface {
	AcquireToken() (string, error)
}
