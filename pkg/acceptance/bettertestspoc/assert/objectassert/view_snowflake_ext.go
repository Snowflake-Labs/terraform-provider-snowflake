package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (w *ViewAssert) HasCreatedOnNotEmpty() *ViewAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return w
}

func (v *ViewAssert) HasNonEmptyText() *ViewAssert {
	v.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		if o.Text == "" {
			return fmt.Errorf("expected non empty text")
		}
		return nil
	})
	return v
}
