package resources

import (
	"database/sql"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
)

var warehouseProperties = []string{"comment", "warehouse_size"}
var warehouseSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"comment": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"warehouse_size": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			// TODO
			return
		},
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
	},
}

func Warehouse() *schema.Resource {
	return &schema.Resource{
		Create: CreateWarehouse,
		Read:   ReadWarehouse,
		Delete: DeleteWarehouse,
		Update: UpdateWarehouse,

		Schema: warehouseSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateWarehouse(data *schema.ResourceData, meta interface{}) error {
	return CreateResource("warehouse", warehouseProperties, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(data, meta)

}

func ReadWarehouse(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()

	err := DBExec(db, `USE WAREHOUSE "%s"`, name)
	if err != nil {
		return err
	}

	err = DBExec(db, "SHOW WAREHOUSES LIKE '%s'", name)
	if err != nil {
		return err
	}

	stmt3 := `select "name", "comment", "size" from table(result_scan(last_query_id()))`
	log.Printf("[DEBUG] stmt %s", stmt3)

	row := db.QueryRow(stmt3)

	var warehouseName, comment, size sql.NullString
	err = row.Scan(&warehouseName, &comment, &size)
	if err != nil {
		return err
	}

	err = data.Set("name", warehouseName.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	err = data.Set("warehouse_size", size.String)

	return err
}

func UpdateWarehouse(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource("warehouse", warehouseProperties, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(data, meta)
}

func DeleteWarehouse(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("warehouse", snowflake.Warehouse)(data, meta)
}
