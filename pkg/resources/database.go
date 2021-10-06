package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var databaseSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: false,
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"data_retention_time_in_days": {
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	},
	"from_share": {
		Type:          schema.TypeMap,
		Description:   "Specify a provider and a share in this map to create a database from a share.",
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"from_database", "from_replica"},
	},
	"from_database": {
		Type:          schema.TypeString,
		Description:   "Specify a database to create a clone from.",
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"from_share", "from_replica"},
	},
	"from_replica": {
		Type:          schema.TypeString,
		Description:   "Specify a fully-qualified path to a database to create a replica from.",
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"from_share", "from_database"},
	},
	"replication_accounts": {
		Type:             schema.TypeSet,
		Elem:             &schema.Schema{Type: schema.TypeString},
		Optional:         true,
		Description:      "A list of accounts to be added to the replication.",
		DiffSuppressFunc: diffCaseInsensitive,
	},
	"replication_failover_accounts": {
		Type:             schema.TypeSet,
		Elem:             &schema.Schema{Type: schema.TypeString},
		Optional:         true,
		Description:      "A list of accounts to be added to the failover replication.",
		DiffSuppressFunc: diffCaseInsensitive,
	},
	"replication_is_primary": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, sets the database as primary for failover.",
		Default:     false,
	},
	"tag": tagReferenceSchema,
}

var databaseProperties = []string{"comment", "data_retention_time_in_days", "tag"}

// Database returns a pointer to the resource representing a database
func Database() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabase,
		Read:   ReadDatabase,
		Delete: DeleteDatabase,
		Update: UpdateDatabase,

		Schema: databaseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateDatabase implements schema.CreateFunc
func CreateDatabase(d *schema.ResourceData, meta interface{}) error {
	if _, ok := d.GetOk("from_share"); ok {
		return createDatabaseFromShare(d, meta)
	}

	if _, ok := d.GetOk("from_database"); ok {
		return createDatabaseFromDatabase(d, meta)
	}

	if _, ok := d.GetOk("from_replica"); ok {
		return createDatabaseFromReplica(d, meta)
	}

	return CreateResource("database", databaseProperties, databaseSchema, snowflake.Database, ReadDatabase)(d, meta)
}

func createDatabaseFromShare(d *schema.ResourceData, meta interface{}) error {
	in := d.Get("from_share").(map[string]interface{})
	prov := in["provider"]
	share := in["share"]

	if prov == nil || share == nil {
		return fmt.Errorf("from_share must contain the keys provider and share, but it had %+v", in)
	}

	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	builder := snowflake.DatabaseFromShare(name, prov.(string), share.(string))

	err := snowflake.Exec(db, builder.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating database %v from share %v.%v", name, prov, share)
	}

	d.SetId(name)

	return ReadDatabase(d, meta)
}

func createDatabaseFromDatabase(d *schema.ResourceData, meta interface{}) error {
	sourceDb := d.Get("from_database").(string)

	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	builder := snowflake.DatabaseFromDatabase(name, sourceDb)

	err := snowflake.Exec(db, builder.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating a clone database %v from database %v", name, sourceDb)
	}

	d.SetId(name)

	return ReadDatabase(d, meta)
}

func createDatabaseFromReplica(d *schema.ResourceData, meta interface{}) error {
	sourceDb := d.Get("from_replica").(string)

	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	builder := snowflake.DatabaseFromReplica(name, sourceDb)

	err := snowflake.Exec(db, builder.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating a secondary database %v from database %v", name, sourceDb)
	}

	d.SetId(name)

	return ReadDatabase(d, meta)
}

func ReadDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()

	stmt := snowflake.Database(name).Show()
	row := snowflake.QueryRow(db, stmt)

	database, err := snowflake.ScanDatabase(row)

	if err != nil {
		if err == sql.ErrNoRows {
			// If not found, mark resource to be removed from statefile during apply or refresh
			log.Printf("[DEBUG] database (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return errors.Wrap(err, "unable to scan row for SHOW DATABASES")
	}

	err = d.Set("name", database.DBName.String)
	if err != nil {
		return err
	}
	err = d.Set("comment", database.Comment.String)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(database.RetentionTime.String, 10, 64)
	if err != nil {
		return err
	}

	err = d.Set("data_retention_time_in_days", i)
	if err != nil {
		return err
	}

	stmt = snowflake.Replication(name).Show()
	row = snowflake.QueryRow(db, stmt)
	replication, err := snowflake.ScanReplication(row)

	if err != nil {
		if err == sql.ErrNoRows {
			// If not found, no replication is enabled
			d.Set("replication_accounts", make([]string, 0))
			d.Set("replication_enable_failover", false)
			d.Set("replication_is_primary", false)
			return nil
		}
		return err
	}

	err = d.Set("replication_accounts", strings.Split(replication.ReplAccounts.String, ", "))
	if err != nil {
		return err
	}

	err = d.Set("replication_failover_accounts", strings.Split(replication.FailoverAccounts.String, ", "))
	if err != nil {
		return err
	}

	return d.Set("replication_is_primary", replication.IsPrimary.Bool)
}

func UpdateDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("name").(string)

	// Change the replication details first - this is a special case and won't work using the generic method
	x := func(resource, replType string, toggleGrant func(db *sql.DB, database, account, replType, toggleType string) error) error {
		old_ver, new_ver := d.GetChange(resource)

		if old_ver == nil {
			old_ver = new(schema.Set)
		}
		if new_ver == nil {
			new_ver = new(schema.Set)
		}
		old_ver_set := old_ver.(*schema.Set)
		new_ver_set := new_ver.(*schema.Set)

		remove := expandStringList(old_ver_set.Difference(new_ver_set).List())
		add := expandStringList(new_ver_set.Difference(old_ver_set).List())

		for _, account := range remove {
			err := toggleReplicationToAcc(db, database, account, replType, "DISABLE")
			if err != nil {
				return err
			}
		}
		for _, account := range add {
			err := toggleReplicationToAcc(db, database, account, replType, "ENABLE")
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := x("replication_accounts", "REPLICATION", toggleReplicationToAcc)
	if err != nil {
		return err
	}

	err = x("replication_failover_accounts", "FAILOVER", toggleReplicationToAcc)
	if err != nil {
		return err
	}

	_, n := d.GetChange("replication_is_primary")

	if n == true {
		err = setPrimaryReplicationSource(db, database)
		if err != nil {
			return err
		}
	}

	return UpdateResource("database", databaseProperties, databaseSchema, snowflake.Database, ReadDatabase)(d, meta)
}

func toggleReplicationToAcc(db *sql.DB, database, account, replType, toggleType string) error {
	err := snowflake.Exec(db, fmt.Sprintf(`ALTER DATABASE "%v" %v %v TO ACCOUNTS "%v"`, database, toggleType, replType, account))
	return err
}

func setPrimaryReplicationSource(db *sql.DB, database string) error {
	err := snowflake.Exec(db, fmt.Sprintf(`ALTER DATABASE "%v" PRIMARY`, database))
	return err
}

func DeleteDatabase(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("database", snowflake.Database)(d, meta)
}
