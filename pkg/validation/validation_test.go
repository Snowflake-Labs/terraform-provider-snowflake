package validation

import (
	"testing"
)

var validAccounts = []string{
	"bob",
	"sue",
	"spaces are valid",
	"special$characters$are$valid",
}

var invalidAccounts = []interface{}{
	"123asd",
	"=**",
	"\"",
	123,
}

func TestValidateAccount(t *testing.T) {
	for _, a := range validAccounts {
		_, errs := ValidateAccount(a, "test_account")
		if len(errs) != 0 {
			t.Errorf("%v failed to validate: %v", a, errs)
		}
	}

	for _, a := range invalidAccounts {
		_, errs := ValidateAccount(a, "test_account")
		if len(errs) == 0 {
			t.Errorf("%v should have failed to validate: %v", a, errs)
		}
	}
}

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
	for _, p := range validPasswords {
		_, errs := ValidatePassword(p, "test_password")
		if len(errs) != 0 {
			t.Errorf("%v failed to validate: %v", p, errs)
		}
	}

	for _, p := range invalidPasswords {
		_, errs := ValidatePassword(p, "test_password")
		if len(errs) == 0 {
			t.Errorf("%v should have failed to validate: %v", p, errs)
		}
	}
}
