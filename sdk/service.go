package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

// Anthropic is an object used by end user to interact with anthropic api.
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

// DoPrompt is a wrapper for Do, which uses default parameters for most fields.
// It is possible to pass in an override model, if none is provided, DefaultModel is used.
func (a *Anthropic) DoPrompt(prompt string, maxTokens uint32, overrideModel ...Model) (*string, error) {
	model := defaultModel
	if len(overrideModel) != 0 {
		model = overrideModel[0]
	}
	response, err := a.Do(Request{
		Prompt:            prompt,
		Model:             model,
		MaxTokensToSample: maxTokens,
	})
	if err != nil {
		return nil, err
	}
	return &response.Completion, nil
}

// Do performs a request to the anthropic api. This is a blocking operation.
func (a *Anthropic) Do(request Request) (*Response, error) {
	values, err := query.Values(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v?%v", apiRoot, values.Encode()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var r ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
			return nil, err
		}
		return nil, errors.Wrap(ErrInternalAnthropic, fmt.Sprintf(
			"code: %v, type: %v, message: %v",
			res.StatusCode, r.Error.Type, r.Error.Message,
		))
	}
	var response Response
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ValidatePrompt check for prompt format and returns an error on validation failure.
func (a *Anthropic) ValidatePrompt(prompt string) error {
	if !promptRegexp.MatchString(prompt) {
		return ErrInvalidPromptFormat
	}
	return nil
}