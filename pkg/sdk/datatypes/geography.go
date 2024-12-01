package datatypes

// GeographyDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-geospatial#geography-data-type
// It does not have synonyms. It does not have any attributes.
type GeographyDataType struct {
	underlyingType string
}

var GeographyDataTypeSynonyms = []string{"GEOGRAPHY"}

func parseGeographyDataTypeRaw(raw sanitizedDataTypeRaw) (*GeographyDataType, error) {
	return &GeographyDataType{raw.matchedByType}, nil
}
