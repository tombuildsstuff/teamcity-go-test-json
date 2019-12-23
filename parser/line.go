package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseLine(input string) (*GoTestJsonLogLine, error) {
	// e.g. "go: downloading github.com/hashicorp/azure-sdk-for-go [..]"
	if strings.HasPrefix(input, "go:") {
		return nil, nil
	}

	var output GoTestJsonLogLine

	if err := json.Unmarshal([]byte(input), &output); err != nil {
		return nil, fmt.Errorf("Error deserializing: %q", input)
	}

	return &output, nil
}
