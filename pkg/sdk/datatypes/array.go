package datatypes

// ArrayDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-semistructured#array
// It does not have synonyms. It does not have any attributes.
type ArrayDataType struct {
	underlyingType string
}

func (t *ArrayDataType) ToSql() string {
	return t.underlyingType
}

func (t *ArrayDataType) ToLegacyDataTypeSql() string {
	return t.underlyingType
}

var ArrayDataTypeSynonyms = []string{"ARRAY"}

func parseArrayDataTypeRaw(raw sanitizedDataTypeRaw) (*ArrayDataType, error) {
	return &ArrayDataType{raw.matchedByType}, nil
}