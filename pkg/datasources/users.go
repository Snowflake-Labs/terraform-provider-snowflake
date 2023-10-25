package datasources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
	userPattern := d.Get("pattern").(string)

	account, err := snowflake.ReadCurrentAccount(db)
	if err != nil {
		log.Print("[DEBUG] unable to retrieve current account")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", account.Account, account.Region))

	currentUsers, err := snowflake.ListUsers(userPattern, db)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] no users found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse users in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	users := []map[string]interface{}{}

	for _, user := range currentUsers {
		userMap := map[string]interface{}{}
		userMap["name"] = user.Name.String
		userMap["login_name"] = user.LoginName.String
		userMap["comment"] = user.Comment.String
		userMap["disabled"] = user.Disabled
		userMap["default_warehouse"] = user.DefaultWarehouse.String
		userMap["default_namespace"] = user.DefaultNamespace.String
		userMap["default_role"] = user.DefaultRole.String
		userMap["default_secondary_roles"] = strings.Split(
			helpers.ListContentToString(user.DefaultSecondaryRoles.String), ",")
		userMap["has_rsa_public_key"] = user.HasRsaPublicKey
		userMap["email"] = user.Email.String
		userMap["display_name"] = user.DisplayName.String
		userMap["first_name"] = user.FirstName.String
		userMap["last_name"] = user.LastName.String

		users = append(users, userMap)
	}

	return d.Set("users", users)
}
