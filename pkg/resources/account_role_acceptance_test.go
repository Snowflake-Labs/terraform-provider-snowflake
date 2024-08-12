package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountRole_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	currentRole, err := acc.Client(t).ContextFunctions.CurrentRole(context.Background())
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accountRoleBasicConfig(id.Name(), ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_account_role.role", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.comment", ""),
				),
			},
			// import - without optionals
			{
				Config:       accountRoleBasicConfig(id.Name(), ""),
				ResourceName: "snowflake_account_role.role",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", ""),
				),
			},
			// set optionals
			{
				Config: accountRoleBasicConfig(id.Name(), comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_account_role.role", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_account_role.role", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.comment", comment),
				),
			},
			// import - complete
			{
				Config:       accountRoleBasicConfig(id.Name(), ""),
				ResourceName: "snowflake_account_role.role",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", comment),
				),
			},
			// unset
			{
				Config: accountRoleBasicConfig(id.Name(), ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_account_role.role", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "comment", ""),

					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_account_role.role", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.comment", ""),
				),
			},
		},
	})
}

func TestAcc_AccountRole_Complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	currentRole, err := acc.Client(t).ContextFunctions.CurrentRole(context.Background())
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			{
				Config: accountRoleBasicConfig(id.Name(), comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "comment", comment),

					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_account_role.role", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.comment", comment),
				),
			},
			{
				Config:       accountRoleBasicConfig(id.Name(), ""),
				ResourceName: "snowflake_account_role.role",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", comment),
				),
			},
			// rename + comment change
			{
				Config: accountRoleBasicConfig(newId.Name(), newComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "comment", newComment),

					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_account_role.role", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "show_output.0.comment", newComment),
				),
			},
		},
	})
}

func accountRoleBasicConfig(name, comment string) string {
	s := `
resource "snowflake_account_role" "role" {
	name = "%s"
	comment = "%s"
}
`
	return fmt.Sprintf(s, name, comment)
}

func TestAcc_AccountRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: accountRoleBasicConfig(id.Name(), comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accountRoleBasicConfig(id.Name(), comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_AccountRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             accountRoleBasicConfig(quotedId, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accountRoleBasicConfig(quotedId, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
		},
	})
}
