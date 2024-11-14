package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TagMaskingPolicyAssociationBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAttachmentConfig(accName, acc.TestDatabaseName, acc.TestSchemaName, tag.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "masking_policy_id", fmt.Sprintf("%s.%s.%s", acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "tag_id", tag.ID().FullyQualifiedName()),
				),
			},
		},
	})
}

func tagAttachmentConfig(n string, databaseName string, schemaName string, tagId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
	name = "%[1]v"
	database = "%[2]s"
	schema = "%[3]s"
	argument {
		name = "val"
		type = "VARCHAR"
	}
	body = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR(16777216)"
	comment = "Terraform acceptance test"
}

resource "snowflake_tag_masking_policy_association" "test" {
	tag_id = "\"%s\".\"%s\".\"%s\""
	masking_policy_id = "${snowflake_masking_policy.test.database}.${snowflake_masking_policy.test.schema}.${snowflake_masking_policy.test.name}"
}
`, n, databaseName, schemaName, tagId.DatabaseName(), tagId.SchemaName(), tagId.Name())
}
