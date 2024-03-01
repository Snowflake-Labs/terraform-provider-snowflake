package datasources

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountRolesSchema = map[string]*schema.Schema{
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the command output by object name.",
	},
	"roles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the roles which you can view across your entire account, including the system-defined roles and any custom roles that exist.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Identifier for the role.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The comment on the role",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The owner of the role",
				},
			},
		},
	},
}

func Roles() *schema.Resource {
	return &schema.Resource{
		ReadContext: ReadAccountRoles,
		Schema:      accountRolesSchema,
	}
}

func ReadAccountRoles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	req := sdk.NewShowRoleRequest()
	if pattern, ok := d.GetOk("pattern"); ok {
		req.WithLike(sdk.NewLikeRequest(pattern.(string)))
	}

	roles, err := client.Roles.Show(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to show account roles",
				Detail:   fmt.Sprintf("Search pattern: %v, err: %s", d.Get("pattern").(string), err),
			},
		}
	}

	mappedRoles := make([]map[string]any, len(roles))
	for i, role := range roles {
		mappedRoles[i] = map[string]any{
			"name":    role.Name,
			"comment": role.Comment,
			"owner":   role.Owner,
		}
	}

	if err := d.Set("roles", mappedRoles); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set roles",
				Detail:   fmt.Sprintf("Search pattern: %v, err: %s", d.Get("pattern").(string), err),
			},
		}
	}

	d.SetId("roles_read")

	return nil
}
