// Code generated by assertions generator; DO NOT EDIT.

package objectassert

import (
	"fmt"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ConnectionAssert struct {
	*assert.SnowflakeObjectAssert[sdk.Connection, sdk.AccountObjectIdentifier]
}

func Connection(t *testing.T, id sdk.AccountObjectIdentifier) *ConnectionAssert {
	t.Helper()
	return &ConnectionAssert{
		assert.NewSnowflakeObjectAssertWithProvider(sdk.ObjectTypeConnection, id, acc.TestClient().Connection.Show),
	}
}

func ConnectionFromObject(t *testing.T, connection *sdk.Connection) *ConnectionAssert {
	t.Helper()
	return &ConnectionAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectTypeConnection, connection.ID(), connection),
	}
}

func (c *ConnectionAssert) HasRegionGroup(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.RegionGroup == nil {
			return fmt.Errorf("expected region group to have value; got: nil")
		}
		if *o.RegionGroup != expected {
			return fmt.Errorf("expected region group: %v; got: %v", expected, *o.RegionGroup)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasSnowflakeRegion(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.SnowflakeRegion != expected {
			return fmt.Errorf("expected snowflake region: %v; got: %v", expected, o.SnowflakeRegion)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasCreatedOn(expected time.Time) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasAccountName(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.AccountName != expected {
			return fmt.Errorf("expected account name: %v; got: %v", expected, o.AccountName)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasName(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasComment(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.Comment == nil {
			return fmt.Errorf("expected comment to have value; got: nil")
		}
		if *o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, *o.Comment)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasIsPrimary(expected bool) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.IsPrimary != expected {
			return fmt.Errorf("expected is primary: %v; got: %v", expected, o.IsPrimary)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasPrimary(expected sdk.ExternalObjectIdentifier) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.Primary != expected {
			return fmt.Errorf("expected primary: %v; got: %v", expected, o.Primary)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasConnectionUrl(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.ConnectionUrl != expected {
			return fmt.Errorf("expected connection url: %v; got: %v", expected, o.ConnectionUrl)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasOrganizationName(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.OrganizationName != expected {
			return fmt.Errorf("expected organization name: %v; got: %v", expected, o.OrganizationName)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasAccountLocator(expected string) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.AccountLocator != expected {
			return fmt.Errorf("expected account locator: %v; got: %v", expected, o.AccountLocator)
		}
		return nil
	})
	return c
}
