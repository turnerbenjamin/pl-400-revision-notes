package reqbuilder

import "net/http"

type ReqBuilder interface {
	AddQueryParam(queryParam string) ReqBuilder
	Build() (*http.Request, error)
}
