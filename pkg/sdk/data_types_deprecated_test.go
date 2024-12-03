package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsStringType(t *testing.T) {
	type test struct {
		input string
		want  bool
	}

	tests := []test{
		// case insensitive.
		{input: "STRING", want: true},
		{input: "string", want: true},
		{input: "String", want: true},

		// varchar types.
		{input: "VARCHAR", want: true},
		{input: "NVARCHAR", want: true},
		{input: "NVARCHAR2", want: true},
		{input: "CHAR", want: true},
		{input: "NCHAR", want: true},
		{input: "CHAR VARYING", want: true},
		{input: "NCHAR VARYING", want: true},
		{input: "TEXT", want: true},

		// with length
		{input: "VARCHAR(100)", want: true},
		{input: "NVARCHAR(100)", want: true},
		{input: "NVARCHAR2(100)", want: true},
		{input: "CHAR(100)", want: true},
		{input: "NCHAR(100)", want: true},
		{input: "CHAR VARYING(100)", want: true},
		{input: "NCHAR VARYING(100)", want: true},
		{input: "TEXT(100)", want: true},

		// binary is not string types.
		{input: "binary", want: false},
		{input: "varbinary", want: false},

		// other types
		{input: "boolean", want: false},
		{input: "number", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := IsStringType(tc.input)
			require.Equal(t, tc.want, got)
		})
	}
}
