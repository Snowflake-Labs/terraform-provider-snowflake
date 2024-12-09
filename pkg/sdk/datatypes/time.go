package datatypes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
)

const DefaultTimePrecision = 9

// TimeDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#time
// It does not have synonyms. It does have optional precision attribute.
type TimeDataType struct {
	precision      int
	underlyingType string
}

func (t *TimeDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.precision)
}

func (t *TimeDataType) ToLegacyDataTypeSql() string {
	return TimeLegacyDataType
}

func (t *TimeDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", TimeLegacyDataType, t.precision)
}

var TimeDataTypeSynonyms = []string{TimeLegacyDataType}

func parseTimeDataTypeRaw(raw sanitizedDataTypeRaw) (*TimeDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		logging.DebugLogger.Printf("[DEBUG] Returning default precision for time")
		return &TimeDataType{DefaultTimePrecision, raw.matchedByType}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		logging.DebugLogger.Printf(`time %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
		return nil, fmt.Errorf(`time %s could not be parsed, use "%s(precision)" format`, raw.raw, raw.matchedByType)
	}
	precisionRaw := r[1 : len(r)-1]
	precision, err := strconv.Atoi(strings.TrimSpace(precisionRaw))
	if err != nil {
		logging.DebugLogger.Printf(`[DEBUG] Could not parse time precision "%s", err: %v`, precisionRaw, err)
		return nil, fmt.Errorf(`could not parse the time's precision: "%s", err: %w`, precisionRaw, err)
	}
	return &TimeDataType{precision, raw.matchedByType}, nil
}

func areTimeDataTypesTheSame(a, b *TimeDataType) bool {
	return a.precision == b.precision
}
