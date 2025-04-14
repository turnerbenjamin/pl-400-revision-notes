package reqbuilder

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type oDataReqBuilder struct {
	httpMethod  string
	path        string
	payload     *bytes.Reader
	queryParams []string
}

func NewODataReqBuilder(httpMethod, path string, payload *bytes.Reader) ReqBuilder {
	return &oDataReqBuilder{
		httpMethod: httpMethod,
		path:       path,
		payload:    payload,
	}
}

func (rb *oDataReqBuilder) AddQueryParam(queryString string) ReqBuilder {
	rb.queryParams = append(rb.queryParams, queryString)
	return rb
}

func (rb *oDataReqBuilder) Build() (*http.Request, error) {
	url := rb.buildUrl()

	var req *http.Request
	var err error
	if rb.payload != nil {
		req, err = http.NewRequest(rb.httpMethod, url, rb.payload)
	} else {
		req, err = http.NewRequest(rb.httpMethod, url, nil)
	}
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (rb *oDataReqBuilder) buildUrl() string {
	if len(rb.queryParams) == 0 {
		return rb.path
	}
	return fmt.Sprintf("%s?%s", rb.path, strings.Join(rb.queryParams, "&"))

}
