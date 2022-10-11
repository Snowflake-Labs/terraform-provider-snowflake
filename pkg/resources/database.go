package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
	"is_transient": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies a database as transient. Transient databases do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
		ForceNew:    true,
	},
	"data_retention_time_in_days": {
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	},
	"from_share": {
		Type:          schema.TypeMap,
		Elem:          &schema.Schema{Type: schema.TypeString},
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
		Description:   "Specify a fully-qualified path to a database to create a replica from. A fully qualified path follows the format of \"<organization_name>\".\"<account_name>\".\"<db_name>\". An example would be: \"myorg1\".\"account1\".\"db1\"",
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"from_share", "from_database"},
	},
	"replication_configuration": {
		Type:        schema.TypeList,
		Description: "When set, specifies the configurations for database replication.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"accounts": {
					Type:     schema.TypeList,
					Required: true,
					MinItems: 1,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"ignore_edition_check": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
			},
		},
	},
	"tag": tagReferenceSchema,
}

// Database returns a pointer to the resource representing a database.
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

func createDatabase(d *schema.ResourceData, builder *snowflake.DatabaseBuilder, meta interface{}) error {
	db := meta.(*sql.DB)
	q := builder.Create()
	name := d.Get("name").(string)

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating database %v", name)
	}

	d.SetId(name)

	return ReadDatabase(d, meta)
}

// CreateDatabase implements schema.CreateFunc.
func CreateDatabase(d *schema.ResourceData, meta interface{}) error {
	// TODO: Migrate database from share and from replica to iterative approach
	if _, ok := d.GetOk("from_share"); ok {
		return createDatabaseFromShare(d, meta)
	}

	if _, ok := d.GetOk("from_replica"); ok {
		return createDatabaseFromReplica(d, meta)
	}

	name := d.Get("name").(string)
	builder := snowflake.Database(name)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("is_transient"); ok && v.(bool) {
		builder.Transient()
	}

	if v, ok := d.GetOk("from_database"); ok {
		builder.Clone(v.(string))
	}

	if v, ok := d.GetOk("data_retention_time_in_days"); ok {
		builder.WithDataRetentionDays(v.(int))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	// If set, verify parameters are valid and attempt to enable replication
	if v, ok := d.GetOk("replication_configuration"); ok {
		replicationConfiguration := v.([]interface{})[0].(map[string]interface{})
		ignoreEditionCheck := replicationConfiguration["ignore_edition_check"].(bool)

		if !ignoreEditionCheck {
			return errors.New("error enabling replication - ignore edition check was set to false")
		}
		resource := createDatabase(d, builder, meta)
		if err := enableReplication(d, meta, replicationConfiguration); err != nil {
			return errors.Wrapf(err, "error enabling replication - account does not exist or System Parameter ENABLE_ACCOUNT_DATABASE_REPLICATION must be set to true")
		}
		return resource
	}

	return createDatabase(d, builder, meta)
}

func enableReplication(d *schema.ResourceData, meta interface{}, replicationConfig map[string]interface{}) error {
	db := meta.(*sql.DB)
	primaryDBName := d.Get("name").(string)
	accounts := replicationConfig["accounts"].([]interface{})
	accountsToEnableReplication := strings.Join(expandStringList(accounts), ", ")
	enableReplicationStmt := fmt.Sprintf(`ALTER DATABASE "%s" ENABLE REPLICATION TO ACCOUNTS %s`, primaryDBName, accountsToEnableReplication)
	return snowflake.Exec(db, enableReplicationStmt)
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

	if comment, ok := d.GetOk("comment"); ok {
		builder.WithComment(comment.(string))
	}

	err := snowflake.Exec(db, builder.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating database %v from share %v.%v", name, prov, share)
	}

	d.SetId(name)

	return ReadDatabase(d, meta)
}

func createDatabaseFromReplica(d *schema.ResourceData, meta interface{}) error {
	sourceDB := d.Get("from_replica").(string)

	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	builder := snowflake.DatabaseFromReplica(name, sourceDB)

	err := snowflake.Exec(db, builder.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating a secondary database %v from database %v", name, sourceDB)
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

	// reset the options before reading back from the DB
	err = d.Set("is_transient", false)
	if err != nil {
		return err
	}

	if opts := database.Options.String; opts != "" {
		for _, opt := range strings.Split(opts, ", ") {
			switch opt {
			case "TRANSIENT":
				err = d.Set("is_transient", true)
				if err != nil {
					return err
				}
			}
		}
	}

	return d.Set("data_retention_time_in_days", i)
}

func UpdateDatabase(d *schema.ResourceData, meta interface{}) error {
	dbName := d.Id()
	builder := snowflake.Database(dbName)
	db := meta.(*sql.DB)

	// If replication configuration changes, need to update accounts that have permission to replicate database
	if d.HasChange("replication_configuration") {
		oldConfig, newConfig := d.GetChange("replication_configuration")
		newConfigLength := len(newConfig.([]interface{}))
		oldConfigLength := len(oldConfig.([]interface{}))
		// Enable replication for any new accounts and disable replication for removed accounts
		if newConfigLength > 0 {
			newAccounts := extractInterfaceFromAttribute(newConfig, "accounts")
			enableQuery := builder.EnableReplicationAccounts(dbName, strings.Join(expandStringList(newAccounts), ", "))
			err := snowflake.Exec(db, enableQuery)
			if err != nil {
				return errors.Wrapf(err, "error enabling replication configuration with statement %v", enableQuery)
			}
		}

		if oldConfigLength > 0 {
			oldAccounts := extractInterfaceFromAttribute(oldConfig, "accounts")
			var accountsToDisableReplication []interface{}
			if newConfigLength > 0 {
				newAccounts := extractInterfaceFromAttribute(newConfig, "accounts")
				accountsToDisableReplication = builder.GetRemovedAccountsFromReplicationConfiguration(oldAccounts, newAccounts)
			} else {
				accountsToDisableReplication = builder.GetRemovedAccountsFromReplicationConfiguration(oldAccounts, nil)
			}
			// If accounts were found to be removed, disable replication
			if len(accountsToDisableReplication) > 0 {
				disableQuery := builder.DisableReplicationAccounts(dbName, strings.Join(expandStringList(accountsToDisableReplication), ", "))
				err := snowflake.Exec(db, disableQuery)
				if err != nil {
					return errors.Wrapf(err, "error disabling replication configuration with statement %v", disableQuery)
				}
			}
		}
	}

	if d.HasChange("name") {
		name := d.Get("name")
		q := builder.Rename(name.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating database name on %v", d.Id())
		}
		d.SetId(fmt.Sprintf("%v", name.(string)))
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating database comment on %v", d.Id())
		}
	}

	if d.HasChange("data_retention_time_in_days") {
		days := d.Get("data_retention_time_in_days")

		q := builder.ChangeDataRetentionDays(days.(int))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating data retention days on %v", d.Id())
		}
	}

	tagChangeErr := handleTagChanges(db, d, builder)
	if tagChangeErr != nil {
		return tagChangeErr
	}

	return ReadDatabase(d, meta)
}

func DeleteDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()

	q := snowflake.Database(name).Drop()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting database %v", d.Id())
	}

	d.SetId("")

	return nil
}

func extractInterfaceFromAttribute(config interface{}, attribute string) []interface{} {
	return config.([]interface{})[0].(map[string]interface{})[attribute].([]interface{})
}
