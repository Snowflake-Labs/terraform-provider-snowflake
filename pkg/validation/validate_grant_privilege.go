package validation

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ValidateGrantPrivilege(valid []string) schema.SchemaValidateDiagFunc {
	f := func(i interface{}, k cty.Path) error {
		v, ok := i.(string)
		if !ok {
			return fmt.Errorf("expected type of %s to be string", k.Index())
		}
	}

	return func(i interface{}, k cty.Path) diag.Diagnostics {
		err := f(i, k)
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

}
