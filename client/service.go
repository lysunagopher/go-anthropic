package client

import (
	"net/http"
)

// Anthropic exposes
type Anthropic struct {
	client HTTPClient
	apiKey string
}

// NewAnthropic instantiates an Anthropic object with provided parameters.
func NewAnthropic(client HTTPClient, apiKey string) (*Anthropic, error) {
	return &Anthropic{
		client: client,
		apiKey: apiKey,
	}, nil
}

// Do performs a request to
func (a *Anthropic) Do(request Request) (Response, error) {
	// TODO: implement
	panic("not implemented")
}

// ValidatePrompt check for prompt format and returns an error on validation failure.
func (_ Anthropic) ValidatePrompt(prompt string) error {
	if !promptRegexp.MatchString(prompt) {
		return ErrInvalidPromptFormat
	}
	return nil
}

// construct generates a http.Request based on provided Request.
func construct(request Request) http.Request {
	// TODO: implement
	panic("not implemented")
}
