package core

import (
	"testing"
)

func TestProcessVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		context  VariableContext
		expected string
		hasError bool
	}{
		{
			name:     "No variables",
			input:    "echo hello world",
			context:  VariableContext{},
			expected: "echo hello world",
			hasError: false,
		},
		{
			name:  "Container name variable",
			input: "docker restart {{.Container.Name}}",
			context: VariableContext{
				Container: ContainerInfo{
					Name: "test-container",
					ID:   "abc123",
				},
			},
			expected: "docker restart test-container",
			hasError: false,
		},
		{
			name:  "Container ID variable",
			input: "docker restart {{.Container.ID}}",
			context: VariableContext{
				Container: ContainerInfo{
					Name: "test-container",
					ID:   "abc123",
				},
			},
			expected: "docker restart abc123",
			hasError: false,
		},
		{
			name:  "Multiple variables",
			input: "echo Container {{.Container.Name}} has ID {{.Container.ID}}",
			context: VariableContext{
				Container: ContainerInfo{
					Name: "test-container",
					ID:   "abc123",
				},
			},
			expected: "echo Container test-container has ID abc123",
			hasError: false,
		},
		{
			name:     "Invalid template",
			input:    "echo {{.Invalid}",
			context:  VariableContext{},
			expected: "echo {{.Invalid}",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ProcessVariables(tt.input, tt.context)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.hasError && result != tt.expected {
				t.Errorf("Expected %q but got %q", tt.expected, result)
			}
		})
	}
}
