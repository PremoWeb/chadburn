package core

import (
	"bytes"
	"text/template"
)

// ContainerInfo holds information about a container that can be used in variable replacements
type ContainerInfo struct {
	Name string
	ID   string
}

// VariableContext holds all the variables that can be used in replacements
type VariableContext struct {
	Container ContainerInfo
}

// ProcessVariables replaces variables in the input string using the provided context
func ProcessVariables(input string, context VariableContext) (string, error) {
	// If the input doesn't contain any template markers, return it as is
	if !bytes.Contains([]byte(input), []byte("{{")) {
		return input, nil
	}

	tmpl, err := template.New("command").Parse(input)
	if err != nil {
		return input, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, context)
	if err != nil {
		return input, err
	}

	return buf.String(), nil
}
