// Package service provides functionality for interacting with Microsoft
// Dataverse APIs.
package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/turnerbenjamin/go_odata/msal"
)

// DataverseService defines the interface for interacting with Dataverse API.
// It provides methods to verify connectivity and execute HTTP requests against
// Dataverse endpoints.
type DataverseService interface {
	// TestConnection validates connectivity to Dataverse by attempting to
	// acquire an authentication token. Returns nil if successful, error
	// otherwise.
	TestConnection() error

	// Execute sends the provided HTTP request to Dataverse with proper
	// authentication.
	// Returns a DataverseResponse containing the response data and status
	// information.
	Execute(req *http.Request) (*DataverseResponse, error)
}

// DataverseServiceOptions contains the configuration options for creating a
// DataverseService.
type DataverseServiceOptions struct {
	// Client is the MSAL client used for authentication with Dataverse.
	// It handles token acquisition and caching.
	Client msal.DataverseClient
}

// dataverseService is the internal implementation of the DataverseService
// interface.
type dataverseService struct {
	client     msal.DataverseClient
	httpClient *http.Client
}

// NewDataverseService creates a new DataverseService instance with the provided
// options.
func NewDataverseService(options DataverseServiceOptions) (DataverseService, error) {
	s := dataverseService{
		client: options.Client,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return &s, nil
}

// DataverseResponse encapsulates the HTTP response from a Dataverse API call.
type DataverseResponse struct {
	// StatusCode contains the HTTP status code returned by the API
	StatusCode int
	// Body contains the raw response bytes from the API call
	Body []byte
	// IsSuccessful indicates whether the request was successful (status code
	// 2xx or 3xx)
	IsSuccessful bool
}

// TestConnection validates connectivity to Dataverse by attempting to acquire
// an authentication token from the Microsoft identity platform. It leverages
// token caching if a valid token already exists.
//
// Returns nil if successful, or a wrapped error with details if connection
// fails.
func (s dataverseService) TestConnection() error {
	_, err := s.client.AcquireToken()
	if err != nil {
		return fmt.Errorf("unable to connect: %w", err)
	}
	return nil
}

// Execute sends the provided HTTP request to Dataverse with proper
// authentication.
// It automatically acquires an access token (using the cached token if valid)
// and adds it to the request as a Bearer token. The method also sets
// appropriate headers for JSON content.
//
// The method handles the complete request lifecycle including sending the
// request, reading the response body, and properly closing resources.
//
// Returns a DataverseResponse containing status code, response body, and
// success indicator, or an error if the request fails at any stage.
func (s dataverseService) Execute(req *http.Request) (*DataverseResponse, error) {

	accessToken, err := s.client.AcquireToken()
	if err != nil {
		return nil, err
	}

	req.Header.Set(acceptHeader, contentTypeJSON)
	req.Header.Set(authHeader, bearerTokenPrefix+accessToken)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &DataverseResponse{
		Body:         body,
		StatusCode:   res.StatusCode,
		IsSuccessful: res.StatusCode >= 200 && res.StatusCode < 400,
	}, nil
}
