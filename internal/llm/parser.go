package llm

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ExtractJSON locates the last valid JSON block in the output string.
// It handles markdown code fences and strips preamble/postscript text.
// Educational Comment: LLMs often wrap JSON in markdown blocks or include conversational text.
// Robust parsing requires finding the actual JSON content within the noise.
func ExtractJSON(output string) (string, error) {
	// Find the last closing brace
	end := strings.LastIndex(output, "}")
	if end == -1 {
		return "", errors.New("no JSON object found in output")
	}

	// Iterate backwards to find the matching opening brace
	balance := 0
	for i := end; i >= 0; i-- {
		char := output[i]
		if char == '}' {
			balance++
		} else if char == '{' {
			balance--
		}

		if balance == 0 {
			// Found the matching opening brace
			return output[i : end+1], nil
		}
	}

	return "", errors.New("no matching opening brace found for JSON object")
}

// ExtractTokenCount attempts to extract token usage from the output string
// using regex patterns like "Input Tokens: 123".
// Returns (inputTokens, outputTokens).
func ExtractTokenCount(output string) (int, int) {
	inputRe := regexp.MustCompile(`(?i)Input Tokens:\s*(\d+)`)
	outputRe := regexp.MustCompile(`(?i)Output Tokens:\s*(\d+)`)

	inputMatch := inputRe.FindStringSubmatch(output)
	outputMatch := outputRe.FindStringSubmatch(output)

	var input, out int

	if len(inputMatch) > 1 {
		input, _ = strconv.Atoi(inputMatch[1])
	}

	if len(outputMatch) > 1 {
		out, _ = strconv.Atoi(outputMatch[1])
	}

	return input, out
}

// CleanOutput removes ANSI escape codes and other common noise.
func CleanOutput(output string) string {
	// Regex for ANSI escape codes
	ansiRe := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRe.ReplaceAllString(output, "")
}
