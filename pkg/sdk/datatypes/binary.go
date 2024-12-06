package datatypes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
)

const DefaultBinarySize = 8388608

// BinaryDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-text#data-types-for-binary-strings
// It does have synonyms that allow specifying size.
type BinaryDataType struct {
	size           int
	underlyingType string
}

func (t *BinaryDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.size)
}

func (t *BinaryDataType) ToLegacyDataTypeSql() string {
	return BinaryLegacyDataType
}

func (t *BinaryDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", BinaryLegacyDataType, t.size)
}

var BinaryDataTypeSynonyms = []string{BinaryLegacyDataType, "VARBINARY"}

func parseBinaryDataTypeRaw(raw sanitizedDataTypeRaw) (*BinaryDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		logging.DebugLogger.Printf("[DEBUG] Returning default size for binary")
		return &BinaryDataType{DefaultBinarySize, raw.matchedByType}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		logging.DebugLogger.Printf(`binary %s could not be parsed, use "%s(size)" format`, raw.raw, raw.matchedByType)
		return nil, fmt.Errorf(`binary %s could not be parsed, use "%s(size)" format`, raw.raw, raw.matchedByType)
	}
	sizeRaw := r[1 : len(r)-1]
	size, err := strconv.Atoi(strings.TrimSpace(sizeRaw))
	if err != nil {
		logging.DebugLogger.Printf(`[DEBUG] Could not parse binary size "%s", err: %v`, sizeRaw, err)
		return nil, fmt.Errorf(`could not parse the binary's size: "%s", err: %w`, sizeRaw, err)
	}
	return &BinaryDataType{size, raw.matchedByType}, nil
}

func areBinaryDataTypesTheSame(a, b *BinaryDataType) bool {
	return a.size == b.size
}
