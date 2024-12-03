package datatypes

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// TODO [SNOW-1843440]: generalize definitions for different types; generalize the ParseDataType function
// TODO [SNOW-1843440]: generalize implementation in types (i.e. the internal struct implementing ToLegacyDataTypeSql and containing the underlyingType)
// TODO [SNOW-1843440]: consider known/unknown to use Snowflake defaults and allow better handling in terraform resources
// TODO [SNOW-1843440]: replace old DataTypes

// DataType is the common interface that represents all Snowflake datatypes documented in https://docs.snowflake.com/en/sql-reference/intro-summary-data-types.
type DataType interface {
	ToSql() string
	ToLegacyDataTypeSql() string
}

type sanitizedDataTypeRaw struct {
	raw           string
	matchedByType string
}

// ParseDataType is the entry point to get the implementation of the DataType from input raw string.
// TODO [SNOW-1843440]: order currently matters (e.g. HasPrefix(TIME) can match also TIMESTAMP*, make the checks more precise and order-independent)
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
	if idx := slices.IndexFunc(TimestampLtzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampLtzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampLtzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimestampNtzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampNtzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampNtzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimestampTzDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimestampTzDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimestampTzDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(TimeDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTimeDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, TimeDataTypeSynonyms[idx]})
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

func AreTheSame(a DataType, b DataType) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch v := a.(type) {
	case *ArrayDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *BinaryDataType:
		return castSuccessfully(v, b, areBinaryDataTypesTheSame)
	case *BooleanDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *DateDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *FloatDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *GeographyDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *GeometryDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *NumberDataType:
		return castSuccessfully(v, b, areNumberDataTypesTheSame)
	case *ObjectDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *TextDataType:
		return castSuccessfully(v, b, areTextDataTypesTheSame)
	case *TimeDataType:
		return castSuccessfully(v, b, areTimeDataTypesTheSame)
	case *TimestampLtzDataType:
		return castSuccessfully(v, b, areTimestampLtzDataTypesTheSame)
	case *TimestampNtzDataType:
		return castSuccessfully(v, b, areTimestampNtzDataTypesTheSame)
	case *TimestampTzDataType:
		return castSuccessfully(v, b, areTimestampTzDataTypesTheSame)
	case *VariantDataType:
		return castSuccessfully(v, b, noArgsDataTypesAreTheSame)
	case *VectorDataType:
		return castSuccessfully(v, b, areVectorDataTypesTheSame)
	}
	return false
}

func IsTextDataType(a DataType) bool {
	_, ok := a.(*TextDataType)
	return ok
}

func castSuccessfully[T any](a T, b DataType, invoke func(a T, b T) bool) bool {
	if dCasted, ok := b.(T); ok {
		return invoke(a, dCasted)
	}
	return false
}

func noArgsDataTypesAreTheSame[T DataType](_ T, _ T) bool {
	return true
}
