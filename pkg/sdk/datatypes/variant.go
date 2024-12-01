package datatypes

// VariantDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-semistructured#variant
// It does not have synonyms. It does not have any attributes.
type VariantDataType struct {
	underlyingType string
}

var VariantDataTypeSynonyms = []string{"VARIANT"}

func parseVariantDataTypeRaw(raw sanitizedDataTypeRaw) (*VariantDataType, error) {
	return &VariantDataType{raw.matchedByType}, nil
}
