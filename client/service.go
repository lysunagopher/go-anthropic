package client

import (
	"net/http"
)

// Anthropic exposes
type Anthropic struct {
	client HTTPClient
}

// NewAnthropic instantiates an Anthropic object with provided parameters.
func NewAnthropic(client HTTPClient, apiKey string) (*Anthropic, error) {
	return &Anthropic{
		client: client,
	}, nil
}

// Do performs a request to
func (a *Anthropic) Do(request Request) Response {
	// TODO: implement
	panic("not implemented")
}

// construct generates a http.Request based on provided Request.
func construct(request Request) http.Request {
	// TODO: implement
	panic("not implemented")
}
