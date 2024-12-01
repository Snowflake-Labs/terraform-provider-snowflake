package datatypes

// TimeDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#time
// It does not have synonyms. It does not have any attributes.
type TimeDataType struct {
	underlyingType string
}

var TimeDataTypeSynonyms = []string{"TIME"}

func parseTimeDataTypeRaw(raw sanitizedDataTypeRaw) (*TimeDataType, error) {
	return &TimeDataType{raw.matchedByType}, nil
}
