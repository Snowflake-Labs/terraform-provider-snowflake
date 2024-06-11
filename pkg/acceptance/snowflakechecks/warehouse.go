package snowflakechecks

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1473412]: consider using size from state instead of passing it
func CheckWarehouseSize(t *testing.T, id sdk.AccountObjectIdentifier, expectedSize sdk.WarehouseSize) func(state *terraform.State) error {
	t.Helper()
	return func(_ *terraform.State) error {
		warehouse, err := acc.TestClient().Warehouse.Show(t, id)
		if err != nil {
			return err
		}
		if warehouse.Size != expectedSize {
			return fmt.Errorf("expected size: %s; got: %s", expectedSize, warehouse.Size)
		}
		return nil
	}
}

func CheckWarehouseType(t *testing.T, id sdk.AccountObjectIdentifier, expectedType sdk.WarehouseType) func(state *terraform.State) error {
	t.Helper()
	return func(_ *terraform.State) error {
		warehouse, err := acc.TestClient().Warehouse.Show(t, id)
		if err != nil {
			return err
		}
		if warehouse.Type != expectedType {
			return fmt.Errorf("expected type: %s; got: %s", expectedType, warehouse.Type)
		}
		return nil
	}
}
