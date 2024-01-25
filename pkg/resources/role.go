package resources

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"tag": tagReferenceSchema,
}

func Role() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateAccountRole,
		ReadContext:   ReadAccountRole,
		DeleteContext: DeleteAccountRole,
		UpdateContext: UpdateAccountRole,

		Schema: accountRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	req := sdk.NewCreateRoleRequest(id)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if _, ok := d.GetOk("tag"); ok {
		req.WithTag(getPropertyTags(d, "tag"))
	}

	err := client.Roles.Create(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create account role",
				Detail:   fmt.Sprintf("Account role name: %s, err: %s", name, err),
			},
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadAccountRole(ctx, d, meta)
}

func ReadAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	accountRole, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(id))
	if err != nil {
		if err.Error() == "object does not exist" {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Account role not found; marking it as removed",
					Detail:   fmt.Sprintf("Account role name: %s, err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to show account role by id",
				Detail:   fmt.Sprintf("Account role name: %s, err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	if err := d.Set("name", accountRole.Name); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set account role name",
				Detail:   fmt.Sprintf("Account role name: %s, err: %s", accountRole.Name, err),
			},
		}
	}

	if err := d.Set("comment", accountRole.Comment); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set account role comment",
				Detail:   fmt.Sprintf("Account role name: %s, comment: %s, err: %s", accountRole.Name, accountRole.Comment, err),
			},
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return nil
}

func UpdateAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithSetComment(v.(string)))
			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to set account role comment",
						Detail:   fmt.Sprintf("Account role name: %s, comment: %s, err: %s", id.FullyQualifiedName(), v, err),
					},
				}
			}
		} else {
			err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithUnsetComment(true))
			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to unset account role comment",
						Detail:   fmt.Sprintf("Account role name: %s, err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				tagNames := make([]string, len(unsetTags))
				for i, v := range unsetTags {
					tagNames[i] = v.FullyQualifiedName()
				}
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to unset account role tags",
						Detail:   fmt.Sprintf("Account role name: %s, tags to unset: %v, err: %s", id.FullyQualifiedName(), tagNames, err),
					},
				}
			}
		}

		if len(setTags) > 0 {
			err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithSetTags(setTags))
			if err != nil {
				tagNames := make([]string, len(unsetTags))
				for i, v := range unsetTags {
					tagNames[i] = v.FullyQualifiedName()
				}
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to set account role tags",
						Detail:   fmt.Sprintf("Account role name: %s, tags to set: %v, err: %s", id.FullyQualifiedName(), tagNames, err),
					},
				}
			}
		}
	}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")

		newId, err := helpers.DecodeSnowflakeParameterID(newName.(string))
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to parse account role name",
					Detail:   fmt.Sprintf("Account role name: %s, err: %s", newName, err),
				},
			}
		}

		err = client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithRenameTo(newId.(sdk.AccountObjectIdentifier)))
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to rename account role name",
					Detail:   fmt.Sprintf("Previous account role name: %s, new account role name: %s, err: %s", id, newName, err),
				},
			}
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
	}

	return nil
}

func DeleteAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.Roles.Drop(ctx, sdk.NewDropRoleRequest(id))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to drop account role",
				Detail:   fmt.Sprintf("Account role name: %s, err: %s", d.Id(), err),
			},
		}
	}

	d.SetId("")

	return nil
}
