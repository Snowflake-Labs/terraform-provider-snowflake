package datasources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var usersSchema = map[string]*schema.Schema{
	"pattern": {
		Type:     schema.TypeString,
		Required: true,
		Description: "Users pattern for which to return metadata. Please refer to LIKE keyword from " +
			"snowflake documentation : https://docs.snowflake.com/en/sql-reference/sql/show-users.html#parameters",
	},
	"users": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The users in the database",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"login_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"disabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"default_warehouse": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"default_namespace": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"default_role": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"default_secondary_roles": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
					Computed: true,
				},
				"has_rsa_public_key": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"email": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"display_name": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"first_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"last_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Users() *schema.Resource {
	return &schema.Resource{
		Read:   ReadUsers,
		Schema: usersSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func ReadUsers(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	userPattern := d.Get("pattern").(string)

	account, err1 := client.ContextFunctions.CurrentAccount(ctx)
	region, err2 := client.ContextFunctions.CurrentRegion(ctx)
	if err1 != nil || err2 != nil {
		log.Print("[DEBUG] unable to retrieve current account")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", account, region))
	extractedUsers, err := client.Users.Show(ctx, &sdk.ShowUserOptions{
		Like: &sdk.Like{Pattern: sdk.String(userPattern)},
	})
	if err != nil {
		log.Printf("[DEBUG] no users found in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	users := make([]map[string]any, len(extractedUsers))

	for i, user := range extractedUsers {
		users[i] = map[string]any{
			"name":                    user.Name,
			"login_name":              user.LoginName,
			"comment":                 user.Comment,
			"disabled":                user.Disabled,
			"default_warehouse":       user.DefaultWarehouse,
			"default_namespace":       user.DefaultNamespace,
			"default_role":            user.DefaultRole,
			"default_secondary_roles": strings.Split(helpers.ListContentToString(user.DefaultSecondaryRoles), ","),
			"has_rsa_public_key":      user.HasRsaPublicKey,
			"email":                   user.Email,
			"display_name":            user.DisplayName,
			"first_name":              user.FirstName,
			"last_name":               user.LastName,
		}
	}

	return d.Set("users", users)
}
