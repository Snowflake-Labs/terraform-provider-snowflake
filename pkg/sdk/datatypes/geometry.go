package datatypes

// GeometryDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-geospatial#geometry-data-type
// It does not have synonyms. It does not have any attributes.
type GeometryDataType struct {
	underlyingType string
}

var GeometryDataTypeSynonyms = []string{"GEOMETRY"}

func parseGeometryDataTypeRaw(raw sanitizedDataTypeRaw) (*GeometryDataType, error) {
	return &GeometryDataType{raw.matchedByType}, nil
}
