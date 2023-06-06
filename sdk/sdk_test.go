package sdk_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/staropshq/go-anthropic/helpers/mock"
	"github.com/staropshq/go-anthropic/sdk"

	"github.com/stretchr/testify/assert"
)

var (
	client    *mock.HTTPClient
	anthropic *sdk.Anthropic
)

func TestMain(m *testing.M) {
	client = mock.NewHTTPClient()
	var err error
	if anthropic, err = sdk.NewAnthropic(client, ""); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestNewAnthropic(t *testing.T) {
	art := assert.New(t)
	// Common case.
	anth, err := sdk.NewAnthropic(mock.NewHTTPClient(), "")
	if art.NoError(err) {
		art.NotNil(anth)
	}
	// Nil resource.
	anth, err = sdk.NewAnthropic(nil, "")
	if art.Error(err) {
		art.Nil(anth)
	}
}

func TestAnthropic_DoPrompt(t *testing.T) {
	art := assert.New(t)
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
}

func TestAnthropic_Do(t *testing.T) {
	art := assert.New(t)
	response := []byte(`{
		"completion":" The sky appears blue to our eyes due to the way the atmosphere interacts with sunlight.",
		"stop_reason":"stop_sequence"
	}`)
	client.RespondWith(response, http.StatusOK, nil)
	completion, err := anthropic.Do(sdk.Request{
		Prompt:            "Why is the sky blue?",
		Model:             sdk.ModelClaude__V1,
		MaxTokensToSample: 255,
	})
	if art.NoError(err) {
		if art.NotNil(completion) {
			art.NotEmpty(*completion)
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
			out: sdk.ErrInvalidPromptFormat,
		},
		{
			// Empty prompt.
			in:  "",
			out: sdk.ErrInvalidPromptFormat,
		},
	}

	for _, tt := range tests {
		out := anthropic.ValidatePrompt(tt.in)
		art.Equal(tt.out, out)
	}
}
