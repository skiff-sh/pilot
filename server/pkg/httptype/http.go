package httptype

import "net/http"

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
