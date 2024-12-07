package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *FunctionAssert) HasCreatedOnNotEmpty() *FunctionAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Function) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected create_on to be not empty")
		}
		return nil
	})
	return a
}
