package datatypes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
)

// TimestampTzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does have optional precision attribute.
type TimestampTzDataType struct {
	precision      int
	underlyingType string
}

func (t *TimestampTzDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
}

func (t *TimestampTzDataType) ToLegacyDataTypeSql() string {
	return TimestampTzLegacyDataType
}

func (t *TimestampTzDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", TimestampTzLegacyDataType, t.precision)
}

var TimestampTzDataTypeSynonyms = []string{TimestampTzLegacyDataType, "TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE"}

func parseTimestampTzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampTzDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		logging.DebugLogger.Printf("[DEBUG] Returning default precision for timestamp tz")
		return &TimestampTzDataType{DefaultTimestampPrecision, raw.matchedByType}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		logging.DebugLogger.Printf(`timestamp tz %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
		return nil, fmt.Errorf(`timestamp tz %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
	}
	precisionRaw := r[1 : len(r)-1]
	precision, err := strconv.Atoi(strings.TrimSpace(precisionRaw))
	if err != nil {
		logging.DebugLogger.Printf(`[DEBUG] Could not parse timestamp tz precision "%s", err: %v`, precisionRaw, err)
		return nil, fmt.Errorf(`could not parse the timestamp's precision: "%s", err: %w`, precisionRaw, err)
	}
	return &TimestampTzDataType{precision, raw.matchedByType}, nil
}

func areTimestampTzDataTypesTheSame(a, b *TimestampTzDataType) bool {
	return a.precision == b.precision
}
