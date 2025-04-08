package msal

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type dataverseClient interface {
	AcquireToken() (string, error)
}

type DataverseService interface {
	Post(endpoint string, payload []byte) (*DataverseResponse, error)
	Get(endpoint string) (*DataverseResponse, error)
}

type DataverseServiceConfig struct {
	ClientId     string
	TenantId     string
	ClientSecret string
	APIBaseUrl   string
	Authority    string
	ResourceUrl  string
}

type dataverseService struct {
	client  dataverseClient
	baseurl string
}

type DataverseResponse struct {
	StatusCode   int
	Body         *[]byte
	IsSuccessful bool
}

func (s dataverseService) Post(endpoint string, payload []byte) (*DataverseResponse, error) {
	path := s.baseurl + endpoint
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return s.makeRequest(req)
}

func (s dataverseService) Get(endpoint string) (*DataverseResponse, error) {
	path := s.baseurl + endpoint
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return s.makeRequest(req)
}

func (s dataverseService) makeRequest(req *http.Request) (*DataverseResponse, error) {

	accessToken, err := s.client.AcquireToken()
	if err != nil {
		msg := "failed to acquire token"
		return nil, errors.New(msg)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("OData-Version", "4.0")
	req.Header.Set("OData-MaxVersion", "4.0")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Prefer", "return=representation")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &DataverseResponse{
		Body:         &body,
		StatusCode:   res.StatusCode,
		IsSuccessful: res.StatusCode >= 200 && res.StatusCode < 400,
	}, nil
}
