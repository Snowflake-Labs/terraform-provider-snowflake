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
