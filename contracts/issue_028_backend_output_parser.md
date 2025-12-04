# Contract: Backend Output Parser (Issue 028)

**Goal**: Implement robust parsing logic to extract JSON and token usage from LLM responses, handling common noise like markdown fences, ANSI codes, and preamble text.

## Scope
1.  **Create `internal/llm/parser.go`**:
    -   `ExtractJSON(output string) (string, error)`:
        -   Locate the *last* valid JSON block in the string.
        -   Handle markdown code fences (e.g., \`\`\`json ... \`\`\`).
        -   Strip any text before or after the JSON block.
    -   `ExtractTokenCount(output string) (int, int)`:
        -   Use regex to find patterns like `Input Tokens: 123` or `Tokens: 456`.
        -   Return (input, output) counts.
    -   `CleanOutput(output string) string`:
        -   Remove ANSI escape codes (colors).
        -   Remove common debug prefixes if necessary.

2.  **Update `internal/llm/gateway.go`**:
    -   Use `parser.ExtractJSON` to sanitize the response before returning it to the handler.
    -   Use `parser.ExtractTokenCount` as a fallback if the provider API doesn't return usage metadata.

## Success Criteria
-   Unit tests in `internal/llm/parser_test.go` covering:
    -   Pure JSON input.
    -   Markdown-wrapped JSON.
    -   JSON with preamble/postscript text.
    -   Broken/Incomplete JSON (should error).
-   Integration with `gateway.go` verified via existing tests.

## Handoff
-   **Signal File**: `handoffs/issue_028_backend_output_parser.json`
-   **Git Branch**: `feat/issue-028-backend-output-parser`
