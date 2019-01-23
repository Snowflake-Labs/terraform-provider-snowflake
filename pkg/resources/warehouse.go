package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func Warehouse() *schema.Resource {
	d := NewResourceWarehouse()
	return &schema.Resource{
		Create: d.Create,
		Read:   d.Read,
		Delete: d.Delete,
		Update: d.Update,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateWarehouseName,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				// TODO validation
			},
		},
	}
}

type warehouse struct{}

func NewResourceWarehouse() *warehouse {
	return &warehouse{}
}

func ValidateWarehouseName(val interface{}, key string) ([]string, []error) {
	return snowflake.ValidateIdentifier(val)
}

func (w *warehouse) Create(data *schema.ResourceData, meta interface{}) error {
	name := data.Get("name").(string)
	comment := data.Get("comment").(string)
	db := meta.(*sql.DB)

	err := DBExec(db, "CREATE WAREHOUSE %s COMMENT='%s", name, comment)

	if err != nil {
		return errors.Wrap(err, "error creating warehouse")
	}

	data.SetId(name)

	return nil
}

func (w *warehouse) Read(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	name := data.Id()

	stmt := fmt.Sprintf("SHOW WAREHOUSES LIKE '%s'", name)
	log.Printf("[DEBUG] stmt %s", stmt)

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	stmt2 := `select "name", "comment" from table(result_scan(last_query_id()))`
	log.Printf("[DEBUG] stmt %s", stmt2)

	row2 := db.QueryRow(stmt2)

	var warehouseName, comment sql.NullString
	err = row2.Scan(&warehouseName, &comment)
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
	return nil
}

func (w *warehouse) Delete(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	err := DBExec(db, "DROP WAREHOUSE %s", name)
	if err != nil {
		return errors.Wrapf(err, "error dropping warehouse %s", name)
	}

	return nil
}

func (w *warehouse) Update(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)

	db := meta.(*sql.DB)
	if data.HasChange("name") {
		// I wish this could be done on one line.
		oldNameI, newNameI := data.GetChange("name")
		oldName := oldNameI.(string)
		newName := newNameI.(string)

		err := DBExec(db, "ALTER WAREHOUSE %s RENAME TO %s", oldName, newName)

		if err != nil {
			return errors.Wrapf(err, "error renaming warehouse %s to %s", oldName, newName)
		}
		data.SetId(newName)
		data.SetPartial("name")
	}

	if data.HasChange("comment") {
		name := data.Get("name").(string)
		comment := data.Get("comment").(string)

		err := DBExec(db, "ALTER WAREHOUSE %s SET COMMENT='%s'", name, snowflake.EscapeString(comment))

		if err != nil {
			return errors.Wrap(err, "error altering warehouse")
		}
		data.SetPartial("comment")
	}
	data.Partial(false)
	return nil
}

func DBExec(db *sql.DB, query string, args ...interface{}) error {
	stmt := fmt.Sprintf(query, args...)
	log.Printf("[DEBUG] stmt %s", stmt)

	_, err := db.Exec(stmt)
	return err
}
