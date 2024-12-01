package datatypes

// BooleanDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-logical
// It does have synonyms. It does not have any attributes.
type BooleanDataType struct {
	underlyingType string
}

var BooleanDataTypeSynonyms = []string{"BOOLEAN", "BOOL"}

func parseBooleanDataTypeRaw(raw sanitizedDataTypeRaw) (*BooleanDataType, error) {
	return &BooleanDataType{raw.matchedByType}, nil
}
