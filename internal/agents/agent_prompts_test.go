package agents

import (
	"strings"
	"testing"
)

func TestGetAgentPrompt_CanonicalNames(t *testing.T) {
	tests := []struct {
		role     string
		expected string
	}{
		{"Architect", SystemPromptArchitect},
		{"Implementation", SystemPromptImplementation},
		{"Test", SystemPromptTest},
		{"Optimizer", SystemPromptOptimizer},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			prompt, err := GetAgentPrompt(tt.role)
			if err != nil {
				t.Fatalf("GetAgentPrompt(%q) returned error: %v", tt.role, err)
			}
			if prompt != tt.expected {
				t.Errorf("GetAgentPrompt(%q) = %q, want %q", tt.role, prompt[:50], tt.expected[:50])
			}
		})
	}
}

func TestGetAgentPrompt_Aliases(t *testing.T) {
	tests := []struct {
		alias    string
		expected string
	}{
		{"coder", SystemPromptImplementation},
		{"planner", SystemPromptArchitect},
		{"developer", SystemPromptImplementation},
		{"dev", SystemPromptImplementation},
		{"tester", SystemPromptTest},
		{"qa", SystemPromptTest},
		{"auditor", SystemPromptOptimizer},
		{"optimizer", SystemPromptOptimizer},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			prompt, err := GetAgentPrompt(tt.alias)
			if err != nil {
				t.Fatalf("GetAgentPrompt(%q) returned error: %v", tt.alias, err)
			}
			if prompt != tt.expected {
				t.Errorf("GetAgentPrompt(%q) returned wrong prompt", tt.alias)
			}
		})
	}
}

func TestGetAgentPrompt_CaseInsensitive(t *testing.T) {
	tests := []struct {
		role     string
		expected string
	}{
		// Canonical names with different cases
		{"ARCHITECT", SystemPromptArchitect},
		{"architect", SystemPromptArchitect},
		{"ArChItEcT", SystemPromptArchitect},
		{"IMPLEMENTATION", SystemPromptImplementation},
		{"implementation", SystemPromptImplementation},
		{"TEST", SystemPromptTest},
		{"test", SystemPromptTest},
		{"OPTIMIZER", SystemPromptOptimizer},
		// Aliases with different cases
		{"PLANNER", SystemPromptArchitect},
		{"Planner", SystemPromptArchitect},
		{"CODER", SystemPromptImplementation},
		{"Coder", SystemPromptImplementation},
		{"TESTER", SystemPromptTest},
		{"Tester", SystemPromptTest},
		{"QA", SystemPromptTest},
		{"Qa", SystemPromptTest},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			prompt, err := GetAgentPrompt(tt.role)
			if err != nil {
				t.Fatalf("GetAgentPrompt(%q) returned error: %v", tt.role, err)
			}
			if prompt != tt.expected {
				t.Errorf("GetAgentPrompt(%q) returned wrong prompt", tt.role)
			}
		})
	}
}

func TestGetAgentPrompt_UnknownRole(t *testing.T) {
	unknownRoles := []string{
		"unknown",
		"invalid",
		"random",
		"",
		"  ",
	}

	for _, role := range unknownRoles {
		t.Run(role, func(t *testing.T) {
			_, err := GetAgentPrompt(role)
			if err == nil {
				t.Errorf("GetAgentPrompt(%q) should have returned an error", role)
			}
			if !strings.Contains(err.Error(), "unknown agent role") {
				t.Errorf("GetAgentPrompt(%q) error = %v, want 'unknown agent role'", role, err)
			}
		})
	}
}

func TestGetAgentPrompt_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		role     string
		expected string
	}{
		{"  coder  ", SystemPromptImplementation},
		{"\tplanner\t", SystemPromptArchitect},
		{" Architect ", SystemPromptArchitect},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			prompt, err := GetAgentPrompt(tt.role)
			if err != nil {
				t.Fatalf("GetAgentPrompt(%q) returned error: %v", tt.role, err)
			}
			if prompt != tt.expected {
				t.Errorf("GetAgentPrompt(%q) returned wrong prompt", tt.role)
			}
		})
	}
}

func TestResolveRole(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"coder", "Implementation"},
		{"CODER", "Implementation"},
		{"planner", "Architect"},
		{"Architect", "Architect"},
		{"unknown", "unknown"}, // Returns as-is for unknown roles
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := resolveRole(tt.input)
			if result != tt.expected {
				t.Errorf("resolveRole(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCanonicalRoles(t *testing.T) {
	roles := GetCanonicalRoles()
	if len(roles) != 4 {
		t.Errorf("GetCanonicalRoles() returned %d roles, want 4", len(roles))
	}

	expected := map[string]bool{
		"Architect":      true,
		"Implementation": true,
		"Test":           true,
		"Optimizer":      true,
	}

	for _, role := range roles {
		if !expected[role] {
			t.Errorf("Unexpected role in GetCanonicalRoles(): %s", role)
		}
	}
}

func TestGetRoleAliases(t *testing.T) {
	aliases := GetRoleAliases()

	expectedAliases := map[string]string{
		"planner":   "Architect",
		"coder":     "Implementation",
		"developer": "Implementation",
		"dev":       "Implementation",
		"tester":    "Test",
		"qa":        "Test",
		"auditor":   "Optimizer",
		"optimizer": "Optimizer",
	}

	if len(aliases) != len(expectedAliases) {
		t.Errorf("GetRoleAliases() returned %d aliases, want %d", len(aliases), len(expectedAliases))
	}

	for alias, canonical := range expectedAliases {
		if aliases[alias] != canonical {
			t.Errorf("GetRoleAliases()[%q] = %q, want %q", alias, aliases[alias], canonical)
		}
	}
}
