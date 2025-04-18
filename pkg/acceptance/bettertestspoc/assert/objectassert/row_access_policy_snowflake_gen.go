// Code generated by assertions generator; DO NOT EDIT.

package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type RowAccessPolicyAssert struct {
	*assert.SnowflakeObjectAssert[sdk.RowAccessPolicy, sdk.SchemaObjectIdentifier]
}

func RowAccessPolicy(t *testing.T, id sdk.SchemaObjectIdentifier) *RowAccessPolicyAssert {
	t.Helper()
	return &RowAccessPolicyAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeRowAccessPolicy, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.RowAccessPolicy, sdk.SchemaObjectIdentifier] {
			return testClient.RowAccessPolicy.Show
		}),
	}
}

func RowAccessPolicyFromObject(t *testing.T, rowAccessPolicy *sdk.RowAccessPolicy) *RowAccessPolicyAssert {
	t.Helper()
	return &RowAccessPolicyAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectTypeRowAccessPolicy, rowAccessPolicy.ID(), rowAccessPolicy),
	}
}

func (r *RowAccessPolicyAssert) HasCreatedOn(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasName(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasDatabaseName(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.DatabaseName != expected {
			return fmt.Errorf("expected database name: %v; got: %v", expected, o.DatabaseName)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasSchemaName(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.SchemaName != expected {
			return fmt.Errorf("expected schema name: %v; got: %v", expected, o.SchemaName)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasKind(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.Kind != expected {
			return fmt.Errorf("expected kind: %v; got: %v", expected, o.Kind)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasOwner(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, o.Owner)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasComment(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasOptions(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.Options != expected {
			return fmt.Errorf("expected options: %v; got: %v", expected, o.Options)
		}
		return nil
	})
	return r
}

func (r *RowAccessPolicyAssert) HasOwnerRoleType(expected string) *RowAccessPolicyAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.RowAccessPolicy) error {
		t.Helper()
		if o.OwnerRoleType != expected {
			return fmt.Errorf("expected owner role type: %v; got: %v", expected, o.OwnerRoleType)
		}
		return nil
	})
	return r
}
