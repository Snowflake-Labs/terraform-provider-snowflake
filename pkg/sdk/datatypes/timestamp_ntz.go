package datatypes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
)

// TimestampNtzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does have optional precision attribute.
type TimestampNtzDataType struct {
	precision      int
	underlyingType string
}

func (t *TimestampNtzDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
}

func (t *TimestampNtzDataType) ToLegacyDataTypeSql() string {
	return TimestampNtzLegacyDataType
}

func (t *TimestampNtzDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", TimestampNtzLegacyDataType, t.precision)
}

var TimestampNtzDataTypeSynonyms = []string{TimestampNtzLegacyDataType, "TIMESTAMPNTZ", "TIMESTAMP WITHOUT TIME ZONE", "DATETIME"}

func parseTimestampNtzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampNtzDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		logging.DebugLogger.Printf("[DEBUG] Returning default precision for timestamp ntz")
		return &TimestampNtzDataType{DefaultTimestampPrecision, raw.matchedByType}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		logging.DebugLogger.Printf(`timestamp ntz %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
		return nil, fmt.Errorf(`timestamp ntz %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
	}
	precisionRaw := r[1 : len(r)-1]
	precision, err := strconv.Atoi(strings.TrimSpace(precisionRaw))
	if err != nil {
		logging.DebugLogger.Printf(`[DEBUG] Could not parse timestamp ntz precision "%s", err: %v`, precisionRaw, err)
		return nil, fmt.Errorf(`could not parse the timestamp's precision: "%s", err: %w`, precisionRaw, err)
	}
	return &TimestampNtzDataType{precision, raw.matchedByType}, nil
}

func areTimestampNtzDataTypesTheSame(a, b *TimestampNtzDataType) bool {
	return a.precision == b.precision
}
