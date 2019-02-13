package resources

import (
	"database/sql"
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
	return CreateResource("database", databaseProperties, databaseSchema, snowflake.Database, ReadDatabase)(data, meta)

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

func UpdateDatabase(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource("database", databaseProperties, databaseSchema, snowflake.Database, ReadDatabase)(data, meta)
}

func DeleteDatabase(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("database", snowflake.Database)(data, meta)
}
