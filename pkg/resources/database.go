package resources

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var databaseSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
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
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Number of days for which Snowflake retains historical data for performing Time Travel actions (SELECT, CLONE, UNDROP) on the object. A value of 0 effectively disables Time Travel for the specified database, schema, or table. For more information, see Understanding & Using Time Travel.",
		Default:     1,
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

// CreateDatabase implements schema.CreateFunc.
func CreateDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	// Is it a Shared Database?
	if fromShare, ok := d.GetOk("from_share"); ok {
		account := fromShare.(map[string]interface{})["provider"].(string)
		share := fromShare.(map[string]interface{})["share"].(string)
		shareID := sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifierFromAccountLocator(account), sdk.NewAccountObjectIdentifier(share))
		opts := &sdk.CreateSharedDatabaseOptions{}
		if v, ok := d.GetOk("comment"); ok {
			opts.Comment = sdk.String(v.(string))
		}
		err := client.Databases.CreateShared(ctx, id, shareID, opts)
		if err != nil {
			return fmt.Errorf("error creating database %v: %w", name, err)
		}
		d.SetId(name)
		if v, ok := d.GetOk("replication_configuration"); ok {
			replicationConfiguration := v.([]interface{})[0].(map[string]interface{})
			accounts := replicationConfiguration["accounts"].([]interface{})
			accountIDs := make([]sdk.AccountIdentifier, len(accounts))
			for i, account := range accounts {
				accountIDs[i] = sdk.NewAccountIdentifierFromAccountLocator(account.(string))
			}
			opts := &sdk.AlterDatabaseReplicationOptions{
				EnableReplication: &sdk.EnableReplication{
					ToAccounts: accountIDs,
				},
			}
			if ignoreEditionCheck, ok := replicationConfiguration["ignore_edition_check"]; ok {
				opts.EnableReplication.IgnoreEditionCheck = sdk.Bool(ignoreEditionCheck.(bool))
			}
			err := client.Databases.AlterReplication(ctx, id, opts)
			if err != nil {
				return fmt.Errorf("error enabling replication for database %v: %w", name, err)
			}
		}
		return ReadDatabase(d, meta)
	}
	// Is it a Secondary Database?
	if primaryName, ok := d.GetOk("from_replica"); ok {
		primaryID := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(primaryName.(string))
		opts := &sdk.CreateSecondaryDatabaseOptions{
			DataRetentionTimeInDays: sdk.Int(d.Get("data_retention_time_in_days").(int)),
		}
		err := client.Databases.CreateSecondary(ctx, id, primaryID, opts)
		if err != nil {
			return fmt.Errorf("error creating database %v: %w", name, err)
		}
		d.SetId(name)
		// todo: add failover_configuration block
		return ReadDatabase(d, meta)
	}

	// Otherwise it is a Standard Database
	opts := sdk.CreateDatabaseOptions{}
	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	if v, ok := d.GetOk("is_transient"); ok && v.(bool) {
		opts.Transient = sdk.Bool(v.(bool))
	}

	if v, ok := d.GetOk("from_database"); ok {
		opts.Clone = &sdk.Clone{
			SourceObject: sdk.NewAccountObjectIdentifier(v.(string)),
		}
	}

	if v, ok := d.GetOk("data_retention_time_in_days"); ok {
		opts.DataRetentionTimeInDays = sdk.Int(v.(int))
	}

	err := client.Databases.Create(ctx, id, &opts)
	if err != nil {
		return fmt.Errorf("error creating database %v: %w", name, err)
	}
	d.SetId(name)
	return ReadDatabase(d, meta)
}

func ReadDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	name := d.Id()
	id := sdk.NewAccountObjectIdentifier(name)

	database, err := client.Databases.ShowByID(ctx, id)
	if err != nil {
		return err
	}
	if err := d.Set("name", database.Name); err != nil {
		return err
	}
	if err := d.Set("comment", database.Comment); err != nil {
		return err
	}

	if err := d.Set("data_retention_time_in_days", database.RetentionTime); err != nil {
		return err
	}

	if err := d.Set("is_transient", database.Transient); err != nil {
		return err
	}

	return nil
}

func UpdateDatabase(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	id := sdk.NewAccountObjectIdentifier(name)
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		opts := &sdk.AlterDatabaseOptions{
			NewName: sdk.NewAccountObjectIdentifier(newName),
		}
		err := client.Databases.Alter(ctx, id, opts)
		if err != nil {
			return fmt.Errorf("error updating database name on %v err = %w", d.Id(), err)
		}
		d.SetId(newName)
		id = sdk.NewAccountObjectIdentifier(newName)
	}

	if d.HasChange("comment") {
		comment := ""
		if c := d.Get("comment"); c != nil {
			comment = c.(string)
		}
		opts := &sdk.AlterDatabaseOptions{
			Set: &sdk.DatabaseSet{
				Comment: sdk.String(comment),
			},
		}
		err := client.Databases.Alter(ctx, id, opts)
		if err != nil {
			return fmt.Errorf("error updating database comment on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("data_retention_time_in_days") {
		days := d.Get("data_retention_time_in_days").(int)
		opts := &sdk.AlterDatabaseOptions{
			Set: &sdk.DatabaseSet{
				DataRetentionTimeInDays: sdk.Int(days),
			},
		}
		err := client.Databases.Alter(ctx, id, opts)
		if err != nil {
			return fmt.Errorf("error updating database data retention time on %v err = %w", d.Id(), err)
		}
	}

	// If replication configuration changes, need to update accounts that have permission to replicate database
	if d.HasChange("replication_configuration") {
		oldConfig, newConfig := d.GetChange("replication_configuration")
		newAccounts := newConfig.([]interface{})[0].(map[string]interface{})["accounts"].([]interface{})
		newAccountIDs := make([]sdk.AccountIdentifier, len(newAccounts))
		for i, account := range newAccounts {
			newAccountIDs[i] = sdk.NewAccountIdentifierFromAccountLocator(account.(string))
		}
		oldAccounts := oldConfig.([]interface{})[0].(map[string]interface{})["accounts"].([]interface{})
		oldAccountIDs := make([]sdk.AccountIdentifier, len(oldAccounts))
		for i, account := range oldAccounts {
			oldAccountIDs[i] = sdk.NewAccountIdentifierFromAccountLocator(account.(string))
		}
		accountsToRemove := make([]sdk.AccountIdentifier, 0)
		accountsToAdd := make([]sdk.AccountIdentifier, 0)
		// Find accounts to remove
		for _, oldAccountID := range oldAccountIDs {
			if !slices.Contains(newAccountIDs, oldAccountID) {
				accountsToRemove = append(accountsToRemove, oldAccountID)
			}
		}

		// Find accounts to add
		for _, newAccountID := range newAccountIDs {
			if !slices.Contains(oldAccountIDs, newAccountID) {
				accountsToAdd = append(accountsToAdd, newAccountID)
			}
		}
		if len(accountsToAdd) > 0 {
			opts := &sdk.AlterDatabaseReplicationOptions{
				EnableReplication: &sdk.EnableReplication{
					ToAccounts: accountsToAdd,
				},
			}
			if ignoreEditionCheck := d.Get("ignore_edition_check").(bool); ignoreEditionCheck {
				opts.EnableReplication.IgnoreEditionCheck = sdk.Bool(ignoreEditionCheck)
			}
			err := client.Databases.AlterReplication(ctx, id, opts)
			if err != nil {
				return fmt.Errorf("error enabling replication configuration on %v err = %w", d.Id(), err)
			}
		}

		if len(accountsToRemove) > 0 {
			opts := &sdk.AlterDatabaseReplicationOptions{
				DisableReplication: &sdk.DisableReplication{
					ToAccounts: accountsToRemove,
				},
			}
			err := client.Databases.AlterReplication(ctx, id, opts)
			if err != nil {
				return fmt.Errorf("error disabling replication configuration on %v err = %w", d.Id(), err)
			}
		}
	}

	return ReadDatabase(d, meta)
}

func DeleteDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	name := d.Id()
	id := sdk.NewAccountObjectIdentifier(name)
	err := client.Databases.Drop(ctx, id, &sdk.DropDatabaseOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
