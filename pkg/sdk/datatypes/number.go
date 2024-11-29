package datatypes

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
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
var NumberDataTypeSubTypes = []string{"INTEGER", "INT", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT"}
var AllNumberDataTypes = append(NumberDataTypeSynonyms, NumberDataTypeSubTypes...)

func parseNumberDataTypeRaw(raw sanitizedDataTypeRaw) (*NumberDataType, error) {
	switch {
	case slices.Contains(NumberDataTypeSynonyms, raw.matchedByType):
		return parseNumberDataTypeWithPrecisionAndScale(raw)
	case slices.Contains(NumberDataTypeSubTypes, raw.matchedByType):
		return parseNumberDataTypeWithoutPrecisionAndScale(raw)
	default:
		return nil, fmt.Errorf("unknown number data type: %s", raw.raw)
	}
}

// parseNumberDataTypeWithPrecisionAndScale extracts precision and scale from the raw number data type input.
// It returns defaults if no arguments were provided. It returns error if any part is not parseable.
func parseNumberDataTypeWithPrecisionAndScale(raw sanitizedDataTypeRaw) (*NumberDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		logging.DebugLogger.Printf("[DEBUG] Returning default number precision and scale")
		return &NumberDataType{DefaultNumberPrecision, DefaultNumberScale}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		logging.DebugLogger.Printf(`number %s could not be parsed, use "NUMBER(precision, scale) format"`, raw.raw)
		return nil, fmt.Errorf(`number %s could not be parsed, use "NUMBER(precision, scale) format"`, raw.raw)
	}
	onlyArgs := r[1 : len(r)-1]
	parts := strings.Split(onlyArgs, ",")
	switch l := len(parts); l {
	case 1:
		precision, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err == nil {
			return &NumberDataType{precision, DefaultNumberScale}, nil
		} else {
			logging.DebugLogger.Printf(`[DEBUG] Could not parse number precision "%s", err: %v`, parts[0], err)
			return nil, fmt.Errorf(`could not parse the number's precision: "%s"`, parts[0])
		}
	case 2:
		precision, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			logging.DebugLogger.Printf(`[DEBUG] Could not parse number precision "%s", err: %v`, parts[0], err)
			return nil, fmt.Errorf(`could not parse the number's precision: "%s"`, parts[0])
		}
		scale, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			logging.DebugLogger.Printf(`[DEBUG] Could not parse number scale "%s", err: %v`, parts[1], err)
			return nil, fmt.Errorf(`could not parse the number's scale: "%s"`, parts[1])
		}
		return &NumberDataType{precision, scale}, nil
	default:
		logging.DebugLogger.Printf("[DEBUG] Unexpected length of number arguments")
		return nil, fmt.Errorf(`number cannot have %d arguments: "%s"; only precision and scale are allowed`, l, onlyArgs)
	}
}

func parseNumberDataTypeWithoutPrecisionAndScale(raw sanitizedDataTypeRaw) (*NumberDataType, error) {
	if raw.raw != raw.matchedByType {
		args := strings.TrimPrefix(raw.raw, raw.matchedByType)
		logging.DebugLogger.Printf("[DEBUG] Number type %s cannot have arguments: %s", raw.matchedByType, args)
		return nil, fmt.Errorf("number type %s cannot have arguments: %s", raw.matchedByType, args)
	} else {
		logging.DebugLogger.Printf("[DEBUG] Returning default number precision and scale")
		return &NumberDataType{DefaultNumberPrecision, DefaultNumberScale}, nil
	}
}
