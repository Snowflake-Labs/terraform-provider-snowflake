package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var rolesSchema = map[string]*schema.Schema{
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"in_class": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: resources.IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		Description:      "Filters the SHOW GRANTS output by class name.",
	},
	"roles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all role details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW ROLES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowRoleSchema,
					},
				},
			},
		},
	},
}

func Roles() *schema.Resource {
	return &schema.Resource{
		ReadContext: ReadRoles,
		Schema:      rolesSchema,
		Description: "Datasource used to get details of filtered roles. Filtering is aligned with the current possibilities for [SHOW ROLES](https://docs.snowflake.com/en/sql-reference/sql/show-roles) query (`like` and `in_class` are all supported). The results of SHOW are encapsulated in one output collection.",
	}
}

func ReadRoles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	req := sdk.NewShowRoleRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		req.WithLike(sdk.NewLikeRequest(likePattern.(string)))
	}

	if className, ok := d.GetOk("in_class"); ok {
		req.WithInClass(sdk.RolesInClass{
			Class: sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(className.(string)),
		})
	}

	roles, err := client.Roles.Show(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to show roles",
				Detail:   fmt.Sprintf("Error: %s", err),
			},
		}
	}

	d.SetId("roles_read")

	flattenedRoles := make([]map[string]any, len(roles))
	for i, role := range roles {
		role := role
		flattenedRoles[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.RoleToSchema(&role)},
		}
	}

	err = d.Set("roles", flattenedRoles)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
