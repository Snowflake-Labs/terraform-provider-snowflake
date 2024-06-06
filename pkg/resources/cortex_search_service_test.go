package resources_test

import (
	"context"
	"fmt"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("snowflake_cortex_search_service", &resource.Sweeper{
		Name: "snowflake_cortex_search_service",
		F: func(profile string) error {
			client, err := sdk.NewDefaultClient()
			if err != nil {
				return fmt.Errorf("error getting default client during sweep: %w", err)
			}
			ctx := context.Background()
			cortexSearchServices, err := client.CortexSearchServices.Show(ctx, sdk.NewShowCortexSearchServiceRequest().WithIn(&sdk.In{
				Schema: acc.TestClient().Ids.SchemaId(),
			}))
			if err != nil {
				return fmt.Errorf("error getting cortex search services during sweep: %w", err)
			}
			for _, cortexSearchService := range cortexSearchServices {
				err := client.CortexSearchServices.Drop(ctx, sdk.NewDropCortexSearchServiceRequest(cortexSearchService.ID()))
				if err != nil {
					return fmt.Errorf("error dropping cortex search service %s during sweep: %w", cortexSearchService.ID(), err)
				}
			}
			return nil
		},
	})
}
