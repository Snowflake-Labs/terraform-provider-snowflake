package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jmoiron/sqlx"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"data_retention_time_in_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateDatabase(data *schema.ResourceData, meta interface{}) error {
	name := data.Get("name").(string)
	comment := data.Get("comment").(string)
	retention, retentionSet := data.GetOk("data_retention_time_in_days")
	db := meta.(*sql.DB)

	stmt := fmt.Sprintf(`CREATE DATABASE "%s" COMMENT='%s'`, name, snowflake.EscapeString(comment))
	if retentionSet {
		stmt = fmt.Sprintf("%s DATA_RETENTION_TIME_IN_DAYS = %d", stmt, retention)
	}
	log.Printf("[DEBUG] stmt %s", stmt)
	_, err := db.Exec(stmt)

	if err != nil {
		return errors.Wrap(err, "error creating database")
	}

	data.SetId(name)

	return ReadDatabase(data, meta)
}

type database struct {
	CreatedOn     sql.NullString `db:"created_on"`
	DBName        sql.NullString `db:"name"`
	IsDefault     sql.NullString `db:"is_default"`
	IsCurrent     sql.NullString `db:"is_current"`
	Origin        sql.NullString `db:"origin"`
	Owner         sql.NullString `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime sql.NullString `db:"retentionTime"`
}

func ReadDatabase(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sdb := sqlx.NewDb(db, "snowflake")

	// TODO Not sure if we should use id or name here.
	name := data.Id()

	// TODO make sure there are no wildcard-y characters here, otherwise it could match more than 1 row.
	stmt := fmt.Sprintf("SHOW DATABASES LIKE '%s'", name)

	log.Printf("[DEBUG] stmt %s", stmt)
	row := sdb.QueryRowx(stmt)

	database := &database{}
	err := row.StructScan(database)

	if err != nil {
		return errors.Wrap(err, "unable to scan row for SHOW DATABASES")
	}

	err = data.Set("name", database.DBName.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", database.Comment.String)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(database.RetentionTime.String, 10, 64)
	if err != nil {
		return err
	}

	err = data.Set("data_retention_time_in_days", i)
	return err
}

func DeleteDatabase(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	stmt := fmt.Sprintf(`DROP DATABASE "%s"`, name)
	log.Printf("[DEBUG] stmt %s", stmt)
	_, err := db.Exec(stmt)
	if err != nil {
		return errors.Wrapf(err, "error dropping database %s", name)
	}

	data.SetId("")
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

		stmt := fmt.Sprintf(`ALTER DATABASE "%s" RENAME TO "%s"`, oldName, newName)
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

		stmt := fmt.Sprintf(`ALTER DATABASE "%s" SET COMMENT='%s'`, name, snowflake.EscapeString(comment))
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

		stmt := fmt.Sprintf(`ALTER DATABASE "%s" SET DATA_RETENTION_TIME_IN_DAYS = %d`, name, retention)
		log.Printf("[DEBUG] stmt %s", stmt)

		_, err := db.Exec(stmt)
		if err != nil {
			return errors.Wrap(err, "Error setting data_retention_time_in_days")
		}
		data.SetPartial("data_retention_time_in_days")
	}
	data.Partial(false)
	return ReadDatabase(data, meta)
}
