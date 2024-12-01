package datatypes

// TimestampNtzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does not have any attributes.
type TimestampNtzDataType struct {
	underlyingType string
}

var TimestampNtzDataTypeSynonyms = []string{"TIMESTAMP_NTZ", "TIMESTAMPNTZ", "TIMESTAMP WITHOUT TIME ZONE", "DATETIME"}

func parseTimestampNtzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampNtzDataType, error) {
	return &TimestampNtzDataType{raw.matchedByType}, nil
}
