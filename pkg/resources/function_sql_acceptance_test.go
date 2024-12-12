package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FunctionSql_InlineBasic(t *testing.T) {
	argName := "abc"
	dataType := testdatatypes.DataTypeFloat
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes()

	definition := acc.TestClient().Function.SampleSqlDefinitionWithArgument(t, argName)

	functionModel := model.FunctionSqlBasicInline("test", id, definition, dataType.ToLegacyDataTypeSql()).
		WithArgument(argName, dataType)
	functionModelRenamed := model.FunctionSqlBasicInline("test", idWithChangedNameButTheSameDataType, definition, dataType.ToLegacyDataTypeSql()).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionSql),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionSqlResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("SQL").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// IMPORT
			{
				ResourceName:            functionModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedFunctionSqlResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "FLOAT")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, functionModelRenamed),
				Check: assert.AssertThat(t,
					resourceassert.FunctionSqlResource(t, functionModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionSql_InlineFull(t *testing.T) {
	argName := "abc"
	comment := random.Comment()
	otherComment := random.Comment()
	dataType := testdatatypes.DataTypeFloat
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := acc.TestClient().Function.SampleSqlDefinitionWithArgument(t, argName)

	functionModel := model.FunctionSqlBasicInline("test", id, definition, dataType.ToLegacyDataTypeSql()).
		WithIsSecure(r.BooleanFalse).
		WithArgument(argName, dataType).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment(comment)
	functionModelRenamed := model.FunctionSqlBasicInline("test", idWithChangedNameButTheSameDataType, definition, dataType.ToLegacyDataTypeSql()).
		WithIsSecure(r.BooleanFalse).
		WithArgument(argName, dataType).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment(otherComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionSql),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionSqlResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasCommentString(comment).
						HasReturnBehaviorString(string(sdk.ReturnResultsBehaviorVolatile)).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("SQL").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// IMPORT
			{
				ResourceName:            functionModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"return_results_behavior"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedFunctionSqlResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "FLOAT")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, functionModelRenamed),
				Check: assert.AssertThat(t,
					resourceassert.FunctionSqlResource(t, functionModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()).
						HasCommentString(otherComment),
				),
			},
		},
	})
}
