# Issue #040: Improve Token Estimation Accuracy

**Priority:** 🟢 LOW  
**Estimated Tokens:** ~1,500 (Medium complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-012 from v1-analysis.md

Current token estimation is naive:
```go
// Simple approximation: len(text) / 4
count := len(req.Text) / 4
```

This can be 20-50% off from actual token counts because:
- Code has different tokenization than prose
- Non-ASCII characters may use more tokens
- OpenAI and Anthropic use different tokenizers

---

## 2. 📋 Acceptance Criteria

### Backend (Go) - Option A: Improved Heuristic
- [ ] Implement word-based estimation (more accurate than character-based)
- [ ] Account for code patterns (brackets, operators = usually single tokens)
- [ ] Handle non-ASCII characters (assume 2-4 tokens each)
- [ ] Add provider parameter to endpoint for provider-specific adjustments

### Backend (Go) - Option B: Real Tokenizer (Preferred)
- [ ] Add `tiktoken-go` library for accurate OpenAI tokenization
- [ ] Add simple Anthropic estimation (similar to GPT-4 tokenizer)
- [ ] Return both estimated and actual counts in API response

### API Enhancement
- [ ] Update `POST /api/tokens/estimate` to accept provider parameter
- [ ] Return breakdown: `{ "estimate": 150, "method": "tiktoken", "provider": "openai" }`

### Frontend (React)
- [ ] Update TokenMeter to optionally show provider-specific counts
- [ ] Add provider selector if multiple providers are configured

### Tests
- [ ] Benchmark test: Compare estimated vs actual tokens (after execution)
- [ ] Unit test: Verify different inputs produce reasonable estimates

---

## 3. 📊 Token Efficiency Strategy

- If using tiktoken-go, it's a single dependency addition
- Alternatively, improved heuristic is ~40 lines of code
- Minimal frontend changes

---

## 4. 🏗️ Technical Specification

### Option A: Improved Heuristic
```go
func EstimateTokens(text string, provider string) int {
    // Base: word count + punctuation
    words := strings.Fields(text)
    wordCount := len(words)
    
    // Code adjustments
    codePatterns := regexp.MustCompile(`[{}\[\]();:,<>=+\-*/]`)
    codeTokens := len(codePatterns.FindAllString(text, -1))
    
    // Non-ASCII penalty
    nonAscii := 0
    for _, r := range text {
        if r > 127 {
            nonAscii++
        }
    }
    
    // Provider adjustments
    multiplier := 1.0
    if provider == "Anthropic" {
        multiplier = 1.1 // Claude tends to tokenize slightly higher
    }
    
    estimate := float64(wordCount) * 1.3 + float64(codeTokens) * 0.5 + float64(nonAscii) * 2
    return int(estimate * multiplier)
}
```

### Option B: Real Tokenizer
```go
import "github.com/pkoukk/tiktoken-go"

func EstimateTokensAccurate(text string, model string) (int, error) {
    encoding, err := tiktoken.EncodingForModel(model)
    if err != nil {
        // Fallback to heuristic
        return EstimateTokens(text, "openai"), nil
    }
    tokens := encoding.Encode(text, nil, nil)
    return len(tokens), nil
}
```

### Updated API Response
```json
{
    "count": 156,
    "method": "tiktoken",
    "model": "gpt-4",
    "breakdown": {
        "words": 42,
        "punctuation": 15,
        "special": 3
    }
}
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `internal/server/ledger.go` (update handler) |
| CREATE | `internal/tokenizer/estimator.go` (new logic) |
| MODIFY | `go.mod` (add tiktoken-go if using Option B) |
| MODIFY | `frontend/src/components/Architect/TokenMeter.tsx` |
| MODIFY | `frontend/src/components/Architect/ArchitectView.tsx` |

---

## 6. ✅ Definition of Done

1. Token estimation within 15% of actual for typical prompts
2. API returns estimation method used
3. Provider parameter is optional (defaults to OpenAI)
4. Frontend displays estimation in TokenMeter
5. Benchmark test documents accuracy improvement
