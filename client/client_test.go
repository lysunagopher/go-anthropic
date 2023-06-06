package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnthropic_Do(t *testing.T) {

}

func TestAnthropic_ValidatePrompt(t *testing.T) {
	art := assert.New(t)
	tests := []struct {
		in  string
		out error
	}{
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

	anthropic := Anthropic{}
	for _, tt := range tests {
		out := anthropic.ValidatePrompt(tt.in)
		art.Equal(tt.out, out)
	}
}
