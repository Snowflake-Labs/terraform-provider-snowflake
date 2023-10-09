package resources_test

import (
	"context"
	"fmt"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("snowflake_dynamic_table", &resource.Sweeper{
		Name: "snowflake_dynamic_table",
		F: func(profile string) error {
			client, err := sdk.NewDefaultClient()
			if err != nil {
				return fmt.Errorf("Error getting default client during sweep: %s", err)
			}
			ctx := context.Background()
			dynamicTables, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithIn(&sdk.In{
				Schema: sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName),
			}))
			if err != nil {
				return fmt.Errorf("Error getting dynamic tables during sweep: %s", err)
			}
			for _, dynamicTable := range dynamicTables {
				err := client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(dynamicTable.ID()))
				if err != nil {
					return fmt.Errorf("Error dropping dynamic table %s during sweep: %s", dynamicTable.ID(), err)
				}
			}
			return nil
		},
	})
}
