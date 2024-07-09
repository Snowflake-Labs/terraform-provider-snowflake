package assert

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type WarehouseAssert struct {
	*SnowflakeObjectAssert[sdk.Warehouse, sdk.AccountObjectIdentifier]
}

func Warehouse(t *testing.T, id sdk.AccountObjectIdentifier) *WarehouseAssert {
	t.Helper()
	return &WarehouseAssert{
		NewSnowflakeObjectAssertWithProvider(sdk.ObjectTypeWarehouse, id, acc.TestClient().Warehouse.Show),
	}
}

func WarehouseFromObject(t *testing.T, warehouse *sdk.Warehouse) *WarehouseAssert {
	t.Helper()
	return &WarehouseAssert{
		NewSnowflakeObjectAssertWithObject(sdk.ObjectTypeWarehouse, warehouse.ID(), warehouse),
	}
}

func (w *WarehouseAssert) HasName(expected string) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasState(expected sdk.WarehouseState) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.State != expected {
			return fmt.Errorf("expected state: %v; got: %v", expected, string(o.State))
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasType(expected sdk.WarehouseType) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.Type != expected {
			return fmt.Errorf("expected type: %v; got: %v", expected, string(o.Type))
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasSize(expected sdk.WarehouseSize) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.Size != expected {
			return fmt.Errorf("expected size: %v; got: %v", expected, string(o.Size))
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasMinClusterCount(expected int) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.MinClusterCount != expected {
			return fmt.Errorf("expected min cluster count: %v; got: %v", expected, o.MinClusterCount)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasMaxClusterCount(expected int) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.MaxClusterCount != expected {
			return fmt.Errorf("expected max cluster count: %v; got: %v", expected, o.MaxClusterCount)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasAutoSuspend(expected int) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.AutoSuspend != expected {
			return fmt.Errorf("expected auto suspend: %v; got: %v", expected, o.AutoSuspend)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasAutoResume(expected bool) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.AutoResume != expected {
			return fmt.Errorf("expected auto resume: %v; got: %v", expected, o.AutoResume)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasComment(expected string) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasEnableQueryAcceleration(expected bool) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.EnableQueryAcceleration != expected {
			return fmt.Errorf("expected enable query acceleration: %v; got: %v", expected, o.EnableQueryAcceleration)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasQueryAccelerationMaxScaleFactor(expected int) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.QueryAccelerationMaxScaleFactor != expected {
			return fmt.Errorf("expected query acceleration max scale factor: %v; got: %v", expected, o.QueryAccelerationMaxScaleFactor)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasResourceMonitor(expected sdk.AccountObjectIdentifier) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.ResourceMonitor.Name() != expected.Name() {
			return fmt.Errorf("expected resource monitor: %v; got: %v", expected.Name(), o.ResourceMonitor.Name())
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasScalingPolicy(expected sdk.ScalingPolicy) *WarehouseAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.ScalingPolicy != expected {
			return fmt.Errorf("expected type: %v; got: %v", expected, string(o.ScalingPolicy))
		}
		return nil
	})
	return w
}
