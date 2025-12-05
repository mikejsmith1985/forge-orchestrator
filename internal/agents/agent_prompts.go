package agents

import (
	"fmt"
	"strings"
)

// --- Agent Persona Definitions ---

// SystemPromptArchitect defines the persona and rules for the high-level planner agent.
// This agent's primary job is to decompose human language into a clean, token-optimized plan.
const SystemPromptArchitect = `
You are the Forge Orchestrator Agent. Your role is to convert verbose, high-level user goals
into concise, structured, actionable JSON contracts for specialized worker agents.
Your focus is token efficiency, logical sequencing, and adherence to the project charter.
DO NOT write code. DO NOT engage in conversation. ONLY output the requested JSON contract.
`

// SystemPromptImplementation defines the persona for the developer agent.
// This agent focuses on writing code and the required Playwright test based on the contract.
const SystemPromptImplementation = `
You are a specialized Forge Implementation Agent (Developer). Your primary language is Go and TypeScript/React.
Your sole task is to implement the feature described in the GitHub Issue contract provided by the Orchestrator.
You MUST adhere to the following rules:
1. Adhere to the Self-Documenting Code principle: All code must include verbose, plain-English comments, understandable to a child.
2. For any UI changes, you MUST create a corresponding Playwright test that verifies the UX.
3. Your final output MUST include both the modified code files AND the Playwright test file.
4. DO NOT ask questions or engage in conversation. If you fail, output a FAILURE_REPORT JSON.
`

// SystemPromptTest defines the persona for the QA agent.
// This agent critiques code and executes tests, focusing on UX validation.
const SystemPromptTest = `
You are the Forge Test Agent (QA). Your sole job is to ruthlessly validate code and output.
You MUST prioritize user experience (UX) and functional correctness over API status codes.
Your analysis must focus on whether the React UI renders correctly and handles user interaction.
DO NOT write new code. DO NOT assume success. Critically analyze the provided code and test results.
`

// SystemPromptOptimizer defines the persona for the cost management agent.
// This agent analyzes the Token Ledger data and suggests improvements.
const SystemPromptOptimizer = `
You are the Forge Token Optimizer Agent. Your goal is to reduce operational cost and token waste.
Analyze the provided execution log (Flow ID, model used, input/output tokens, and failure reason, if any).
Your output MUST be a JSON object containing a 'suggestion' field and an 'estimated_savings' field.
Suggestion must be concrete (e.g., 'Change prompt to use JSON instead of Markdown').
`

// GetAgentPrompt retrieves the correct prompt string based on the agent's role.
func GetAgentPrompt(role string) (string, error) {
	resolvedRole := resolveRole(role)

	switch resolvedRole {
	case "Architect":
		return SystemPromptArchitect, nil
	case "Implementation":
		return SystemPromptImplementation, nil
	case "Test":
		return SystemPromptTest, nil
	case "Optimizer":
		return SystemPromptOptimizer, nil
	default:
		return "", fmt.Errorf("unknown agent role: %s", role)
	}
}

// roleAliases maps informal/alternative role names to canonical names
var roleAliases = map[string]string{
	"planner":   "Architect",
	"coder":     "Implementation",
	"developer": "Implementation",
	"dev":       "Implementation",
	"tester":    "Test",
	"qa":        "Test",
	"auditor":   "Optimizer",
	"optimizer": "Optimizer",
}

// canonicalRoles contains the list of valid canonical role names
var canonicalRoles = []string{"Architect", "Implementation", "Test", "Optimizer"}

// resolveRole converts an alias or informal role name to its canonical form.
// It handles case-insensitive matching for both aliases and canonical names.
func resolveRole(role string) string {
	trimmed := strings.TrimSpace(role)
	normalized := strings.ToLower(trimmed)

	// Check if it's an alias
	if canonical, ok := roleAliases[normalized]; ok {
		return canonical
	}

	// Check if it's already a canonical name (case-insensitive)
	for _, name := range canonicalRoles {
		if strings.EqualFold(trimmed, name) {
			return name
		}
	}

	return role // Return as-is, GetAgentPrompt will error
}

// GetCanonicalRoles returns the list of valid canonical role names
func GetCanonicalRoles() []string {
	return canonicalRoles
}

// GetRoleAliases returns a copy of the role aliases map
func GetRoleAliases() map[string]string {
	aliases := make(map[string]string)
	for k, v := range roleAliases {
		aliases[k] = v
	}
	return aliases
}
