package datasources

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	RoleNameKey    = "name"
	RoleCommentKey = "comment"
)

var roleSchema = map[string]*schema.Schema{
	RoleNameKey: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The role for which to return metadata.",
	},
	RoleCommentKey: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The comment on the role.",
	},
}

// Role Snowflake Role resource.
func Role() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRole,
		Schema: roleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadRole Reads the database metadata information.
func ReadRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	role, err := client.Roles.ShowByID(ctx, objectIdentifier)
	if err != nil {
		return err
	}

	d.SetId(role.Name)

	return errors.Join(
		d.Set(RoleNameKey, role.Name),
		d.Set(RoleCommentKey, role.Comment),
	)
}
