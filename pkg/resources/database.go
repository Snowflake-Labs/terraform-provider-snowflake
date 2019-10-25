package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	"from_share": &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Specify a provider and a share in this map to create a database from a share.",
		Optional:    true,
		ForceNew:    true,
		ConflictsWith: []string{"from_database"},
	},
	"from_database": &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a database to create a clone from.",
		Optional:    true,
		ForceNew:    true,
		ConflictsWith: []string{"from_share"},
	},
}

var databaseProperties = []string{"comment", "data_retention_time_in_days"}

// Database returns a pointer to the resource representing a database
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

// CreateDatabase implements schema.CreateFunc
func CreateDatabase(data *schema.ResourceData, meta interface{}) error {
	if _, ok := data.GetOk("from_share"); ok {
		return createDatabaseFromShare(data, meta)
	}
	
	if _, ok := data.GetOk("from_database"); ok {
		return createDatabaseFromDatabase(data, meta)
	}

	return CreateResource("database", databaseProperties, databaseSchema, snowflake.Database, ReadDatabase)(data, meta)
}

func createDatabaseFromShare(data *schema.ResourceData, meta interface{}) error {
	in := data.Get("from_share").(map[string]interface{})
	prov := in["provider"]
	share := in["share"]

	if prov == nil || share == nil {
		return fmt.Errorf("from_share must contain the keys provider and share, but it had %+v", in)
	}

	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	builder := snowflake.DatabaseFromShare(name, prov.(string), share.(string))

	err := DBExec(db, builder.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating database %v from share %v.%v", name, prov, share)
	}

	data.SetId(name)

	return ReadDatabase(data, meta)
}

func createDatabaseFromDatabase(data *schema.ResourceData, meta interface{}) error {
	sourceDb := data.Get("from_database").(string)

	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	builder := snowflake.DatabaseFromDatabase(name, sourceDb)

	err := DBExec(db, builder.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating a clone database %v from database %v", name, sourceDb)
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
		if err == sql.ErrNoRows {
			log.Printf("[WARN] database %v not found, removing from state file", name)
			data.SetId("")
			return nil
		}
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
