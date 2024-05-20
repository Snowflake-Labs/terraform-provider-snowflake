package snowflakesql

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool_Scan(t *testing.T) {
	type testCase struct {
		from     any
		expected Bool
		err      error
	}

	for _, tc := range []testCase{
		{
			// passing nil will result in invalid Bool without errors
			from:     nil,
			expected: Bool{},
			err:      nil,
		},
		{
			from:     "1",
			expected: Bool{Valid: true, Bool: true},
		},
		{
			from:     &[]string{"1"}[0], // pointer to string
			expected: Bool{Valid: true, Bool: true},
		},
		{
			from:     "2",
			expected: Bool{},
			err:      errors.New("sql/driver: couldn't convert \"2\" into type bool"),
		},
		{
			from:     &[]string{"2"}[0], // pointer to string
			expected: Bool{},
			err:      errors.New("sql/driver: couldn't convert \"2\" into type bool"),
		},
		{
			from:     "",
			expected: Bool{},
			err:      errors.New("sql/driver: couldn't convert \"\" into type bool"),
		},
		{
			from:     &[]string{""}[0], // pointer to string
			expected: Bool{},
			err:      errors.New("sql/driver: couldn't convert \"\" into type bool"),
		},
		{
			from:     true,
			expected: Bool{Valid: true, Bool: true},
		},
		{
			from:     &[]bool{true}[0], // pointer to bool
			expected: Bool{Valid: true, Bool: true},
		},
		{
			from:     false,
			expected: Bool{Valid: true, Bool: false},
		},
		{
			from:     &[]bool{false}[0], // pointer to bool
			expected: Bool{Valid: true, Bool: false},
		},
		{
			from:     int64(123),
			expected: Bool{},
			err:      errors.New("sql/driver: couldn't convert 123 into type bool"),
		},
		{
			from:     &[]int64{123}[0], // pointer to int64
			expected: Bool{},
			err:      errors.New("sql/driver: couldn't convert 123 into type bool"),
		},
		{
			from:     io.Copy,
			expected: Bool{},
			err:      errors.New("(func(io.Writer, io.Reader) (int64, error)) into type bool"),
		},
	} {
		t.Run(fmt.Sprint(tc.from), func(t *testing.T) {
			var res Bool
			err := res.Scan(tc.from)
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err.Error())
			}
			assert.Exactly(t, tc.expected, res)
		})
	}
}
