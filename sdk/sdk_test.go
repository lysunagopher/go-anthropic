package sdk

import (
	"net/http"
	"os"
	"testing"

	"github.com/staropshq/go-anthropic/helpers/mock"

	"github.com/stretchr/testify/assert"
)

var (
	client    *mock.HTTPClient
	anthropic *Anthropic
)

func TestMain(m *testing.M) {
	client = mock.NewHTTPClient()
	var err error
	if anthropic, err = NewAnthropic(client, "", ""); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestNewAnthropic(t *testing.T) {
	art := assert.New(t)
	// Common case.
	anth, err := NewAnthropic(mock.NewHTTPClient(), "", "")
	if art.NoError(err) {
		art.NotNil(anth)
	}
	// Nil resource.
	anth, err = NewAnthropic(nil, "", "")
	if art.Error(err) {
		art.Nil(anth)
	}
}

func TestAnthropic_DoPrompt(t *testing.T) {
	art := assert.New(t)
	// Common case.
	response := []byte(`{
		"completion":" The sky appears blue to our eyes due to the way the atmosphere interacts with sunlight.",
		"stop_reason":"stop_sequence"
	}`)
	client.RespondWith(response, http.StatusOK, nil)
	completion, err := anthropic.Answer("Why is the sky blue?", 255)
	if art.NoError(err) {
		if art.NotNil(completion) {
			art.NotEmpty(*completion)
		}
	}
	// Override model.
	completion, err = anthropic.Answer("Why is the sky blue?", 255)
	if art.NoError(err) {
		if art.NotNil(completion) {
			art.NotEmpty(*completion)
		}
	}
}

func TestAnthropic_Do(t *testing.T) {
	type test struct {
		name     string
		response []byte
		status   int
		prompt   string
		model    Model
		tokens   uint32
		expected error
	}
	tests := []test{
		{
			name: "common case",
			response: []byte(`{
				"completion":" The sky appears blue to our eyes due to the way the atmosphere interacts with sunlight.",
				"stop_reason":"stop_sequence"
			}`),
			status:   http.StatusOK,
			prompt:   anthropic.formatPrompt("Why is the sky blue?"),
			model:    ModelClaude__V1,
			tokens:   255,
			expected: nil,
		},
		{
			name: "internal anthropic error",
			response: []byte(`{
				"type": "invalid_request_error",
				"message":"field required"
			}`),
			status:   http.StatusBadRequest,
			prompt:   anthropic.formatPrompt("Why is the sky blue?"),
			model:    ModelClaude__V1,
			tokens:   0,
			expected: ErrInternalAnthropic,
		},
		{
			name:     "internal prompt format",
			response: []byte(`{}`),
			status:   0,
			prompt:   "Why is the sky blue?",
			model:    ModelClaude__V1,
			tokens:   255,
			expected: ErrInvalidPromptFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_t *testing.T) {
			art := assert.New(_t)
			client.RespondWith(tt.response, tt.status, nil)
			resp, err := anthropic.Do(Request{
				Prompt:            tt.prompt,
				Model:             tt.model,
				MaxTokensToSample: tt.tokens,
			})
			if tt.expected != nil {
				if art.Error(err) {
					art.ErrorContains(err, tt.expected.Error())
					art.Nil(resp)
				}
			} else {
				if art.NoError(err) {
					art.NotNil(resp)
				}
			}
		})
	}
}

func TestAnthropic_ValidatePrompt(t *testing.T) {
	type test struct {
		name string
		in   string
		out  error
	}
	tests := []test{
		{
			name: "common case",
			in:   "\n\nHuman: abc123\n\nAssistant:",
			out:  nil,
		},
		{
			name: "multiline message",
			in:   "\n\nHuman: abc123\nabc123\n\nAssistant:",
			out:  nil,
		},
		{
			name: "empty message",
			in:   "\n\nHuman: \n\nAssistant:",
			out:  nil,
		},
		{
			name: "non formatted message",
			in:   "abc123",
			out:  ErrInvalidPromptFormat,
		},
		{
			name: "empty prompt",
			in:   "",
			out:  ErrInvalidPromptFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_t *testing.T) {
			art := assert.New(_t)
			out := anthropic.validatePrompt(tt.in)
			art.Equal(tt.out, out)
		})
	}
}
