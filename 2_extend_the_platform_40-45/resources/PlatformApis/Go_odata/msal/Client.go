package msal

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type ClientConfig struct {
	ClientId     string
	ClientSecret string
	ResourceUrl  string
	APIBaseUrl   string
	Authority    string
}

type dataverseClient interface {
	AcquireToken() (string, error)
}

type DataverseService interface {
	Post(endpoint string, payload []byte) (*DataverseResponse, error)
	Get(endpoint string) (*DataverseResponse, error)
	GetMany(endpoint string) (*DataverseResponse, error)
	Connect() error
	GetNext(url string) (*DataverseResponse, error)
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

func (s dataverseService) Connect() error {
	_, err := s.client.AcquireToken()
	return err
}

func (s dataverseService) Post(endpoint string, payload []byte) (*DataverseResponse, error) {
	path := s.baseurl + endpoint
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")
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

func (s dataverseService) GetMany(endpoint string) (*DataverseResponse, error) {
	path := s.baseurl + endpoint
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Prefer", "odata.maxpagesize=5")
	return s.makeRequest(req)
}

func (s dataverseService) GetNext(url string) (*DataverseResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Prefer", "odata.maxpagesize=2")
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
