package helpers

import (
	"encoding/csv"
	"fmt"
	"strings"
)

const (
	ParameterIDDelimiter = '.'
)

func ParseIdentifierString(identifier string) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(identifier))
	reader.Comma = ParameterIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read identifier: %s, err = %w", identifier, err)
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("incompatible identifier: %s", identifier)
	}
	return lines[0], nil
}
