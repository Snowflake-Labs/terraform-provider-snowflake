package snowflakechecks

import (
	"fmt"
	"slices"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func DoesNotContainPublicSchema(t *testing.T, id sdk.AccountObjectIdentifier) resource.TestCheckFunc {
	t.Helper()
	return func(state *terraform.State) error {
		if slices.ContainsFunc(acc.TestClient().Database.Describe(t, id).Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to not contain public schema", id.FullyQualifiedName())
		}
		return nil
	}
}

func ContainsPublicSchema(t *testing.T, id sdk.AccountObjectIdentifier) resource.TestCheckFunc {
	t.Helper()
	return func(state *terraform.State) error {
		if !slices.ContainsFunc(acc.TestClient().Database.Describe(t, id).Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to contain public schema", id.FullyQualifiedName())
		}
		return nil
	}
}
