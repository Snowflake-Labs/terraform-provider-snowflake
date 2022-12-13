package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
		Description: "The Snowflake Primary Role ID; as returned by CURRENT_ROLE()."
	},
}

func CurrentRole() *schema.Resource {
	return &schema.Resource{
		Read:   ReadCurrentRole,
		Schema: currentRoleSchema,
	}
}

func ReadCurrentRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	role, err := snowflake.ReadCurrentRole(db)

	if err != nil {
		log.Printf("[DEBUG] current_role failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf(role.Role))
	d.Set("name", role.Role)

	return nil
}
