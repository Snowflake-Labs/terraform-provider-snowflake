package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	tfconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicy(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"

	body := "case when current_role() in ('ANALYST') then true else false end"
	changedBody := "case when current_role() in ('CHANGED') then true else false end"
	argument := []sdk.RowAccessPolicyArgument{
		{
			Name: "A",
			Type: sdk.DataTypeVARCHAR,
		},
		{
			Name: "B",
			Type: sdk.DataTypeVARCHAR,
		},
	}
	changedArgument := []sdk.RowAccessPolicyArgument{
		{
			Name: "C",
			Type: sdk.DataTypeBoolean,
		},
		{
			Name: "D",
			Type: sdk.DataTypeTimestampNTZ,
		},
	}
	policyModel := model.RowAccessPolicy("test", argument, body, id.DatabaseName(), id.Name(), id.SchemaName()).WithComment("Terraform acceptance test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test").
					HasBodyString(body).
					HasArguments(argument),
					resourceshowoutputassert.RowAccessPolicyShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindRowAccessPolicy)).
						HasName(id.Name()).
						HasOptions("").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()).
						HasComment("Terraform acceptance test"),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.body", body)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.return_type", "BOOLEAN")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.name", "B")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.type", string(sdk.DataTypeVARCHAR))),
				),
			},
			// change comment and expression
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithBody(changedBody).WithComment("Terraform acceptance test - changed comment")),
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(changedBody).
					HasArguments(argument),
				),
			},
			// change signature
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithArgument(changedArgument)),
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			// external change on signature
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				PreConfig: func() {
					arg := sdk.NewCreateRowAccessPolicyArgsRequest("A", sdk.DataTypeBoolean)
					createRequest := sdk.NewCreateRowAccessPolicyRequest(id, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, "case when current_role() in ('ANALYST') then false else true end")
					acc.TestClient().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *createRequest.WithOrReplace(sdk.Pointer(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			// external change on body
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				PreConfig: func() {
					acc.TestClient().RowAccessPolicy.Alter(t, *sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.Pointer("case when current_role() in ('EXTERNAL') then false else true end")))
				},
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			{
				ConfigVariables:   tfconfig.ConfigVariablesFromModel(t, policyModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset comment
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithComment("")),
				PreConfig: func() {
					acc.TestClient().RowAccessPolicy.Alter(t, *sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.Pointer("case when current_role() in ('EXTERNAL') then false else true end")))
				},
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("").
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			// IMPORT
			{
				ConfigVariables:   tfconfig.ConfigVariablesFromModel(t, policyModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2053 is fixed
func TestAcc_RowAccessPolicy_Issue2053(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", []sdk.RowAccessPolicyArgument{
		{
			Name: "A",
			Type: sdk.DataTypeVARCHAR,
		},
	}, body, id.DatabaseName(), id.Name(), id.SchemaName())
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.95.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				// these configs have "weird" format on purpose - to test against handling new lines during diff correctly
				Config: rowAccessPolicy_v0_95_0_WithHeredoc(id, `    case
      when current_role() in ('ANALYST') then true
      else false
    end
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ExpectNonEmptyPlan: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables:          tfconfig.ConfigVariablesFromModel(t, policyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasBodyString(`case
  when current_role() in ('ANALYST') then true
  else false
end`),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_Rename(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", []sdk.RowAccessPolicyArgument{
		{
			Name: "a",
			Type: sdk.DataTypeVARCHAR,
		},
	}, body, id.DatabaseName(), id.Name(), id.SchemaName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// rename
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithName(newId.Name())),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(newId.Name()).
					HasFullyQualifiedNameString(newId.FullyQualifiedName()),
				),
			},
		},
	})
}

func rowAccessPolicy_v0_95_0(id sdk.SchemaObjectIdentifier, expr string) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  signature = {
    A = "VARCHAR",
    b = "VARCHAR",
  }
  row_access_expression = "%s"
}`, id.Name(), id.DatabaseName(), id.SchemaName(), expr)
}

func rowAccessPolicy_v0_95_0_WithHeredoc(id sdk.SchemaObjectIdentifier, expr string) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  signature = {
    A = "VARCHAR",
  }
  row_access_expression = <<-EOT
%s
EOT
}`, id.Name(), id.DatabaseName(), id.SchemaName(), expr)
}

func rowAccessPolicy_v0_96_0(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  argument {
    name = "A"
    type = "VARCHAR"
  }
  row_access_expression = <<-EOT
    case
      when current_role() in ('ANALYST') then true
      else false
    end
EOT
}`, id.Name(), id.DatabaseName(), id.SchemaName())
}

func TestAcc_RowAccessPolicy_InvalidDataType(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", []sdk.RowAccessPolicyArgument{
		{
			Name: "a",
			Type: "invalid-type",
		},
	}, body, id.DatabaseName(), id.Name(), id.SchemaName())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				ExpectError:     regexp.MustCompile(`invalid data type: invalid-type`),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_DataTypeAliases(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", []sdk.RowAccessPolicyArgument{
		{
			Name: "A",
			Type: "TEXT",
		},
	}, body, id.DatabaseName(), id.Name(), id.SchemaName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasArguments([]sdk.RowAccessPolicyArgument{
						{
							Name: "A",
							Type: sdk.DataTypeVARCHAR,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_view_migrateFromVersion_0_95_0_LowercaseArgName(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", []sdk.RowAccessPolicyArgument{
		{
			Name: "A",
			Type: sdk.DataTypeVARCHAR,
		},
		{
			Name: "b",
			Type: sdk.DataTypeVARCHAR,
		},
	}, body, id.DatabaseName(), id.Name(), id.SchemaName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.95.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: rowAccessPolicy_v0_95_0(id, body),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						// expect change - arg name is lower case which causes a diff
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "row_access_expression", body)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.A", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.B", string(sdk.DataTypeVARCHAR))),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables:          tfconfig.ConfigVariablesFromModel(t, policyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasBodyString(body).
					HasArguments([]sdk.RowAccessPolicyArgument{
						{
							Name: "A",
							Type: sdk.DataTypeVARCHAR,
						},
						{
							Name: "b",
							Type: sdk.DataTypeVARCHAR,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_view_migrateFromVersion_0_95_0_UppercaseArgName(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_row_access_policy.test"
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", []sdk.RowAccessPolicyArgument{
		{
			Name: "A",
			Type: sdk.DataTypeVARCHAR,
		},
		{
			Name: "B",
			Type: sdk.DataTypeVARCHAR,
		},
	}, body, id.DatabaseName(), id.Name(), id.SchemaName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.95.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: rowAccessPolicy_v0_95_0(id, body),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						// expect change - arg name is lower case which causes a diff
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "row_access_expression", body)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.A", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.B", string(sdk.DataTypeVARCHAR))),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables:          tfconfig.ConfigVariablesFromModel(t, policyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: assert.AssertThat(t, resourceassert.RowAccessPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasBodyString(body).
					HasArguments([]sdk.RowAccessPolicyArgument{
						{
							Name: "A",
							Type: sdk.DataTypeVARCHAR,
						},
						{
							Name: "B",
							Type: sdk.DataTypeVARCHAR,
						},
					}),
				),
			},
		},
	})
}
