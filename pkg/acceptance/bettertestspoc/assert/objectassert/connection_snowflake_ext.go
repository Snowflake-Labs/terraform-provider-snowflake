package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionAssert) HasFailoverAllowedToAccounts(expected []string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if !slices.Equal(expected, o.FailoverAllowedToAccounts) {
			return fmt.Errorf("expected failover_allowed_to_accounts to be: %v; got: %v", expected, o.FailoverAllowedToAccounts)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasNoComment() *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to have nil; got: %s", *o.Comment)
		}
		return nil
	})
	return c
}
