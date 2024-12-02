package datatypes

import (
	"fmt"
	"slices"
	"strings"
)

type DataType interface {
	ToSql() string
	ToLegacyDataTypeSql() string
}

type sanitizedDataTypeRaw struct {
	raw           string
	matchedByType string
}

// https://docs.snowflake.com/en/sql-reference/intro-summary-data-types
// Session-configurable TIMESTAMP alias is currenlty not supported (https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp).
func ParseDataType(raw string) (DataType, error) {
	dataTypeRaw := strings.TrimSpace(strings.ToUpper(raw))

	if idx := slices.IndexFunc(AllNumberDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseNumberDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, AllNumberDataTypes[idx]})
	}
	if idx := slices.Index(FloatDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseFloatDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, FloatDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(AllTextDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTextDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, AllTextDataTypes[idx]})
	}
	if idx := slices.IndexFunc(BinaryDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseBinaryDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, BinaryDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(BooleanDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseBooleanDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, BooleanDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(FloatDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseFloatDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, FloatDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(DateDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseDateDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, DateDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(TimeDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseTimeDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimeDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimestampLtzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampLtzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampLtzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimestampNtzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampNtzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampNtzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimestampTzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampTzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampTzDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(VariantDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseVariantDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, VariantDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(ObjectDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseObjectDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, ObjectDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(ArrayDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseArrayDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, ArrayDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(GeographyDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseGeographyDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, GeographyDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(GeometryDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseGeometryDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, GeometryDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(VectorDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseVectorDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, VectorDataTypeSynonyms[idx]})
	}

	return nil, fmt.Errorf("invalid data type: %s", raw)
}