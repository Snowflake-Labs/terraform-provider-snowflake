package resources

import (
	"database/sql"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var sqlScriptSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the sql script.",
		ForceNew:    true,
	},
	"lifecycle_commands": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"create": {
					Type:     schema.TypeString,
					Required: true,
				},
				"update": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"read": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"delete": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
}

func SqlScript() *schema.Resource {
	return &schema.Resource{
		Create: CreateSqlScript,
		Read:   ReadSqlScript,
		Update: UpdateSqlScript,
		Delete: DeleteSqlScript,

		Schema: sqlScriptSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateSqlScript(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	l := data.Get("lifecycle_commands").([]interface{})
	c := l[0].(map[string]interface{})
	script := c["create"].(string)

	err := snowflake.Exec(db, script)
	if err != nil {
		return errors.Wrapf(err, "error with create sql script '%v': %v", name, script)
	}

	data.SetId(name)
	return ReadSqlScript(data, meta)
}

func ReadSqlScript(data *schema.ResourceData, meta interface{}) error {
	l := data.Get("lifecycle_commands").([]interface{})
	c := l[0].(map[string]interface{})
	script := c["read"].(string)

	if len(script) != 0 {
		log.Printf("[WARN] snowflake_sql read is not implemented and does nothing.")
	}

	return nil
}

func UpdateSqlScript(data *schema.ResourceData, meta interface{}) error {
	l := data.Get("lifecycle_commands").([]interface{})
	c := l[0].(map[string]interface{})
	script := c["update"].(string)

	if len(script) != 0 {
		log.Printf("[WARN] snowflake_sql update is not implemented and does nothing.")
	}

	return nil
}

func DeleteSqlScript(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()
	l := data.Get("lifecycle_commands").([]interface{})
	c := l[0].(map[string]interface{})
	script := c["delete"].(string)

	err := snowflake.Exec(db, script)
	if err != nil {
		return errors.Wrapf(err, "error with delete sql script '%v': %v", name, script)
	}

	data.SetId("")
	return nil
}
