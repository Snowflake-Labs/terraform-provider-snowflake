package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataTypeFromString(t *testing.T) {
	type test struct {
		input string
		want  DataType
	}

	tests := []test{
		// case insensitive.
		{input: "STRING", want: DataTypeVARCHAR},
		{input: "string", want: DataTypeVARCHAR},
		{input: "String", want: DataTypeVARCHAR},

		// number types.
		{input: "number", want: DataTypeNumber},
		{input: "decimal", want: DataTypeNumber},
		{input: "numeric", want: DataTypeNumber},
		{input: "int", want: DataTypeNumber},
		{input: "integer", want: DataTypeNumber},
		{input: "bigint", want: DataTypeNumber},
		{input: "smallint", want: DataTypeNumber},
		{input: "tinyint", want: DataTypeNumber},
		{input: "byteint", want: DataTypeNumber},

		// float types.
		{input: "float", want: DataTypeFloat},
		{input: "float4", want: DataTypeFloat},
		{input: "float8", want: DataTypeFloat},
		{input: "double", want: DataTypeFloat},
		{input: "double precision", want: DataTypeFloat},
		{input: "real", want: DataTypeFloat},

		// varchar types.
		{input: "varchar", want: DataTypeVARCHAR},
		{input: "char", want: DataTypeVARCHAR},
		{input: "character", want: DataTypeVARCHAR},
		{input: "string", want: DataTypeVARCHAR},
		{input: "text", want: DataTypeVARCHAR},

		// binary types.
		{input: "binary", want: DataTypeBinary},
		{input: "varbinary", want: DataTypeBinary},
		{input: "boolean", want: DataTypeBoolean},

		// boolean types.
		{input: "boolean", want: DataTypeBoolean},
		{input: "bool", want: DataTypeBoolean},

		// timestamp ntz types.
		{input: "datetime", want: DataTypeTimestampNTZ},
		{input: "timestamp", want: DataTypeTimestampNTZ},
		{input: "timestamp_ntz", want: DataTypeTimestampNTZ},

		// all othertypes
		{input: "date", want: DataTypeDate},
		{input: "time", want: DataTypeTime},
		{input: "timestamp_ltz", want: DataTypeTimestampLTZ},
		{input: "timestamp_tz", want: DataTypeTimestampTZ},
		{input: "variant", want: DataTypeVariant},
		{input: "object", want: DataTypeObject},
		{input: "array", want: DataTypeArray},
		{input: "geography", want: DataTypeGeography},
		{input: "geometry", want: DataTypeGeometry},
		{input: "invalid", want: DataTypeUnknown},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := DataTypeFromString(tc.input)
			require.Equal(t, tc.want, got)
		})
	}
}
