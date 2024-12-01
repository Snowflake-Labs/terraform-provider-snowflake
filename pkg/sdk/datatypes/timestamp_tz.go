package datatypes

// TimestampTzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does not have any attributes.
type TimestampTzDataType struct {
	underlyingType string
}

var TimestampTzDataTypeSynonyms = []string{"TIMESTAMP_TZ", "TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE"}

func parseTimestampTzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampTzDataType, error) {
	return &TimestampTzDataType{raw.matchedByType}, nil
}
