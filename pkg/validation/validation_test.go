package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
