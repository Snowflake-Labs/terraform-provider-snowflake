package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

// TODO write a test for a candidate that's not castable to a string
func TestValidateIdentifier(t *testing.T) {
	cases := []struct {
		candidate string
		valid     bool
	}{
		{"word", true},
		{"_1", true},
		{"Aword", true},
		{"azAZ09_$", true},
		{"invalidcharacter!", false},
		{"1startwithnumber", false},
		{"$startwithdollar", false},
	}

	for _, tc := range cases {
		t.Run(tc.candidate, func(t *testing.T) {
			_, errs := snowflake.ValidateIdentifier(tc.candidate)
			actual := len(errs) == 0

			if actual == tc.valid {
				return
			}

			if tc.valid {
				t.Fatalf("identifier %s should pass validation", tc.candidate)
			} else {
				t.Fatalf("identifier %s should fail validation", tc.candidate)
			}
		})
	}
}
