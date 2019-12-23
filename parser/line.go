package parser

import (
	"encoding/json"
	"fmt"
)

func ParseLine(input string) (*GoTestJsonLogLine, error) {
	var output GoTestJsonLogLine

	if err := json.Unmarshal([]byte(input), &output); err != nil {
		return nil, fmt.Errorf("Error deserializing: %q", input)
	}

	return &output, nil
}
