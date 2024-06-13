package snowflakechecks

import (
	"errors"
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckDatabaseDataRetentionTimeInDays(t *testing.T, databaseId sdk.AccountObjectIdentifier, level string, value string) resource.TestCheckFunc {
	t.Helper()
	return func(state *terraform.State) error {
		param := helpers.FindParameter(t, acc.TestClient().Parameter.ShowDatabaseParameters(t, databaseId), sdk.AccountParameterDataRetentionTimeInDays)
		var errs []error
		if param.Level != sdk.ParameterType(level) {
			errs = append(errs, fmt.Errorf("expected parameter level %s, got %s", sdk.ParameterType(level), param.Level))
		}
		if param.Value != value {
			errs = append(errs, fmt.Errorf("expected parameter value %s, got %s", sdk.ParameterType(level), param.Level))
		}
		return errors.Join(errs...)
	}
}
