package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HTTPClient mocks the simplest http client interface.
type HTTPClient struct {
	body   []byte
	status int
	err    error
}

// NewHTTPClient instantiates an empty mock HTTPClient.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

// Do resolves the request. This is a blocking operation.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	body := ioutil.NopCloser(bytes.NewReader(c.body))
	return &http.Response{
		StatusCode:    c.status,
		Body:          body,
		ContentLength: int64(len(c.body)),
	}, c.err
}

// RespondWith sets the response that is supposed to be sent to callers.
// Persists until overwritten.
func (c *HTTPClient) RespondWith(body []byte, status int, err error) {
	c.body = make([]byte, len(body))
	copy(c.body, body)
	c.status = status
	c.err = err
}
