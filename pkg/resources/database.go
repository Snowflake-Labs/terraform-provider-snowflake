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

var databaseSchema = map[string]*schema.Schema{
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
}

var databaseProperties = []string{"comment", "data_retention_time_in_days"}

func Database() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabase,
		Read:   ReadDatabase,
		Delete: DeleteDatabase,
		Update: UpdateDatabase,

		Schema: databaseSchema,
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
	RetentionTime sql.NullString `db:"retention_time"`
}

func ReadDatabase(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sdb := sqlx.NewDb(db, "snowflake")

	name := data.Id()

	stmt := snowflake.Database(name).Show()

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

	stmt := snowflake.Database(name).Drop()

	err := DBExec(db, stmt)
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

		stmt := snowflake.Database(oldName).Rename(newName)

		err := DBExec(db, stmt)
		if err != nil {
			return errors.Wrapf(err, "error renaming database %s to %s", oldName, newName)
		}

		data.SetId(newName)
		data.SetPartial("name")
	}
	data.Partial(false)

	// TODO this was c/p from user.go we can probably refactor to a common implementation
	changes := []string{}

	for _, prop := range databaseProperties {
		if data.HasChange(prop) {
			changes = append(changes, prop)
		}
	}
	if len(changes) > 0 {
		name := data.Get("name").(string)
		qb := snowflake.Database(name).Alter()

		for _, field := range changes {
			val := data.Get(field)
			switch databaseSchema[field].Type {
			case schema.TypeString:
				valStr := val.(string)
				qb.SetString(field, valStr)
			case schema.TypeBool:
				valBool := val.(bool)
				qb.SetBool(field, valBool)
			case schema.TypeInt:
				valInt := val.(int)
				qb.SetInt(field, valInt)
			}
		}

		err := DBExec(db, qb.Statement())
		if err != nil {
			return errors.Wrap(err, "error altering database")
		}
	}
	return ReadDatabase(data, meta)
}
