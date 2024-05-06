package datasources_test

import (
	"fmt"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Databases(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databases(databaseName, comment),
				Check: resource.ComposeTestCheckFunc(
					checkDatabases(databaseName, comment),
				),
			},
		},
	})
}

func databases(databaseName, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_database" {
			name = "%v"
			comment = "%v"
		}
		data snowflake_databases "t" {
			depends_on = [snowflake_database.test_database]
		}
	`, databaseName, comment)
}

func checkDatabases(databaseName string, comment string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.Modules[0].Resources["data.snowflake_databases.t"]
		if resourceState == nil {
			return fmt.Errorf("resource not found in state")
		}
		instanceState := resourceState.Primary
		if instanceState == nil {
			return fmt.Errorf("resource has no primary instance")
		}
		if instanceState.ID != "databases_read" {
			return fmt.Errorf("expected ID to be 'databases_read', got %s", instanceState.ID)
		}
		nDbs, err := strconv.Atoi(instanceState.Attributes["databases.#"])
		if err != nil {
			return fmt.Errorf("expected a number for field 'databases', got %s", instanceState.Attributes["databases.#"])
		}
		if nDbs == 0 {
			return fmt.Errorf("expected databases to be greater or equal to 1, got %s", instanceState.Attributes["databases.#"])
		}
		dbIdx := -1
		for i := 0; i < nDbs; i++ {
			idxName := fmt.Sprintf("databases.%d.name", i)
			if instanceState.Attributes[idxName] == databaseName {
				dbIdx = i
				break
			}
		}
		if dbIdx == -1 {
			return fmt.Errorf("database %s not found", databaseName)
		}
		idxComment := fmt.Sprintf("databases.%d.comment", dbIdx)
		if instanceState.Attributes[idxComment] != comment {
			return fmt.Errorf("expected comment '%s', got '%s'", comment, instanceState.Attributes[idxComment])
		}
		idxCreatedOn := fmt.Sprintf("databases.%d.created_on", dbIdx)
		if instanceState.Attributes[idxCreatedOn] == "" {
			return fmt.Errorf("expected 'created_on' to be set")
		}
		idxOwner := fmt.Sprintf("databases.%d.owner", dbIdx)
		if instanceState.Attributes[idxOwner] == "" {
			return fmt.Errorf("expected 'owner' to be set")
		}
		idxRetentionTime := fmt.Sprintf("databases.%d.retention_time", dbIdx)
		if instanceState.Attributes[idxRetentionTime] == "" {
			return fmt.Errorf("expected 'retention_time' to be set")
		}
		idxIsCurrent := fmt.Sprintf("databases.%d.is_current", dbIdx)
		if instanceState.Attributes[idxIsCurrent] == "" {
			return fmt.Errorf("expected 'is_current' to be set")
		}
		idxIsDefault := fmt.Sprintf("databases.%d.is_default", dbIdx)
		if instanceState.Attributes[idxIsDefault] == "" {
			return fmt.Errorf("expected 'is_default' to be set")
		}
		return nil
	}
}
