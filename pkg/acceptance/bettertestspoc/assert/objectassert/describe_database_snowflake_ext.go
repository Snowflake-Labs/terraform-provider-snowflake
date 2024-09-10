package objectassert

import (
	"fmt"
	"slices"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO [SNOW-1501905]: this file should be fully regenerated when adding and option to assert the results of describe
type DatabaseDescribeAssert struct {
	*assert.SnowflakeObjectAssert[sdk.DatabaseDetails, sdk.AccountObjectIdentifier]
}

func DatabaseDescribe(t *testing.T, id sdk.AccountObjectIdentifier) *DatabaseDescribeAssert {
	t.Helper()
	return &DatabaseDescribeAssert{
		assert.NewSnowflakeObjectAssertWithProvider(sdk.ObjectType("DATABASE_DETAILS"), id, acc.TestClient().Database.Describe),
	}
}

func (d *DatabaseDescribeAssert) DoesNotContainPublicSchema() *DatabaseDescribeAssert {
	d.AddAssertion(func(t *testing.T, o *sdk.DatabaseDetails) error {
		t.Helper()
		if slices.ContainsFunc(o.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to not contain public schema", d.GetId())
		}
		return nil
	})
	return d
}

func (d *DatabaseDescribeAssert) ContainsPublicSchema() *DatabaseDescribeAssert {
	d.AddAssertion(func(t *testing.T, o *sdk.DatabaseDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to contain public schema", d.GetId())
		}
		return nil
	})
	return d
}
