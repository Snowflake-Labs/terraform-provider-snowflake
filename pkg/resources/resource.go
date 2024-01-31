package resources

import (
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DeleteResource(t string, builder func(string) *snowflake.Builder) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		db := meta.(*sql.DB)
		name := d.Get("name").(string)

		stmt := builder(name).Drop()
		if err := snowflake.Exec(db, stmt); err != nil {
			return fmt.Errorf("error dropping %s %s err = %w", t, name, err)
		}

		d.SetId("")
		return nil
	}
}
