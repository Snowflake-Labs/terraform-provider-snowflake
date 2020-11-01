package resources

import (
	"log"
	"database/sql"

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
		// ForceNew: true,
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

// SqlScript returns a pointer to the resource representing a network policy
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

// CreateSqlScript implements schema.CreateFunc
func CreateSqlScript(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	l := data.Get("lifecycle_commands").([]interface{})
	c := l[0].(map[string]interface{})
	script := c["create"].(string)

	log.Printf("[DEBUG] create sql script: %v", script)
	err := snowflake.Exec(db, script)
	if err != nil {
		return errors.Wrapf(err, "error with create sql script '%v': %v", name, script)
	}

	data.SetId(name)
	return nil
}

// ReadSqlScript implements schema.ReadFunc
func ReadSqlScript(data *schema.ResourceData, meta interface{}) error {
	// db := meta.(*sql.DB)
	// name := data.Get("name").(string)
	// l := data.Get("lifecycle_commands").([]interface{})
	// c := l[0].(map[string]interface{})
	// script := c["read"].(string)

	// log.Printf("[DEBUG] read sql script: %v", script)
	// rows, err := snowflake.Query(db, script)
	// if err != nil {
	// 	return errors.Wrapf(err, "error with read sql script '%v': %v", name, script)
	// }

	// err = data.Set("rows", rows)
	// if err != nil {
	// 	return err
	// }

	return nil	
}

// UpdateSqlScript implements schema.UpdateFunc
func UpdateSqlScript(data *schema.ResourceData, meta interface{}) error {
	// db := meta.(*sql.DB)
	// name := data.Get("name").(string)
	// l := data.Get("lifecycle_commands").([]interface{})
	// c := l[0].(map[string]interface{})
	// script := c["update"].(string)

	// log.Printf("[DEBUG] update sql script: %v", script)
	// err := snowflake.Exec(db, script)
	// if err != nil {
	// 	return errors.Wrapf(err, "error with update sql script '%v': %v", name, script)
	// }

	// err = ReadSqlScript(data, meta)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// DeleteSqlScript implements schema.DeleteFunc
func DeleteSqlScript(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()
	l := data.Get("lifecycle_commands").([]interface{})
	c := l[0].(map[string]interface{})
	script := c["delete"].(string)

	log.Printf("[DEBUG] delete sql script: %v", script)
	err := snowflake.Exec(db, script)
	if err != nil {
		return errors.Wrapf(err, "error with delete sql script '%v': %v", name, script)
	}

	data.SetId("")
	return nil
}

