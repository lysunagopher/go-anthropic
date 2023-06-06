package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Anthropic is an object used by end user to interact with anthropic api.
type Anthropic struct {
	client  HTTPClient
	apiRoot string
	apiKey  string
}

// NewAnthropic instantiates an Anthropic object with provided parameters.
// It is possible to pass in an override api root. If none is provided, apiRoot is used.
func NewAnthropic(client HTTPClient, apiKey string, apiRoot ...string) (*Anthropic, error) {
	if client == nil {
		return nil, ErrNilResource
	}
	root := defaultAPIRoot
	if len(apiRoot) != 0 {
		root = apiRoot[0]
	}
	return &Anthropic{
		client:  client,
		apiRoot: root,
		apiKey:  apiKey,
	}, nil
}

// Answer is a wrapper for Do, which uses default parameters for most fields.
// It is possible to pass in an override model. If none is provided, DefaultModel is used.
func (a *Anthropic) Answer(question string, maxTokens uint32, overrideModel ...Model) (*string, error) {
	model := defaultModel
	if len(overrideModel) != 0 {
		model = overrideModel[0]
	}
	response, err := a.Do(Request{
		Prompt:            a.FormatPrompt(question),
		Model:             model,
		MaxTokensToSample: maxTokens,
	})
	if err != nil {
		return nil, err
	}
	return &response.Completion, nil
}

// Do performs a request to the anthropic api. This is a blocking operation.
// if response doesn't indicate success, return it as error instead.
func (a *Anthropic) Do(request Request) (*SuccessResponse, error) {
	if err := a.ValidatePrompt(request.Prompt); err != nil {
		return nil, err
	}
	j, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, a.apiRoot+"/v1/complete", bytes.NewReader(j))
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
		return nil, errors.Wrap(fmt.Errorf(
			"code: %v, type: %v, message: %v",
			res.StatusCode, r.Error.Type, r.Error.Message,
		), ErrInternalAnthropic.Error())
	}
	var response SuccessResponse
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

// FormatPrompt wraps front into required human-assistant format.
func (a *Anthropic) FormatPrompt(prompt string) string {
	return fmt.Sprintf(promptFormat, prompt)
}
