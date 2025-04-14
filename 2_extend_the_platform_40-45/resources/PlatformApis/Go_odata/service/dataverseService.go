package service

import (
	"errors"
	"io"
	"net/http"

	"github.com/turnerbenjamin/go_odata/msal"
)

type DataverseService interface {
	Connect() error
	Execute(req *http.Request) (*DataverseResponse, error)
}

type DataverseServiceOptions struct {
	Client  msal.DataverseClient
	BaseUrl string
}

type dataverseService struct {
	client  msal.DataverseClient
	baseurl string
}

func NewDataverseService(options DataverseServiceOptions) (DataverseService, error) {
	s := dataverseService{
		client:  options.Client,
		baseurl: options.BaseUrl,
	}

	err := s.Connect()
	if err != nil {
		return nil, err
	}

	return &s, nil
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

func (s dataverseService) Execute(req *http.Request) (*DataverseResponse, error) {

	accessToken, err := s.client.AcquireToken()
	if err != nil {
		msg := "failed to acquire token"
		return nil, errors.New(msg)
	}

	req.Header.Set("Accept", "application/json")
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
