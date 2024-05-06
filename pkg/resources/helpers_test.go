package resources_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

type grantType int

const (
	normal grantType = iota
	onFuture
	onAll
)

func TestGetPropertyAsPointer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}, map[string]interface{}{
		"integer": 123,
		"string":  "some string",
		"boolean": true,
		"invalid": true,
	})

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))
}

// queriedAccountRolePrivilegesEqualTo will check if all the privileges specified in the argument are granted in Snowflake.
func queriedPrivilegesEqualTo(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			if (grant.GrantTo == sdk.ObjectTypeDatabaseRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole) && grant.Privilege == "USAGE" {
				continue
			}
			if !slices.Contains(privileges, grant.Privilege) {
				return fmt.Errorf("grant not expected, grant: %v, not in %v", grants, privileges)
			}
		}

		return nil
	}
}

// queriedAccountRolePrivilegesContainAtLeast will check if all the privileges specified in the argument are granted in Snowflake.
// Any additional grants will be ignored.
func queriedPrivilegesContainAtLeast(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), roleName sdk.ObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		var grantedPrivileges []string
		for _, grant := range grants {
			grantedPrivileges = append(grantedPrivileges, grant.Privilege)
		}
		notAllPrivilegesInGrantedPrivileges := slices.ContainsFunc(privileges, func(privilege string) bool {
			return !slices.Contains(grantedPrivileges, privilege)
		})
		if notAllPrivilegesInGrantedPrivileges {
			return fmt.Errorf("not every privilege from the list: %v was found in grant privileges: %v, for role name: %s", privileges, grantedPrivileges, roleName.FullyQualifiedName())
		}

		return nil
	}
}
