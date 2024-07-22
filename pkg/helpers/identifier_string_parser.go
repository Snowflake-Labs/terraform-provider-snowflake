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
	return parseIdentifierStringWithOpts(identifier, func(r *csv.Reader) {
		r.Comma = ParameterIDDelimiter
	})
}

func parseIdentifierStringWithOpts(identifier string, opts func(*csv.Reader)) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(identifier))
	if opts != nil {
		opts(reader)
	}
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read identifier: %s, err = %w", identifier, err)
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("incompatible identifier: %s", identifier)
	}
	return lines[0], nil
}
