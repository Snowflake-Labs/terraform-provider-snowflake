package datatypes

import (
	"fmt"
	"slices"
)

const (
	DefaultNumberPrecision = 38
	DefaultNumberScale     = 0
)

// NumberDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-numeric#data-types-for-fixed-point-numbers
// It does have synonyms that allow specifying precision and scale; here called synonyms.
// It does have synonyms that does not allow specifying precision and scale; here called subtypes.
type NumberDataType struct {
	precision int
	scale     int
}

var NumberDataTypeSynonyms = []string{"NUMBER", "DECIMAL", "DEC", "NUMERIC"}
var NumberDataTypeSubTypes = []string{"INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT"}
var AllNumberDataTypes = append(NumberDataTypeSynonyms, NumberDataTypeSubTypes...)

func parseNumberDataType(raw sanitizedDataTypeRaw) (*NumberDataType, error) {
	switch {
	case slices.Contains(NumberDataTypeSynonyms, raw.matchedByType):
		// TODO [this PR]: parse precision and scale
		return nil, nil
	case slices.Contains(NumberDataTypeSubTypes, raw.matchedByType):
		// TODO [this PR]: precision and scale are not allowed
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown number data type: %s", raw.raw)
	}
}
