package client

import (
	"errors"
	"net/http"
	"regexp"
)

type (
	// HTTPClient is an abstraction for http client present solely for mocking purposes.
	HTTPClient interface {
		// Do resolves the request. This is a blocking operation.
		Do(req *http.Request) (*http.Response, error)
	}

	// Request wraps all the request parameters (required and optional).
	Request struct {
		// Prompt you want Claude to complete.
		Prompt string `json:"prompt"`
		// Model represents a unique claude's version control identifier.
		Model Model `json:"model"`
		// MaxTokensToSample is a maximum number of tokens to generate before stopping.
		MaxTokensToSample int `json:"max_tokens_to_sample"`
		// StopSequences may include additional strings that will cause the model
		// to stop generating. By default, models stop on "\n\nHuman:".
		StopSequences []string `json:"stop_sequences,omitempty"`
		// Stream indicates whether to incrementally stream the response using SSE.
		Stream *bool `json:"stream,omitempty"`
		// Temperature is the amount of randomness injected into the response. Ranges from 0 to 1. Use temp closer
		// to 0 for analytical / multiple choice, and temp closer to 1 for creative and generative tasks.
		Temperature *float64 `json:"temperature,omitempty"`
		// TopK signals to only sample from the top K options for each subsequent token. Used to remove "long tail"
		// low probability responses. Defaults to -1, which disables it.
		// See: https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277.
		TopK *int `json:"top_k,omitempty"`
		// TopP to do nucleus sampling, in which to compute the cumulative distribution over all the options for each
		// subsequent token in decreasing probability order and cut it off once it reaches a particular probability
		// specified by TopP. Defaults to -1, which disables it. Note that you should either alter Temperature
		// or TopP, but not both.
		TopP *float64 `json:"top_p,omitempty"`
		// Metadata is an object describing metadata about the request.
		Metadata RequestMetadata `json:"metadata,omitempty"`
	}

	// Response wraps all the response fields.
	Response struct {
		// Completion is the resulting completion up to and excluding the stop sequences.
		Completion string `json:"completion"`
		// StopReason is the reason sampling stopped.
		StopReason StopReason `json:"stop_reason"`
	}

	// RequestMetadata is an object describing metadata about the request
	RequestMetadata struct {
		// UserID is an uuid, hash value, or other external identifier for the user who is associated with the request.
		// Anthropic may use this id to help detect abuse. Do not include any identifying information such as name,
		// email address, or phone number.
		UserID string `json:"user_id,omitempty"`
	}

	// Model represents a unique claude's version control identifier.
	// e.g.: "claude-v1", "claude-v1-100k", "claude-instant-v1", ...
	Model string

	// StopReason represents the reason sampling stopped. Always one of "stop_sequence" and "max_tokens".
	StopReason string
)

const (
	// ModelClaude__V1 "claude-v1": Our largest model, ideal for a wide range of more complex tasks.
	ModelClaude__V1 Model = "claude-v1"
	// ModelClaude__V1__100k "claude-v1-100k": An enhanced version of claude-v1 with a 100,000 token (roughly 75,000 word) context window. Ideal for summarizing, analyzing, and querying long documents and conversations for nuanced understanding of complex topics and relationships across very long spans of text.
	ModelClaude__V1__100k Model = "claude-v1-100k"
	// ModelClaude__V1__Instant "claude-instant-v1": A smaller model with far lower latency, sampling at roughly 40 words/sec! Its output quality is somewhat lower than the latest claude-v1 model, particularly for complex tasks. However, it is much less expensive and blazing fast. We believe that this model provides more than adequate performance on a range of tasks including text classification, summarization, and lightweight chat applications, as well as search result summarization.
	ModelClaude__V1__Instant Model = "claude-instant-v1"
	// ModelClaude__V1__Instant100k "claude-instant-v1-100k": An enhanced version of claude-instant-v1 with a 100,000 token context window that retains its performance. Well-suited for high throughput use cases needing both speed and additional context, allowing deeper understanding from extended conversations and documents.
	ModelClaude__V1__Instant100k Model = "claude-instant-v1-100k"
)
const (
	// ModelClaude__V1_3 "claude-v1.3": Compared to claude-v1.2, it's more robust against red-team inputs, better at precise instruction-following, better at code, and better and non-English dialogue and writing.
	ModelClaude__V1_3 Model = "claude-v1.3"
	// ModelClaude__V1_3__100k "claude-v1.3-100k": An enhanced version of claude-v1.3 with a 100,000 token (roughly 75,000 word) context window.
	ModelClaude__V1_3__100k Model = "claude-v1.3-100k"
	// ModelClaude__V1_2 "claude-v1.2": An improved version of claude-v1. It is slightly improved at general helpfulness, instruction following, coding, and other tasks. It is also considerably better with non-English languages. This model also has the ability to role play (in harmless ways) more consistently, and it defaults to writing somewhat longer and more thorough responses.
	ModelClaude__V1_2 Model = "claude-v1.2"
	// ModelClaude__V1_0 "claude-v1.0": An earlier version of claude-v1.
	ModelClaude__V1_0 Model = "claude-v1.0"
	// ModelClaude__V1_1__Instant "claude-instant-v1.1": Our latest version of claude-instant-v1. It is better than claude-instant-v1.0 at a wide variety of tasks including writing, coding, and instruction following. It performs better on academic benchmarks, including math, reading comprehension, and coding tests. It is also more robust against red-teaming inputs.
	ModelClaude__V1_1__Instant Model = "claude-instant-v1.1"
	// ModelClaude__V1_1__Instant100k "claude-instant-v1.1-100k": An enhanced version of claude-instant-v1.1 with a 100,000 token context window that retains its lightning fast 40 word/sec performance.
	ModelClaude__V1_1__Instant100k Model = "claude-instant-v1.1-100k"
	// ModelClaude__V1_0__Instant "claude-instant-v1.0": An earlier version of claude-instant-v1.
	ModelClaude__V1_0__Instant Model = "claude-instant-v1.0"
)

const (
	// StopReasonStopSequence "stop_sequence": we reached a stop sequence — either provided by you via the
	// stop_sequences parameter, or a stop sequence built into the model
	StopReasonStopSequence StopReason = "stop_sequence"
	// StopReasonMaxTokens "max_tokens": we exceeded max_tokens_to_sample or the model's maximum.
	StopReasonMaxTokens StopReason = "stop_sequence"
)

var (
	// ErrInvalidPromptFormat indicates that provided prompt doesn't follow required format.
	ErrInvalidPromptFormat = errors.New("invalid prompt: prompts have to be of following format: `\n\nHuman: ${prompt}\n\nAssistant:`")
)

var (
	// promptRegexp is used in prompt validation.
	promptRegexp = regexp.MustCompile(`\n{2}Human: (.|\n)*\n{2}Assistant:`)
)
