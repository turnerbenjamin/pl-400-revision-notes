// Package requestbuilder provides tools for building HTTP requests,
// specifically designed to work with OData APIs.
package requestbuilder

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// RequestBuilder defines an interface for constructing HTTP requests.
// It uses the builder pattern to allow for chained method calls.
type RequestBuilder interface {
	// AddQueryParam adds a query parameter to the request.
	// Returns the builder for method chaining.
	AddQueryParam(key, value string) RequestBuilder

	// AddHeader adds an HTTP header to the request.
	// Returns the builder for method chaining.
	AddHeader(key, value string) RequestBuilder

	// Build constructs and returns the final http.Request object.
	// Returns an error if the request cannot be created.
	Build() (*http.Request, error)
}

// requestBuilder implements the RequestBuilder interface.
type requestBuilder struct {
	httpMethod  string
	path        string
	payload     io.Reader
	queryParams url.Values
	headers     http.Header
}

// NewRequestBuilder creates a new RequestBuilder instance with the specified
// HTTP method, path, and payload. It initializes empty collections for query
// parameters and headers.
func NewRequestBuilder(httpMethod, path string, payload io.Reader) RequestBuilder {
	return &requestBuilder{
		httpMethod:  httpMethod,
		path:        path,
		payload:     payload,
		queryParams: make(url.Values),
		headers:     make(http.Header),
	}
}

// AddQueryParam adds a query parameter with the specified key and value to the
// request.
// Returns the builder for method chaining.
func (rb *requestBuilder) AddQueryParam(key, value string) RequestBuilder {
	rb.queryParams.Add(key, value)
	return rb
}

// AddHeader adds an HTTP header with the specified key and value to the request.
// Returns the builder for method chaining.
func (rb *requestBuilder) AddHeader(key, value string) RequestBuilder {
	rb.headers.Add(key, value)
	return rb
}

// Build constructs and returns the final http.Request object using the
// configured parameters, headers, and payload. It returns an error if the
// request cannot be created.
func (rb *requestBuilder) Build() (*http.Request, error) {
	url := rb.buildURL()

	req, err := http.NewRequest(rb.httpMethod, url, rb.payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, values := range rb.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return req, nil
}

// buildURL constructs the complete URL from the path and query parameters.
// If there are no query parameters, it returns just the path.
func (rb *requestBuilder) buildURL() string {
	if len(rb.queryParams) == 0 {
		return rb.path
	}
	return fmt.Sprintf("%s?%s", rb.path, rb.queryParams.Encode())
}
