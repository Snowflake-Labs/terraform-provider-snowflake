package helpers

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(client *sdk.Client, ctx context.Context) error {
	log.Printf("[DEBUG] Making sure QUOTED_IDENTIFIERS_IGNORE_CASE parameter is set correctly")
	param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterQuotedIdentifiersIgnoreCase)
	if err != nil {
		return fmt.Errorf("checking QUOTED_IDENTIFIERS_IGNORE_CASE resulted in error: %w", err)
	}
	if param.Value != "false" {
		return fmt.Errorf("parameter QUOTED_IDENTIFIERS_IGNORE_CASE has value %s, expected: false", param.Value)
	}
	return nil
}

func EnsureScimProvisionerRolesExist(client *sdk.Client, ctx context.Context) error {
	log.Printf("[DEBUG] Making sure Scim Provisioner roles exist")
	roleIDs := []sdk.AccountObjectIdentifier{snowflakeroles.GenericScimProvisioner, snowflakeroles.OktaProvisioner}
	currentRoleID, err := client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		_, err := client.Roles.ShowByID(ctx, roleID)
		if err != nil {
			return err
		}
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				Role: roleID,
			},
		})
		if err != nil {
			return err
		}
		if !hasGranteeName(grants, currentRoleID) {
			return fmt.Errorf("role %s not granted to %s", currentRoleID.Name(), roleID.Name())
		}
	}
	return nil
}

func hasGranteeName(grants []sdk.Grant, role sdk.AccountObjectIdentifier) bool {
	for _, grant := range grants {
		if grant.GranteeName == role {
			return true
		}
	}
	return false
}

// MatchAllStringsInOrderNonOverlapping returns a regex matching every string in parts. Matchings are non overlapping.
func MatchAllStringsInOrderNonOverlapping(parts []string) *regexp.Regexp {
	escapedParts := make([]string, len(parts))
	for i := range parts {
		escapedParts[i] = regexp.QuoteMeta(parts[i])
	}
	return regexp.MustCompile(strings.Join(escapedParts, "((.|\n)*)"))
}

// AssertErrorContainsPartsFunc returns a function asserting error message contains each string in parts
func AssertErrorContainsPartsFunc(t *testing.T, parts []string) resource.ErrorCheckFunc {
	t.Helper()
	return func(err error) error {
		for _, part := range parts {
			assert.Contains(t, err.Error(), part)
		}
		return nil
	}
}
