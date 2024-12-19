package datasources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	testconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tags(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	model := model.Tag("test", id.DatabaseName(), id.Name(), id.SchemaName()).
		WithComment("foo").
		WithAllowedValuesValue(tfconfig.ListVariable(tfconfig.StringVariable("foo"), tfconfig.StringVariable(""), tfconfig.StringVariable("bar")))

	dsName := "data.snowflake_tags.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Tags/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, model),

				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "tags.#", "1")),

					resourceshowoutputassert.TagsDatasourceShowOutput(t, "snowflake_tags.test").
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
		},
	})
}

func tagsDatasource(like, resourceName string) string {
	return fmt.Sprintf(`
data "snowflake_tags" "test" {
	depends_on = [%s]

	like = "%s"
}
`, resourceName, like)
}

func TestAcc_Tags_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	id1 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.Tag("test_1", id1.DatabaseName(), id1.Name(), id1.SchemaName())
	model2 := model.Tag("test_2", id2.DatabaseName(), id2.Name(), id2.SchemaName())
	model3 := model.Tag("test_3", id3.DatabaseName(), id3.Name(), id3.SchemaName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testconfig.FromModels(t, model1) + testconfig.FromModels(t, model2) + testconfig.FromModels(t, model3) + tagsDatasourceLike(id1.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tags.test", "tags.#", "1"),
				),
			},
			{
				Config: testconfig.FromModels(t, model1) + testconfig.FromModels(t, model2) + testconfig.FromModels(t, model3) + tagsDatasourceLike(prefix+"%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tags.test", "tags.#", "2"),
				),
			},
		},
	})
}

func tagsDatasourceLike(like string) string {
	return fmt.Sprintf(`
data "snowflake_tags" "test" {
	depends_on = [snowflake_tag.test_1, snowflake_tag.test_2, snowflake_tag.test_3]

	like = "%s"
}
`, like)
}

func TestAcc_Tags_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      tagsDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func tagsDatasourceEmptyIn() string {
	return `
data "snowflake_tags" "test" {
  in {
  }
}
`
}

func TestAcc_Tags_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Tags/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one tag"),
			},
		},
	})
}
