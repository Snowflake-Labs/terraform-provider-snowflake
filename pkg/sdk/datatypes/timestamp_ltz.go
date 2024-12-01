package datatypes

// TimestampLtzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does not have any attributes.
type TimestampLtzDataType struct {
	underlyingType string
}

var TimestampLtzDataTypeSynonyms = []string{"TIMESTAMP_LTZ", "TIMESTAMPLTZ", "TIMESTAMP WITH LOCAL TIME ZONE"}

func parseTimestampLtzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampLtzDataType, error) {
	return &TimestampLtzDataType{raw.matchedByType}, nil
}
