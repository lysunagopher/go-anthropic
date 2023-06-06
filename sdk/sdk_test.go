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
	if anthropic, err = NewAnthropic(client, ""); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestNewAnthropic(t *testing.T) {
	art := assert.New(t)
	// Common case.
	anth, err := NewAnthropic(mock.NewHTTPClient(), "")
	if art.NoError(err) {
		art.NotNil(anth)
	}
	// Nil resource.
	anth, err = NewAnthropic(nil, "")
	if art.Error(err) {
		art.Nil(anth)
	}
	// Override root.
	anth, err = NewAnthropic(mock.NewHTTPClient(), "", "http://127.0.0.1")
	if art.NoError(err) {
		art.NotNil(anth)
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
	completion, err = anthropic.Answer("Why is the sky blue?", 255, ModelClaude__V1_0__Instant)
	if art.NoError(err) {
		if art.NotNil(completion) {
			art.NotEmpty(*completion)
		}
	}
}

func TestAnthropic_Do(t *testing.T) {
	art := assert.New(t)
	type test struct {
		response []byte
		status   int
		prompt   string
		model    Model
		tokens   uint32
		expected error
	}
	tests := []test{
		// Common case.
		{
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
		// Internal anthropic error.
		{
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
		// Invalid prompt format.
		{
			response: []byte(`{}`),
			status:   0,
			prompt:   "Why is the sky blue?",
			model:    ModelClaude__V1,
			tokens:   255,
			expected: ErrInvalidPromptFormat,
		},
	}

	for _, tt := range tests {
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
	}
}

func TestAnthropic_ValidatePrompt(t *testing.T) {
	art := assert.New(t)
	type test struct {
		in  string
		out error
	}
	tests := []test{
		{
			// Common case.
			in:  "\n\nHuman: abc123\n\nAssistant:",
			out: nil,
		},
		{
			// Multiline message.
			in:  "\n\nHuman: abc123\nabc123\n\nAssistant:",
			out: nil,
		},
		{
			// Empty message.
			in:  "\n\nHuman: \n\nAssistant:",
			out: nil,
		},
		{
			// Non formatted message.
			in:  "abc123",
			out: ErrInvalidPromptFormat,
		},
		{
			// Empty prompt.
			in:  "",
			out: ErrInvalidPromptFormat,
		},
	}

	for _, tt := range tests {
		out := anthropic.validatePrompt(tt.in)
		art.Equal(tt.out, out)
	}
}
