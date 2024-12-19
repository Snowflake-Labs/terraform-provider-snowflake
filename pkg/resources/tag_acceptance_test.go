package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tag_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	maskingPolicy, maskingPolicyCleanup := acc.TestClient().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicyCleanup)

	maskingPolicy2, maskingPolicy2Cleanup := acc.TestClient().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicy2Cleanup)

	baseModel := model.Tag("test", id.DatabaseName(), id.Name(), id.SchemaName())

	modelWithExtraFields := model.Tag("test", id.DatabaseName(), id.Name(), id.SchemaName()).
		WithComment("foo").
		WithAllowedValues("foo", "", "bar").
		WithMaskingPolicies(maskingPolicy.ID())

	modelWithDifferentListOrder := model.Tag("test", id.DatabaseName(), id.Name(), id.SchemaName()).
		WithComment("foo").
		WithAllowedValues("", "bar", "foo").
		WithMaskingPolicies(maskingPolicy.ID())

	modelWithDifferentValues := model.Tag("test", id.DatabaseName(), id.Name(), id.SchemaName()).
		WithComment("bar").
		WithAllowedValues("abc", "def", "").
		WithMaskingPolicies(maskingPolicy2.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			// base model
			{
				Config: config.FromModels(t, baseModel),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, baseModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("").
					HasNoMaskingPolicies().
					HasNoAllowedValues(),
					resourceshowoutputassert.TagShowOutput(t, baseModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasNoAllowedValues(),
				),
			},
			// import without optionals
			{
				Config:            config.FromModels(t, baseModel),
				ResourceName:      baseModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// set all fields
			{
				Config: config.FromModels(t, modelWithExtraFields),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, modelWithExtraFields.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("foo"),
					assert.Check(resource.TestCheckResourceAttr(modelWithExtraFields.ResourceReference(), "masking_policies.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "masking_policies.*", maskingPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithExtraFields.ResourceReference(), "allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "allowed_values.*", "bar")),
					resourceshowoutputassert.TagShowOutput(t, modelWithExtraFields.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("foo").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.*", "bar")),
				),
			},
			// external change
			{
				PreConfig: func() {
					acc.TestClient().Tag.Alter(t, sdk.NewAlterTagRequest(id).WithDrop([]string{"foo"}))
				},
				Config: config.FromModels(t, modelWithExtraFields),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, modelWithExtraFields.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("foo"),
					assert.Check(resource.TestCheckResourceAttr(modelWithExtraFields.ResourceReference(), "masking_policies.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "masking_policies.*", maskingPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithExtraFields.ResourceReference(), "allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "allowed_values.*", "bar")),
					resourceshowoutputassert.TagShowOutput(t, modelWithExtraFields.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("foo").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithExtraFields.ResourceReference(), "show_output.0.allowed_values.*", "bar")),
				),
			},
			// different set ordering
			{
				Config: config.FromModels(t, modelWithDifferentListOrder),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithDifferentListOrder.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assert.AssertThat(t, resourceassert.TagResource(t, modelWithDifferentListOrder.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("foo"),
					assert.Check(resource.TestCheckResourceAttr(modelWithDifferentListOrder.ResourceReference(), "masking_policies.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentListOrder.ResourceReference(), "masking_policies.*", maskingPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithDifferentListOrder.ResourceReference(), "allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentListOrder.ResourceReference(), "allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentListOrder.ResourceReference(), "allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentListOrder.ResourceReference(), "allowed_values.*", "bar")),
					resourceshowoutputassert.TagShowOutput(t, modelWithDifferentListOrder.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("foo").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(modelWithDifferentListOrder.ResourceReference(), "show_output.0.allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentListOrder.ResourceReference(), "show_output.0.allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentListOrder.ResourceReference(), "show_output.0.allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentListOrder.ResourceReference(), "show_output.0.allowed_values.*", "bar")),
				),
			},
			// change some values
			{
				Config: config.FromModels(t, modelWithDifferentValues),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, modelWithDifferentValues.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("bar"),
					assert.Check(resource.TestCheckResourceAttr(modelWithDifferentValues.ResourceReference(), "masking_policies.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentValues.ResourceReference(), "masking_policies.*", maskingPolicy2.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelWithDifferentValues.ResourceReference(), "allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentValues.ResourceReference(), "allowed_values.*", "abc")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentValues.ResourceReference(), "allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentValues.ResourceReference(), "allowed_values.*", "def")),
					resourceshowoutputassert.TagShowOutput(t, modelWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("bar").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(modelWithDifferentValues.ResourceReference(), "show_output.0.allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentValues.ResourceReference(), "show_output.0.allowed_values.*", "abc")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentValues.ResourceReference(), "show_output.0.allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(modelWithDifferentValues.ResourceReference(), "show_output.0.allowed_values.*", "def")),
				),
			},
			// unset optionals
			{
				Config: config.FromModels(t, baseModel),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, baseModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("").
					HasMaskingPoliciesLength(0).
					HasAllowedValuesLength(0),
					resourceshowoutputassert.TagShowOutput(t, baseModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasNoAllowedValues(),
				),
			},
		},
	})
}

func TestAcc_Tag_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	maskingPolicy, maskingPolicyCleanup := acc.TestClient().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(maskingPolicyCleanup)

	model := model.Tag("test", id.DatabaseName(), id.Name(), id.SchemaName()).
		WithComment("foo").
		WithAllowedValuesValue(tfconfig.ListVariable(tfconfig.StringVariable("foo"), tfconfig.StringVariable(""), tfconfig.StringVariable("bar"))).
		WithMaskingPoliciesValue(tfconfig.ListVariable(tfconfig.StringVariable(maskingPolicy.ID().FullyQualifiedName())))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, model),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, model.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("foo"),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "masking_policies.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "masking_policies.*", maskingPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "allowed_values.*", "bar")),
					resourceshowoutputassert.TagShowOutput(t, model.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("foo").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "show_output.0.allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "show_output.0.allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "show_output.0.allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "show_output.0.allowed_values.*", "bar")),
				),
			},
			{
				Config:            config.FromModels(t, model),
				ResourceName:      model.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_Tag_Rename(t *testing.T) {
	oldId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	modelWithOldId := model.Tag("test", oldId.DatabaseName(), oldId.Name(), oldId.SchemaName())
	modelWithNewId := model.Tag("test", newId.DatabaseName(), newId.Name(), newId.SchemaName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Tag),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, modelWithOldId),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, modelWithOldId.ResourceReference()).
					HasNameString(oldId.Name()).
					HasDatabaseString(oldId.DatabaseName()).
					HasSchemaString(oldId.SchemaName()),
				),
			},
			{
				Config: config.FromModels(t, modelWithNewId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithOldId.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.TagResource(t, modelWithNewId.ResourceReference()).
					HasNameString(newId.Name()).
					HasDatabaseString(newId.DatabaseName()).
					HasSchemaString(newId.SchemaName()),
				),
			},
		},
	})
}

func TestAcc_Tag_migrateFromVersion_0_98_0(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	model := model.Tag("test", id.DatabaseName(), id.Name(), id.SchemaName()).
		WithAllowedValuesValue(tfconfig.ListVariable(tfconfig.StringVariable("foo"), tfconfig.StringVariable("bar")))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.98.0"),
				Config:            tag_v_0_98_0(id),
				Check: assert.AssertThat(t, resourceassert.TagResource(t, model.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "allowed_values.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "allowed_values.0", "bar")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "allowed_values.1", "foo")),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, model),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assert.AssertThat(t, resourceassert.TagResource(t, model.ResourceReference()).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "allowed_values.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(model.ResourceReference(), "allowed_values.*", "bar")),
				),
			},
		},
	})
}

func tag_v_0_98_0(id sdk.SchemaObjectIdentifier) string {
	s := `
resource "snowflake_tag" "test" {
	name					= "%[1]s"
	database				= "%[2]s"
	schema				    = "%[3]s"
	allowed_values			= ["bar", "foo"]
}
`
	return fmt.Sprintf(s, id.Name(), id.DatabaseName(), id.SchemaName())
}
