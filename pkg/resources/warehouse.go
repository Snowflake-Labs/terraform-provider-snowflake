package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

var properties = []string{"comment", "warehouse_size"}

func Warehouse() *schema.Resource {
	d := NewResourceWarehouse()
	return &schema.Resource{
		Create: d.Create,
		Read:   d.Read,
		Delete: d.Delete,
		Update: d.Update,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					return snowflake.ValidateIdentifier(val)
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToUpper(old) == strings.ToUpper(new)
				},
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				// TODO validation
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
		},
	}
}

type warehouse struct{}

func NewResourceWarehouse() *warehouse {
	return &warehouse{}
}

func (w *warehouse) Create(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	var sb strings.Builder

	_, err := sb.WriteString(fmt.Sprintf("CREATE WAREHOUSE %s", name))
	if err != nil {
		return err
	}

	for _, field := range properties {
		val, ok := data.GetOk(field)
		valStr := val.(string)
		if ok {
			_, e := sb.WriteString(fmt.Sprintf(" %s='%s'", strings.ToUpper(field), snowflake.EscapeString(valStr)))
			if e != nil {
				return e
			}
		}
	}
	err = DBExec(db, sb.String())


	if err != nil {
		return errors.Wrap(err, "error creating warehouse")
	}

	data.SetId(name)

	return nil
}

func (w *warehouse) Read(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Id()

	err := DBExec(db, `USE WAREHOUSE %s`, name)
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

	db := meta.(*sql.DB)
	if data.HasChange("name") {
		data.Partial(true)
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
		data.Partial(false)
	}

	changes := []string{}

	for _, prop := range properties {
		if data.HasChange(prop) {
			changes = append(changes, prop)
		}
	}
	if len(changes) > 0 {
		name := data.Get("name").(string)
		var sb strings.Builder
		_, err := sb.WriteString(fmt.Sprintf("ALTER WAREHOUSE %s SET", name))
		if err != nil {
			return err
		}

		for _, change := range changes {
			val := data.Get(change).(string)
			_, err := sb.WriteString(fmt.Sprintf(" %s='%s'",
				strings.ToUpper(change), snowflake.EscapeString(val)))
			if err != nil {
				return err
			}
		}

		err = DBExec(db, sb.String())
		if err != nil {
			return errors.Wrap(err, "error altering warehouse")
		}
	}
	return nil
}

func DBExec(db *sql.DB, query string, args ...interface{}) error {
	stmt := fmt.Sprintf(query, args...)
	log.Printf("[DEBUG] stmt %s", stmt)

	_, err := db.Exec(stmt)
	return err
}
