package datatypes

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
)

// VectorDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-vector#vector
// It does not have synonyms. It does have type (int or float) and dimension required attributes.
type VectorDataType struct {
	innerType      string
	dimension      int
	underlyingType string
}

func (t *VectorDataType) ToSql() string {
	return fmt.Sprintf("%s(%s, %d)", t.underlyingType, t.innerType, t.dimension)
}

// ToLegacyDataTypeSql for vector is the only one correct because in the old implementation it was returned as DataType(dType), so a proper format.
func (t *VectorDataType) ToLegacyDataTypeSql() string {
	return t.ToSql()
}

func (t *VectorDataType) Canonical() string {
	return t.ToSql()
}

var (
	VectorDataTypeSynonyms  = []string{"VECTOR"}
	VectorAllowedInnerTypes = []string{"INT", "FLOAT"}
)

// parseVectorDataTypeRaw extracts type and dimension from the raw vector data type input.
// Both attributes are required so no defaults are returned in case any of them is missing.
func parseVectorDataTypeRaw(raw sanitizedDataTypeRaw) (*VectorDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" || (!strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")")) {
		logging.DebugLogger.Printf(`vector %s could not be parsed, use "%s(type, dimension)" format`, raw.raw, raw.matchedByType)
		return nil, fmt.Errorf(`vector %s could not be parsed, use "%s(type, dimension)" format`, raw.raw, raw.matchedByType)
	}
	onlyArgs := r[1 : len(r)-1]
	parts := strings.Split(onlyArgs, ",")
	switch l := len(parts); l {
	case 2:
		vectorType := strings.TrimSpace(parts[0])
		if !slices.Contains(VectorAllowedInnerTypes, vectorType) {
			logging.DebugLogger.Printf(`[DEBUG] Inner type for vector could not be recognized: "%s"; use one of %s`, parts[0], strings.Join(VectorAllowedInnerTypes, ","))
			return nil, fmt.Errorf(`could not parse vector's inner type': "%s"; use one of %s`, parts[0], strings.Join(VectorAllowedInnerTypes, ","))
		}
		dimension, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			logging.DebugLogger.Printf(`[DEBUG] Could not parse vector's dimension "%s", err: %v`, parts[1], err)
			return nil, fmt.Errorf(`could not parse the vector's dimension: "%s", err: %w`, parts[1], err)
		}
		return &VectorDataType{vectorType, dimension, raw.matchedByType}, nil
	default:
		logging.DebugLogger.Printf("[DEBUG] Unexpected length of vector arguments")
		return nil, fmt.Errorf(`vector cannot have %d arguments: "%s"; use "%s(type, dimension)" format`, l, onlyArgs, raw.matchedByType)
	}
}

func areVectorDataTypesTheSame(a, b *VectorDataType) bool {
	return a.innerType == b.innerType && a.dimension == b.dimension
}
