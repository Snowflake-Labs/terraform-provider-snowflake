package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func Database() *schema.Resource {
	d := newResourceDatabase()
	return &schema.Resource{
		Create: d.Create,
		Read:   d.Read,
		Delete: d.Delete,
		Update: d.Update,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				Description:  "TODO",
				ValidateFunc: ValidateDatabaseName,
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

type database struct{}

func newResourceDatabase() *database {
	return &database{}
}

func ValidateDatabaseName(val interface{}, key string) ([]string, []error) {
	return snowflake.ValidateIdentifier(val)

}

func (d *database) Create(data *schema.ResourceData, meta interface{}) error {
	name := data.Get("name").(string)
	comment := data.Get("comment").(string)
	db := meta.(*sql.DB)

	// TODO prepared statements don't appear to work for DDL statements, so we might need to do all this ourselves
	// TODO escape name
	// TODO escape comment
	// TODO name appears to get normalized to uppercase, should we do that? or maybe just consider it
	// 	case-insensitive?
	stmt := fmt.Sprintf("CREATE DATABASE %s COMMENT='%s'", name, snowflake.EscapeString(comment))
	log.Printf("[DEBUG] stmt %s", stmt)
	_, err := db.Exec(stmt)

	if err != nil {
		return errors.Wrap(err, "error creating database")
	}

	data.SetId(name)

	return d.Read(data, meta)
}

func (d *database) Read(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	// TODO Not sure if we should use id or name here.
	name := data.Id()

	// TODO make sure there are no wildcard-y characters here, otherwise it could match more than1 row.
	stmt := fmt.Sprintf("SHOW DATABASES LIKE '%s'", name)
	log.Printf("[DEBUG] stmt %s", stmt)

	// TODO if we try to read a row and there is none, this will return an error. We should probably
	//      do something more graceful
	row := db.QueryRow(stmt)

	var createdOn, dbname, isDefault, isCurrent, origin, owner, comment, options, retentionTime sql.NullString

	err := row.Scan(
		&createdOn, &dbname, &isDefault, &isCurrent, &origin, &owner, &comment, &options, &retentionTime,
	)

	if err != nil {
		return errors.Wrap(err, "unable to scan row for SHOW DATABASES")
	}

	data.Set("name", dbname)
	return nil
}

func (d *database) Delete(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	stmt := fmt.Sprintf("DROP DATABASE %s", name)
	log.Printf("[DEBUG] stmt %s", stmt)
	_, err := db.Exec(stmt)
	if err != nil {
		return errors.Wrapf(err, "error dropping database %s", name)
	}

	return nil
}

func (d *database) Update(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	db := meta.(*sql.DB)
	if data.HasChange("name") {
		// I wish this could be done on one line.
		oldNameI, newNameI := data.GetChange("name")
		oldName := oldNameI.(string)
		newName := newNameI.(string)

		stmt := fmt.Sprintf("ALTER DATABASE %s RENAME TO %s", oldName, newName)
		log.Printf("[DEBUG] stmt %s", stmt)

		_, err := db.Exec(stmt)
		if err != nil {
			return errors.Wrapf(err, "error renaming database %s to %s", oldName, newName)
		}
		data.SetId(newName)
		data.SetPartial("name")
	}

	if data.HasChange("comment") {
		name := data.Get("name").(string)
		comment := data.Get("comment").(string)

		stmt := fmt.Sprintf("ALTER DATABASE %s SET COMMENT='%s'", name, snowflake.EscapeString(comment))
		log.Printf("[DEBUG] stmt %s", stmt)

		_, err := db.Exec(stmt)
		if err != nil {
			return errors.Wrap(err, "error altering database")
		}
		data.SetPartial("comment")
	}
	data.Partial(false)
	return nil
}
