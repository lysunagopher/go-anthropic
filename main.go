package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type SamplingParameters struct {
	Prompt            string            `json:"prompt"`
	Temperature       *float64          `json:"temperature,omitempty"`
	MaxTokensToSample int               `json:"max_tokens_to_sample"`
	StopSequences     []string          `json:"stop_sequences"`
	TopK              *int              `json:"top_k,omitempty"`
	TopP              *float64          `json:"top_p,omitempty"`
	Model             string            `json:"model"`
	Tags              map[string]string `json:"tags,omitempty"`
}

type OnOpen func(response *http.Response) error
type OnUpdate func(completion *CompletionResponse) error

const (
	HumanPrompt   = "\n\nHuman:"
	AIPrompt      = "\n\nAssistant:"
	ClientID      = "anthropic-go/0.1.0"
	DefaultAPIURL = "https://api.anthropic.com"
	DoneMessage   = "[DONE]"
)

type Event string

const (
	EventPing Event = "ping"
)

type CompletionResponse struct {
	Completion string `json:"completion"`
	Stop       string `json:"stop"`
	StopReason string `json:"stop_reason"`
	Truncated  bool   `json:"truncated"`
	Exception  string `json:"exception"`
	LogID      string `json:"log_id"`
}

type Client struct {
	apiKey string
	apiURL string
}

func NewClient(apiKey string, apiURL ...string) *Client {
	url := DefaultAPIURL
	if len(apiURL) > 0 {
		url = apiURL[0]
	}
	return &Client{apiKey: apiKey, apiURL: url}
}

func (c *Client) Complete(ctx context.Context, params *SamplingParameters) (*CompletionResponse, error) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/complete", c.apiURL), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client", ClientID)
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sampling error: %d %s", resp.StatusCode, resp.Status)
	}

	var completion CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return nil, err
	}

	return &completion, nil
}

func (c *Client) CompleteStream(ctx context.Context, params *SamplingParameters, onOpen OnOpen, onUpdate OnUpdate) (*CompletionResponse, error) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/complete", c.apiURL), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client", ClientID)
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to open sampling stream, HTTP status code %d: %s", resp.StatusCode, resp.Status)
	}

	if onOpen != nil {
		if err := onOpen(resp); err != nil {
			return nil, err
		}
	}

	dec := json.NewDecoder(resp.Body)
	var completion *CompletionResponse
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if err := dec.Decode(&completion); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if onUpdate != nil {
			if err := onUpdate(completion); err != nil {
				return nil, err
			}
		}

		if completion.StopReason != "" {
			return completion, nil
		}
	}

	return nil, errors.New("unexpected end of stream")
}
