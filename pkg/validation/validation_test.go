package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var validPasswords = []string{
	"bob123BOB",
	"actually-quite-a-strong-pASSword-90210",
}

var invalidPasswords = []interface{}{
	"password",
	"password123",
	"birthday",
	"1982",
	"alM0st!",
	123,
}

func TestValidatePassword(t *testing.T) {
	r := require.New(t)
	for _, p := range validPasswords {
		_, errs := ValidatePassword(p, "test_password")
		r.Len(errs, 0, "%v failed to validate: %v", p, errs)
	}

	for _, p := range invalidPasswords {
		_, errs := ValidatePassword(p, "test_password")
		r.NotZero(len(errs), "%v should have failed to validate: %v", p, errs)
	}
}

func TestValidatePrivilege(t *testing.T) {
	r := require.New(t)

	// even if we "allow" ALL, error out
	w, errs := ValidatePrivilege([]string{"ALL"}, true)("ALL", "unused")
	r.Empty(w)
	r.Len(errs, 1)
	r.Equal(
		"the ALL privilege is deprecated, see https://github.com/chanzuckerberg/terraform-provider-snowflake/discussions/318",
		errs[0].Error(),
	)

	// fail if not in set
	w, errs = ValidatePrivilege([]string{"YES"}, true)("NO", "unused")
	r.Empty(w)
	r.Len(errs, 1)
	r.Equal("expected unused to be one of [YES], got NO", errs[0].Error())

	// success
	w, errs = ValidatePrivilege([]string{"YES"}, true)("YES", "unused")
	r.Empty(w)
	r.Empty(errs)
}
