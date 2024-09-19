package resources

import (
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	config "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

// The following tests are showing the behavior of the provider in cases where objects higher in the hierarchy
// like database or schema are renamed. Learn more on that in the document (TODO(SNOW-): link).
//
// Shallow hierarchy (database + schema)
// - is in config - renamed internally
// - is in config - renamed externally
// - is not in config - renamed internally
// - is not in config - renamed externally
//
// Deep hierarchy (database + schema + schema object)
// - only database is in config - renamed internally
// - only database is in config - renamed externally
// - only schema is in config - renamed internally
// - only schema is in config - renamed externally
// - both database and schema are in config - renamed internally
// - both database and schema are in config - renamed externally
// - both database and schema are not in config - renamed internally
// - both database and schema are not in config - renamed externally

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(id)
	databaseConfigModel := model.Database("test", id.Name())
	schemaConfigModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	schemaConfigModel.SetDependsOn([]string{})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseConfigModel) + config.FromModel(t, schemaConfigModel),
			},
		},
	})
}
