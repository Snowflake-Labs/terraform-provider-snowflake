package bettertestspoc

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (w *WarehouseAssert) HasStateOneOf(expected ...sdk.WarehouseState) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if !slices.Contains(expected, o.State) {
			return fmt.Errorf("expected state one of: %v; got: %v", expected, string(o.State))
		}
		return nil
	})
	return w
}
