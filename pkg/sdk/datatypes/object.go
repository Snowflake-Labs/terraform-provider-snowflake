package datatypes

// ObjectDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-semistructured#object
// It does not have synonyms. It does not have any attributes.
type ObjectDataType struct {
	underlyingType string
}

var ObjectDataTypeSynonyms = []string{"OBJECT"}

func parseObjectDataTypeRaw(raw sanitizedDataTypeRaw) (*ObjectDataType, error) {
	return &ObjectDataType{raw.matchedByType}, nil
}
