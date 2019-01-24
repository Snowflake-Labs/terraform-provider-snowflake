package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func Database() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabase,
		Read:   ReadDatabase,
		Delete: DeleteDatabase,
		Update: UpdateDatabase,

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
			"data_retention_time_in_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func ValidateDatabaseName(val interface{}, key string) ([]string, []error) {
	return snowflake.ValidateIdentifier(val)
}

func CreateDatabase(data *schema.ResourceData, meta interface{}) error {
	name := data.Get("name").(string)
	comment := data.Get("comment").(string)
	retention := data.Get("data_retention_time_in_days")
	db := meta.(*sql.DB)

	// TODO prepared statements don't appear to work for DDL statements, so we might need to do all this ourselves
	// TODO name appears to get normalized to uppercase, should we do that? or maybe just consider it
	// 	case-insensitive?
	stmt := fmt.Sprintf("CREATE DATABASE %s COMMENT='%s'", name, snowflake.EscapeString(comment))
	if retention != nil {
		stmt = fmt.Sprintf("%s DATA_RETENTION_TIME_IN_DAYS = %d", stmt, retention)
	}
	log.Printf("[DEBUG] stmt %s", stmt)
	_, err := db.Exec(stmt)

	if err != nil {
		return errors.Wrap(err, "error creating database")
	}

	data.SetId(name)

	// return ReadDatabase(data, meta)
	return nil
}

func ReadDatabase(data *schema.ResourceData, meta interface{}) error {
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

	err = data.Set("name", dbname.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(retentionTime.String, 10, 64)
	if err != nil {
		return err
	}

	err = data.Set("data_retention_time_in_days", i)
	return err
}

func DeleteDatabase(data *schema.ResourceData, meta interface{}) error {
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

func UpdateDatabase(data *schema.ResourceData, meta interface{}) error {
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

	// TODO collapse these two conditionals into a loop that generates a single statement.
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

	if data.HasChange("data_retention_time_in_days") {
		name := data.Get("name").(string)
		retention := data.Get("data_retention_time_in_days").(int)

		stmt := fmt.Sprintf("ALTER DATABASE %s SET DATA_RETENTION_TIME_IN_DAYS = %d", name, retention)
		log.Printf("[DEBUG] stmt %s", stmt)

		_, err := db.Exec(stmt)
		if err != nil {
			return errors.Wrap(err, "Error setting data_retention_time_in_days")
		}
		data.SetPartial("data_retention_time_in_days")
	}
	data.Partial(false)
	return nil
}
