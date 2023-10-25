// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

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

var validAccounts = []string{
	"testOrg.testAcc",
	"testingOrg2.testingAcc2",
	"abc.12345",
}

var invalidAccounts = []interface{}{
	"abc12345",
	"xyz56789",
}

func TestValidateIsNotAccountLocator(t *testing.T) {
	r := require.New(t)
	for _, p := range validAccounts {
		_, errs := ValidateIsNotAccountLocator(p, "test_valid_accounts")
		r.Len(errs, 0, "account locators are not allowed - please use 'organization_name.account_name]", p, errs)
	}

	for _, p := range invalidAccounts {
		_, errs := ValidateIsNotAccountLocator(p, "test_invalid_accounts")
		r.NotZero(len(errs), "account locators are not allowed - please use 'organization_name.account_name]", p, errs)
	}
}
