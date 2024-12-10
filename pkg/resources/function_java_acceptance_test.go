package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FunctionJava_BasicFlows(t *testing.T) {
	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100
	// differentDataType := testdatatypes.DataTypeNumber_36_2

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	// idWithSameNameButDifferentDataType := acc.TestClient().Ids.NewSchemaObjectIdentifierWithArgumentsNewDataTypes(idWithChangedNameButTheSameDataType.Name(), differentDataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	functionModelNoAttributes := model.FunctionJavaWithId("w", id, dataType, handler, definition)
	functionModelNoAttributesRenamed := model.FunctionJavaWithId("w", idWithChangedNameButTheSameDataType, dataType, handler, definition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.ResourceFromModel(t, functionModelNoAttributes),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModelNoAttributes.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(sdk.DefaultFunctionComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModelNoAttributes.ResourceReference()).
						HasIsSecure(false),
				),
			},
			// RENAME
			{
				Config: config.ResourceFromModel(t, functionModelNoAttributesRenamed),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModelNoAttributesRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
			//// IMPORT
			//{
			//	ResourceName:            userModelNoAttributesRenamed.ResourceReference(),
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"password", "disable_mfa", "days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "login_name", "display_name", "disabled", "must_change_password"},
			//	ImportStateCheck: assert.AssertThatImport(t,
			//		resourceassert.ImportedUserResource(t, id2.Name()).
			//			HasLoginNameString(strings.ToUpper(id.Name())).
			//			HasDisplayNameString(id.Name()).
			//			HasDisabled(false).
			//			HasMustChangePassword(false),
			//	),
			//},
			//// DESTROY
			//{
			//	Config:  config.FromModel(t, userModelNoAttributes),
			//	Destroy: true,
			//},
			//// CREATE WITH ALL ATTRIBUTES
			//{
			//	Config: config.FromModel(t, userModelAllAttributes),
			//	Check: assert.AssertThat(t,
			//		resourceassert.UserResource(t, userModelAllAttributes.ResourceReference()).
			//			HasNameString(id.Name()).
			//			HasPasswordString(pass).
			//			HasLoginNameString(fmt.Sprintf("%s_login", id.Name())).
			//			HasDisplayNameString("Display Name").
			//			HasFirstNameString("Jan").
			//			HasMiddleNameString("Jakub").
			//			HasLastNameString("Testowski").
			//			HasEmailString("fake@email.com").
			//			HasMustChangePassword(true).
			//			HasDisabled(false).
			//			HasDaysToExpiryString("8").
			//			HasMinsToUnlockString("9").
			//			HasDefaultWarehouseString("some_warehouse").
			//			HasDefaultNamespaceString("some.namespace").
			//			HasDefaultRoleString("some_role").
			//			HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
			//			HasMinsToBypassMfaString("10").
			//			HasRsaPublicKeyString(key1).
			//			HasRsaPublicKey2String(key2).
			//			HasCommentString(comment).
			//			HasDisableMfaString(r.BooleanTrue).
			//			HasFullyQualifiedNameString(id.FullyQualifiedName()),
			//	),
			//},
			//// CHANGE PROPERTIES
			//{
			//	Config: config.FromModel(t, userModelAllAttributesChanged(id.Name()+"_other_login")),
			//	Check: assert.AssertThat(t,
			//		resourceassert.UserResource(t, userModelAllAttributesChanged(id.Name()+"_other_login").ResourceReference()).
			//			HasNameString(id.Name()).
			//			HasPasswordString(newPass).
			//			HasLoginNameString(fmt.Sprintf("%s_other_login", id.Name())).
			//			HasDisplayNameString("New Display Name").
			//			HasFirstNameString("Janek").
			//			HasMiddleNameString("Kuba").
			//			HasLastNameString("Terraformowski").
			//			HasEmailString("fake@email.net").
			//			HasMustChangePassword(false).
			//			HasDisabled(true).
			//			HasDaysToExpiryString("12").
			//			HasMinsToUnlockString("13").
			//			HasDefaultWarehouseString("other_warehouse").
			//			HasDefaultNamespaceString("one_part_namespace").
			//			HasDefaultRoleString("other_role").
			//			HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionAll).
			//			HasMinsToBypassMfaString("14").
			//			HasRsaPublicKeyString(key2).
			//			HasRsaPublicKey2String(key1).
			//			HasCommentString(newComment).
			//			HasDisableMfaString(r.BooleanFalse).
			//			HasFullyQualifiedNameString(id.FullyQualifiedName()),
			//	),
			//},
			//// IMPORT
			//{
			//	ResourceName:            userModelAllAttributesChanged(id.Name() + "_other_login").ResourceReference(),
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"password", "disable_mfa", "days_to_expiry", "mins_to_unlock", "mins_to_bypass_mfa", "default_namespace", "login_name", "show_output.0.days_to_expiry"},
			//	ImportStateCheck: assert.AssertThatImport(t,
			//		resourceassert.ImportedUserResource(t, id.Name()).
			//			HasDefaultNamespaceString("ONE_PART_NAMESPACE").
			//			HasLoginNameString(fmt.Sprintf("%s_OTHER_LOGIN", id.Name())),
			//	),
			//},
			//// CHANGE PROP TO THE CURRENT SNOWFLAKE VALUE
			//{
			//	PreConfig: func() {
			//		acc.TestClient().User.SetLoginName(t, id, id.Name()+"_different_login")
			//	},
			//	Config: config.FromModel(t, userModelAllAttributesChanged(id.Name()+"_different_login")),
			//	ConfigPlanChecks: resource.ConfigPlanChecks{
			//		PostApplyPostRefresh: []plancheck.PlanCheck{
			//			plancheck.ExpectEmptyPlan(),
			//		},
			//	},
			//},
			//// UNSET ALL
			//{
			//	Config: config.FromModel(t, userModelNoAttributes),
			//	Check: assert.AssertThat(t,
			//		resourceassert.UserResource(t, userModelNoAttributes.ResourceReference()).
			//			HasNameString(id.Name()).
			//			HasPasswordString("").
			//			HasLoginNameString("").
			//			HasDisplayNameString("").
			//			HasFirstNameString("").
			//			HasMiddleNameString("").
			//			HasLastNameString("").
			//			HasEmailString("").
			//			HasMustChangePasswordString(r.BooleanDefault).
			//			HasDisabledString(r.BooleanDefault).
			//			HasDaysToExpiryString("0").
			//			HasMinsToUnlockString(r.IntDefaultString).
			//			HasDefaultWarehouseString("").
			//			HasDefaultNamespaceString("").
			//			HasDefaultRoleString("").
			//			HasDefaultSecondaryRolesOption(sdk.SecondaryRolesOptionDefault).
			//			HasMinsToBypassMfaString(r.IntDefaultString).
			//			HasRsaPublicKeyString("").
			//			HasRsaPublicKey2String("").
			//			HasCommentString("").
			//			HasDisableMfaString(r.BooleanDefault).
			//			HasFullyQualifiedNameString(id.FullyQualifiedName()),
			//		resourceshowoutputassert.UserShowOutput(t, userModelNoAttributes.ResourceReference()).
			//			HasLoginName(strings.ToUpper(id.Name())).
			//			HasDisplayName(""),
			//	),
			//},
		},
	})
}
