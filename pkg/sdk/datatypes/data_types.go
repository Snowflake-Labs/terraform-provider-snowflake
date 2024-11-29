package datatypes

import (
	"fmt"
	"slices"
	"strings"
)

const DefaultVarcharLength = 16777216

type DataType interface {
}

type sanitizedDataTypeRaw struct {
	raw           string
	matchedByType string
}

// TODO [this PR]: test
// TODO [this PR]: support all data types
func ParseDataType(raw string) (DataType, error) {
	dataTypeRaw := strings.TrimSpace(strings.ToUpper(raw))

	if idx := slices.IndexFunc(AllNumberDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseNumberDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, AllNumberDataTypes[idx]})
	}
	if idx := slices.Index(FloatDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseFloatDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, FloatDataTypeSynonyms[idx]})
	}
	return nil, fmt.Errorf("invalid data type: %s", raw)
}

// TODO [this PR]: support all data types
type VarcharDataType struct{}
type BinaryDataType struct{}
type BooleanDataType struct{}
type TimestampLTZDataType struct{}
type TimestampTZDataType struct{}
type TimestampNTZDataType struct{}
type TimeDataType struct{}
type VectorDataType struct{}
