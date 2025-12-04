package llm

import (
	"testing"
)

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Pure JSON",
			input:   `{"key": "value"}`,
			want:    `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "Markdown Wrapped",
			input:   "Here is the JSON:\n```json\n{\"key\": \"value\"}\n```",
			want:    `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "Preamble and Postscript",
			input:   "Sure, here it is:\n{\"key\": \"value\"}\nHope that helps!",
			want:    `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "Multiple JSON blocks (last one wins)",
			input:   `{"first": 1} ... {"second": 2}`,
			want:    `{"second": 2}`,
			wantErr: false,
		},
		{
			name:    "No JSON",
			input:   "Just some text.",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Nested JSON",
			input:   `{"outer": {"inner": "value"}}`,
			want:    `{"outer": {"inner": "value"}}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTokenCount(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantInput  int
		wantOutput int
	}{
		{
			name:       "Standard Format",
			input:      "Input Tokens: 100\nOutput Tokens: 50",
			wantInput:  100,
			wantOutput: 50,
		},
		{
			name:       "Case Insensitive",
			input:      "input tokens: 10\noutput tokens: 5",
			wantInput:  10,
			wantOutput: 5,
		},
		{
			name:       "Missing Output",
			input:      "Input Tokens: 20",
			wantInput:  20,
			wantOutput: 0,
		},
		{
			name:       "No Tokens",
			input:      "Just text",
			wantInput:  0,
			wantOutput: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInput, gotOutput := ExtractTokenCount(tt.input)
			if gotInput != tt.wantInput {
				t.Errorf("ExtractTokenCount() input = %v, want %v", gotInput, tt.wantInput)
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("ExtractTokenCount() output = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestCleanOutput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "ANSI Colors",
			input: "\x1b[31mError\x1b[0m",
			want:  "Error",
		},
		{
			name:  "No ANSI",
			input: "Clean text",
			want:  "Clean text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanOutput(tt.input)
			if got != tt.want {
				t.Errorf("CleanOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
