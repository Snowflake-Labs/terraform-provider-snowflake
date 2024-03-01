package datasources

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of the [primary role](https://docs.snowflake.com/en/user-guide/security-access-control-overview.html#label-access-control-role-enforcement) in use for the current session.",
	},
}

func CurrentRole() *schema.Resource {
	return &schema.Resource{
		Read:   ReadCurrentRole,
		Schema: currentRoleSchema,
	}
}

func ReadCurrentRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	role, err := client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		log.Printf("[DEBUG] current_role failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(role)
	err = d.Set("name", role)
	if err != nil {
		return err
	}
	return nil
}
