package resources

import (
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func isNotOwnershipGrant() func(value any, path cty.Path) diag.Diagnostics {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if privilege, ok := value.(string); ok && strings.ToUpper(privilege) == "OWNERSHIP" {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Unsupported privilege 'OWNERSHIP'",
				Detail:        "Granting ownership is only allowed in snowflake_grant_ownership resource.",
				AttributePath: nil,
			})
		}
		return diags
	}
}
