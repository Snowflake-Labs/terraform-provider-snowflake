package datatypes

// DateDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#date
// It does not have synonyms. It does not have any attributes.
type DateDataType struct {
	underlyingType string
}

var DateDataTypeSynonyms = []string{"DATE"}

func parseDateDataTypeRaw(raw sanitizedDataTypeRaw) (*DateDataType, error) {
	return &DateDataType{raw.matchedByType}, nil
}
