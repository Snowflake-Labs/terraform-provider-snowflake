package resources_test

import (
	"context"
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
	"regexp"
	"strings"
	"testing"
)

func TestAcc_Provider_UseSecondaryRoles(t *testing.T) {
	providerRole := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	providerUseSecondaryRolesSetup(t, providerRole)

	providerConfigVariables := config.Variables{
		"profile":             config.StringVariable("default"),
		"role":                config.StringVariable(providerRole),
		"use_secondary_roles": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: providerUseSecondarySchemaConfig(providerConfig(t, providerConfigVariables), databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.test", "name", databaseName),
				),
			},
		},
	})
}

func TestAcc_Provider_UseSecondaryRolesUnchecked(t *testing.T) {
	providerRole := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	//providerUseSecondaryRolesSetup(t, providerRole)

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	providerConfigVariables := config.Variables{
		"profile":             config.StringVariable("default"),
		"role":                config.StringVariable(providerRole),
		"use_secondary_roles": config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      providerUseSecondarySchemaConfig(providerConfig(t, providerConfigVariables), databaseName),
				ExpectError: regexp.MustCompile("Insufficient privileges"),
			},
		},
	})
}

func providerUseSecondarySchemaConfig(providerConfig string, databaseName string) string {
	return fmt.Sprintf(`
%s

resource "snowflake_database" "test" {
  name = "%s"
}`, providerConfig, databaseName)
}

func providerConfig(t *testing.T, variables config.Variables) string {
	t.Helper()
	var builder strings.Builder
	for k, v := range variables {
		builder.WriteString(k)
		builder.WriteString(" = ")
		value, err := v.MarshalJSON()
		require.NoError(t, err)
		builder.Write(value)
		builder.WriteByte('\n')
	}
	return fmt.Sprintf(`provider "snowflake" {
%s
}`, builder.String())
}

func providerUseSecondaryRolesSetup(t *testing.T, providerRole string) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)

	ctx := context.Background()
	secondaryRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	secondaryRoleId := sdk.NewAccountObjectIdentifier(secondaryRoleName)

	providerRoleId := sdk.NewAccountObjectIdentifier(providerRole)

	createTestRole(t, client, ctx, providerRole)
	createTestRole(t, client, ctx, secondaryRoleName)

	currentUserName, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(secondaryRoleId, sdk.GrantRole{
		User: sdk.Pointer(sdk.NewAccountObjectIdentifier(currentUserName)),
	}))
	require.NoError(t, err)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(providerRoleId, sdk.GrantRole{
		User: sdk.Pointer(sdk.NewAccountObjectIdentifier(currentUserName)),
	}))
	require.NoError(t, err)

	err = client.Grants.GrantPrivilegesToAccountRole(ctx, &sdk.AccountRoleGrantPrivileges{
		GlobalPrivileges: []sdk.GlobalPrivilege{
			sdk.GlobalPrivilegeCreateDatabase,
		},
	},
		&sdk.AccountRoleGrantOn{
			Account: sdk.Bool(true),
		},
		sdk.NewAccountObjectIdentifier(secondaryRoleName),
		&sdk.GrantPrivilegesToAccountRoleOptions{},
	)
	require.NoError(t, err)

	//err = client.Sessions.UseRole(ctx, secondaryRoleId)
	//require.NoError(t, err)
	//
	//err = client.Databases.Create(ctx, sdk.NewAccountObjectIdentifier(databaseName), new(sdk.CreateDatabaseOptions))
	//require.NoError(t, err)
	//
	//t.Cleanup(func() {
	//	err = client.Databases.Drop(ctx, sdk.NewAccountObjectIdentifier(databaseName), new(sdk.DropDatabaseOptions))
	//	require.NoError(t, err)
	//})
}

func createTestRole(t *testing.T, client *sdk.Client, ctx context.Context, roleName string) {
	t.Helper()
	id := sdk.NewAccountObjectIdentifier(roleName)

	err := client.Roles.Create(ctx, sdk.NewCreateRoleRequest(id))
	require.NoError(t, err)

	t.Cleanup(func() {
		err := client.Roles.Drop(ctx, sdk.NewDropRoleRequest(id))
		require.NoError(t, err)
	})
	require.NoError(t, err)
}
