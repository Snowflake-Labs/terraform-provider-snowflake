package sdk_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO write a test for a candidate that's not castable to a string.
func TestValidateIdentifier(t *testing.T) {
	cases := []struct {
		candidate string
		valid     bool
	}{
		{"word", true},
		{"_1", true},
		{"Aword", true},
		{"azAZ09_$", true},
		{"-30-Ab-", true},
		{"invalidcharacter!", false},
		{"1startwithnumber", true},
		{"$startwithdollar", false},
		{"[]includingBracket", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.candidate, func(t *testing.T) {
			_, errs := sdk.ValidateIdentifier(tc.candidate, []string{})
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
