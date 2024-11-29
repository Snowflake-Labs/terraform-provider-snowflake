package datatypes

import (
	"fmt"
	"slices"
	"strings"
)

const (
	DefaultNumberPrecision = 38
	DefaultNumberScale     = 0
	DefaultVarcharLength   = 16777216
)

// TODO [this PR]: do we need common struct/interface?
type PreciseDataType interface {
}

// NumberDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-numeric#data-types-for-fixed-point-numbers
// It does have synonyms that allow specifying precision and scale; here called synonyms.
// It does have synonyms that does not allow specifying precision and scale; here called subtypes.
type NumberDataType struct {
	precision int
	scale     int
}

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

var NumberDataTypeSynonyms = []string{"NUMBER", "DECIMAL", "DEC", "NUMERIC"}
var NumberDataTypeSubTypes = []string{"INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT"}
var AllNumberDataTypes = append(NumberDataTypeSynonyms, NumberDataTypeSubTypes...)

type sanitizedDataTypeRaw struct {
	raw           string
	matchedByType string
}

// TODO [this PR]: test
// TODO [this PR]: support all data types
func ParsePreciseDataType(raw string) (PreciseDataType, error) {
	dataTypeRaw := strings.TrimSpace(strings.ToUpper(raw))

	if idx := slices.IndexFunc(AllNumberDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseNumberDataType(sanitizedDataTypeRaw{dataTypeRaw, AllNumberDataTypes[idx]})
	}
	return nil, fmt.Errorf("invalid data type: %s", raw)
}

// TODO [this PR]: support all data types
type FloatDataType struct{}
type VarcharDataType struct{}
type BinaryDataType struct{}
type BooleanDataType struct{}
type TimestampLTZDataType struct{}
type TimestampTZDataType struct{}
type TimestampNTZDataType struct{}
type TimeDataType struct{}
type VectorDataType struct{}
