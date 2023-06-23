package resources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"golang.org/x/exp/slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var failoverGroupSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the failover group. The identifier must start with an alphabetic character and cannot contain spaces or special characters unless the identifier string is enclosed in double quotes (e.g. \"My object\"). Identifiers enclosed in double quotes are also case-sensitive.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(old, new)
		},
	},
	"object_types": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Type(s) of objects for which you are enabling replication and failover from the source account to the target account. The following object types are supported: \"ACCOUNT PARAMETERS\", \"DATABASES\", \"INTEGRATIONS\", \"NETWORK POLICIES\", \"RESOURCE MONITORS\", \"ROLES\", \"SHARES\", \"USERS\", \"WAREHOUSES\"",
	},
	"allowed_databases": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the database or list of databases for which you are enabling replication and failover from the source account to the target account. The OBJECT_TYPES list must include DATABASES to set this parameter.",
	},
	"allowed_shares": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the share or list of shares for which you are enabling replication and failover from the source account to the target account. The OBJECT_TYPES list must include SHARES to set this parameter.",
	},
	"allowed_integration_types": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Type(s) of integrations for which you are enabling replication and failover from the source account to the target account. This property requires that the OBJECT_TYPES list include INTEGRATIONS to set this parameter. The following integration types are supported: \"SECURITY INTEGRATIONS\", \"API INTEGRATIONS\"",
	},
	"allowed_accounts": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the target account or list of target accounts to which replication and failover of specified objects from the source account is enabled. Secondary failover groups in the target accounts in this list can be promoted to serve as the primary failover group in case of failover. Expected in the form <org_name>.<target_account_name>",
	},
	"ignore_edition_check": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"from_replica"},
		Description:   "Allows replicating objects to accounts on lower editions.",
	},
	"from_replica": {
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		MaxItems:      1,
		ConflictsWith: []string{"object_types", "allowed_accounts", "allowed_databases", "allowed_shares", "allowed_integration_types", "ignore_edition_check", "replication_schedule"},
		Description:   "Specifies the name of the replica to use as the source for the failover group.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"organization_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of your Snowflake organization.",
				},
				"source_account_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Source account from which you are enabling replication and failover of the specified objects.",
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Identifier for the primary failover group in the source account.",
				},
			},
		},
	},
	"replication_schedule": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		Description:   "Specifies the schedule for refreshing secondary failover groups.",
		ConflictsWith: []string{"from_replica"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cron": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"replication_schedule.interval"},
					Description:   "Specifies the cron expression for the replication schedule. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"expression": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the cron expression for the replication schedule. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
							},
							"time_zone": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the time zone for secondary group refresh.",
							},
						},
					},
				},
				"interval": {
					Type:          schema.TypeInt,
					Optional:      true,
					ConflictsWith: []string{"replication_schedule.cron"},
					Description:   "Specifies the interval in minutes for the replication schedule. The interval must be greater than 0 and less than 1440 (24 hours).",
				},
			},
		},
	},
}

// FailoverGroup returns a pointer to the resource representing a failover group.
func FailoverGroup() *schema.Resource {
	return &schema.Resource{
		Create: CreateFailoverGroup,
		Read:   ReadFailoverGroup,
		Update: UpdateFailoverGroup,
		Delete: DeleteFailoverGroup,

		Schema: failoverGroupSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateFailoverGroup implements schema.CreateFunc.
func CreateFailoverGroup(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	// getting required attributes
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	// if from_replica is set, then we are creating a failover group from an existing replica
	if v, ok := d.GetOk("from_replica"); ok {
		fromReplica := v.([]interface{})[0].(map[string]interface{})
		organizationName := fromReplica["organization_name"].(string)
		sourceAccountName := fromReplica["source_account_name"].(string)
		sourceFailoverGroupName := fromReplica["name"].(string)

		primaryFailoverGroupID := sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(organizationName, sourceAccountName), sdk.NewAccountObjectIdentifier(sourceFailoverGroupName))
		err := client.FailoverGroups.CreateSecondaryReplicationGroup(ctx, id, primaryFailoverGroupID, nil)
		if err != nil {
			return err
		}
		d.SetId(name)
		return ReadFailoverGroup(d, meta)
	}

	// these two are required attributes if from_replica is not set
	if _, ok := d.GetOk("object_types"); !ok {
		return errors.New("object_types is required when not creating from a replica")
	}
	objectTypesList := expandStringList(d.Get("object_types").(*schema.Set).List())
	objectTypes := make([]sdk.PluralObjectType, len(objectTypesList))
	for i, v := range objectTypesList {
		objectTypes[i] = sdk.PluralObjectType(v)
	}

	if _, ok := d.GetOk("allowed_accounts"); !ok {
		return errors.New("allowed_accounts is required when not creating from a replica")
	}
	aaList := expandStringList(d.Get("allowed_accounts").(*schema.Set).List())
	allowedAccounts := make([]sdk.AccountIdentifier, len(aaList))
	for i, v := range aaList {
		// validation since we cannot do that in the ValidateFunc
		parts := strings.Split(v, ".")
		if len(parts) != 2 {
			return fmt.Errorf("allowed_account %s cannot be an account locator and must be of the format <org_name>.<target_account_name>", allowedAccounts[i])
		}
		organizationName := parts[0]
		accountName := parts[1]
		allowedAccounts[i] = sdk.NewAccountIdentifier(organizationName, accountName)
	}

	var opts sdk.CreateFailoverGroupOptions
	// setting optional attributes
	if v, ok := d.GetOk("allowed_databases"); ok {
		allowedDatabasesList := expandStringList(v.(*schema.Set).List())
		allowedDatabaseIdentifiers := make([]sdk.AccountObjectIdentifier, len(allowedDatabasesList))
		for i, v := range allowedDatabasesList {
			allowedDatabaseIdentifiers[i] = sdk.NewAccountObjectIdentifier(v)
		}
		opts.AllowedDatabases = allowedDatabaseIdentifiers
	}

	if v, ok := d.GetOk("allowed_shares"); ok {
		allowedSharesList := expandStringList(v.(*schema.Set).List())
		allowedShareIdentifiers := make([]sdk.AccountObjectIdentifier, len(allowedSharesList))
		for i, v := range allowedSharesList {
			allowedShareIdentifiers[i] = sdk.NewAccountObjectIdentifier(v)
		}
		opts.AllowedShares = allowedShareIdentifiers
	}

	if v, ok := d.GetOk("allowed_integration_types"); ok {
		allowedIntegrationTypesList := expandStringList(v.(*schema.Set).List())
		allowedIntegrationTypes := make([]sdk.IntegrationType, len(allowedIntegrationTypesList))
		for i, v := range allowedIntegrationTypesList {
			allowedIntegrationTypes[i] = sdk.IntegrationType(v)
		}
		opts.AllowedIntegrationTypes = allowedIntegrationTypes
	}

	if v, ok := d.GetOk("ignore_edition_check"); ok {
		opts.IgnoreEditionCheck = sdk.Bool(v.(bool))
	}

	if v, ok := d.GetOk("replication_schedule"); ok {
		replicationSchedule := v.([]interface{})[0].(map[string]interface{})
		if v, ok := replicationSchedule["cron"]; ok {
			c := v.([]interface{})
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)
				if v, ok := cron["time_zone"]; ok {
					timeZone := v.(string)
					cronExpression += " " + timeZone
				}
				opts.ReplicationSchedule = sdk.String("USING CRON " + cronExpression)
			}
		}
		if v, ok := replicationSchedule["interval"]; ok {
			interval := v.(int)
			if interval > 1 {
				opts.ReplicationSchedule = sdk.String(fmt.Sprintf("%d MINUTE", interval))
			}
		}
	}

	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, &opts)
	if err != nil {
		return err
	}

	d.SetId(name)
	return nil
}

// ReadFailoverGroup implements schema.ReadFunc.
func ReadFailoverGroup(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	name := d.Id()
	id := sdk.NewAccountObjectIdentifier(name)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	if err != nil {
		return err
	}

	if err := d.Set("name", failoverGroup.Name); err != nil {
		return err
	}
	// if the failover group is created from a replica, then we do not want to get the other values
	if _, ok := d.GetOk("from_replica"); ok {
		return nil
	}

	replicationSchedule := failoverGroup.ReplicationSchedule
	if replicationSchedule != "" {
		if strings.Contains(replicationSchedule, "MINUTE") {
			interval, err := strconv.Atoi(strings.TrimSuffix(replicationSchedule, " MINUTE"))
			if err != nil {
				return err
			}
			err = d.Set("replication_schedule", []interface{}{
				map[string]interface{}{
					"interval": interval,
				},
			})
			if err != nil {
				return err
			}
		} else {
			repScheduleParts := strings.Split(replicationSchedule, " ")
			timeZone := repScheduleParts[len(repScheduleParts)-1]
			expression := strings.TrimSuffix(strings.TrimPrefix(replicationSchedule, "USING CRON "), " "+timeZone)
			err = d.Set("replication_schedule", []interface{}{
				map[string]interface{}{
					"cron": []interface{}{
						map[string]interface{}{
							"expression": expression,
							"time_zone":  timeZone,
						},
					},
				},
			})
			if err != nil {
				return err
			}
		}
	}
	// object types
	objectTypes := make([]interface{}, len(failoverGroup.ObjectTypes))
	for i, v := range failoverGroup.ObjectTypes {
		objectTypes[i] = string(v)
	}
	objectTypesSet := schema.NewSet(schema.HashString, objectTypes)
	if err := d.Set("object_types", objectTypesSet); err != nil {
		return err
	}

	// integration types
	allowedIntegrationTypes := make([]interface{}, len(failoverGroup.AllowedIntegrationTypes))
	for i, v := range failoverGroup.AllowedIntegrationTypes {
		allowedIntegrationTypes[i] = string(v)
	}

	allowedIntegrationsTypesSet := schema.NewSet(schema.HashString, allowedIntegrationTypes)
	if err := d.Set("allowed_integration_types", allowedIntegrationsTypesSet); err != nil {
		return err
	}

	// allowed accounts
	allowedAccounts := make([]interface{}, len(failoverGroup.AllowedAccounts))
	for i, v := range failoverGroup.AllowedAccounts {
		allowedAccounts[i] = v.Name()
	}
	allowedAccountsSet := schema.NewSet(schema.HashString, allowedAccounts)
	if err := d.Set("allowed_accounts", allowedAccountsSet); err != nil {
		return err
	}

	// allowed databases
	databases, err := client.FailoverGroups.ShowDatabases(ctx, id)
	if err != nil {
		return err
	}
	allowedDatabases := make([]interface{}, len(databases))
	for i, database := range databases {
		allowedDatabases[i] = database.Name()
	}
	allowedDatabasesSet := schema.NewSet(schema.HashString, allowedDatabases)
	if len(allowedDatabases) > 0 {
		if err := d.Set("allowed_databases", allowedDatabasesSet); err != nil {
			return err
		}
	} else {
		if err := d.Set("allowed_databases", nil); err != nil {
			return err
		}
	}

	// allowed shares
	shares, err := client.FailoverGroups.ShowShares(ctx, id)
	if err != nil {
		return err
	}
	allowedShares := make([]interface{}, len(shares))
	for i, share := range shares {
		allowedShares[i] = share.Name()
	}
	allowedSharesSet := schema.NewSet(schema.HashString, allowedShares)
	if len(allowedShares) > 0 {
		if err := d.Set("allowed_shares", allowedSharesSet); err != nil {
			return err
		}
	} else {
		if err := d.Set("allowed_shares", nil); err != nil {
			return err
		}
	}

	return nil
}

// UpdateFailoverGroup implements schema.UpdateFunc.
func UpdateFailoverGroup(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	name := d.Id()
	id := sdk.NewAccountObjectIdentifier(name)

	// alter failover group <name> set ...
	opts := &sdk.AlterSourceFailoverGroupOptions{
		Set: &sdk.FailoverGroupSet{},
	}
	runSet := false

	if d.HasChange("object_types") {
		_, n := d.GetChange("object_types")
		newObjectTypes := expandStringList(n.(*schema.Set).List())
		objectTypes := make([]sdk.PluralObjectType, len(newObjectTypes))
		for i, v := range newObjectTypes {
			objectTypes[i] = sdk.PluralObjectType(v)
		}
		opts.Set.ObjectTypes = objectTypes
		if slices.Contains(objectTypes, sdk.PluralObjectTypeIntegrations) {
			ait := expandStringList(d.Get("allowed_integration_types").(*schema.Set).List())
			allowedIntegrationTypes := make([]sdk.IntegrationType, len(ait))
			for i, v := range ait {
				allowedIntegrationTypes[i] = sdk.IntegrationType(v)
			}
			opts.Set.AllowedIntegrationTypes = allowedIntegrationTypes
		}
		runSet = true
	}

	if d.HasChange("replication_schedule") {
		_, n := d.GetChange("replication_schedule")
		replicationSchedule := n.([]interface{})[0].(map[string]interface{})
		c := replicationSchedule["cron"].([]interface{})
		if len(c) > 0 {
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)
				cronExpression = "USING CRON " + cronExpression
				if v, ok := cron["time_zone"]; ok {
					timeZone := v.(string)
					if timeZone != "" {
						cronExpression = cronExpression + " " + timeZone
					}
				}
				opts.Set.ReplicationSchedule = &cronExpression
			}
		} else {
			opts.Set.ReplicationSchedule = sdk.String(fmt.Sprintf("%d MINUTE", replicationSchedule["interval"].(int)))
		}
		runSet = true
	}

	if d.HasChange("allowed_integration_types") {
		_, n := d.GetChange("allowed_integration_types")
		newAllowedIntegrationTypes := expandStringList(n.(*schema.Set).List())
		allowedIntegrationTypes := make([]sdk.IntegrationType, len(newAllowedIntegrationTypes))
		for i, v := range newAllowedIntegrationTypes {
			allowedIntegrationTypes[i] = sdk.IntegrationType(v)
		}
		opts.Set.AllowedIntegrationTypes = allowedIntegrationTypes
		runSet = true
	}
	if runSet {
		if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
			return err
		}
	}

	if d.HasChange("allowed_databases") {
		o, n := d.GetChange("allowed_databases")
		oad := expandStringList(o.(*schema.Set).List())
		oldAllowedDatabases := make([]sdk.AccountObjectIdentifier, len(oad))
		for i, v := range oad {
			oldAllowedDatabases[i] = sdk.NewAccountObjectIdentifier(v)
		}
		nad := expandStringList(n.(*schema.Set).List())
		newAllowedDatabases := make([]sdk.AccountObjectIdentifier, len(nad))
		for i, v := range nad {
			newAllowedDatabases[i] = sdk.NewAccountObjectIdentifier(v)
		}

		var removedDatabases []sdk.AccountObjectIdentifier
		for _, v := range oldAllowedDatabases {
			if !slices.Contains(newAllowedDatabases, v) {
				removedDatabases = append(removedDatabases, v)
			}
		}
		if len(removedDatabases) > 0 {
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Remove: &sdk.FailoverGroupRemove{
					AllowedDatabases: removedDatabases,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
				return fmt.Errorf("error removing allowed databases for failover group %v err = %w", name, err)
			}
		}

		var addedDatabases []sdk.AccountObjectIdentifier
		for _, v := range newAllowedDatabases {
			if !slices.Contains(oldAllowedDatabases, v) {
				addedDatabases = append(addedDatabases, v)
			}
		}

		if len(addedDatabases) > 0 {
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Add: &sdk.FailoverGroupAdd{
					AllowedDatabases: addedDatabases,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
				return fmt.Errorf("error removing allowed databases for failover group %v err = %w", name, err)
			}
		}
	}

	if d.HasChange("allowed_shares") {
		o, n := d.GetChange("allowed_shares")
		oad := expandStringList(o.(*schema.Set).List())
		oldAllowedShares := make([]sdk.AccountObjectIdentifier, len(oad))
		for i, v := range oad {
			oldAllowedShares[i] = sdk.NewAccountObjectIdentifier(v)
		}
		nad := expandStringList(n.(*schema.Set).List())
		newAllowedShares := make([]sdk.AccountObjectIdentifier, len(nad))
		for i, v := range nad {
			newAllowedShares[i] = sdk.NewAccountObjectIdentifier(v)
		}

		var removedShares []sdk.AccountObjectIdentifier
		for _, v := range oldAllowedShares {
			if !slices.Contains(newAllowedShares, v) {
				removedShares = append(removedShares, v)
			}
		}
		if len(removedShares) > 0 {
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Remove: &sdk.FailoverGroupRemove{
					AllowedShares: removedShares,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
				return fmt.Errorf("error removing allowed shares for failover group %v err = %w", name, err)
			}
		}

		var addedShares []sdk.AccountObjectIdentifier
		for _, v := range newAllowedShares {
			if !slices.Contains(oldAllowedShares, v) {
				addedShares = append(addedShares, v)
			}
		}

		if len(addedShares) > 0 {
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Add: &sdk.FailoverGroupAdd{
					AllowedShares: addedShares,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
				return fmt.Errorf("error removing allowed shares for failover group %v err = %w", name, err)
			}
		}
	}

	if d.HasChange("allowed_accounts") {
		o, n := d.GetChange("allowed_accounts")
		oad := expandStringList(o.(*schema.Set).List())
		oldAllowedAccounts := make([]sdk.AccountIdentifier, len(oad))
		for i, v := range oad {
			parts := strings.Split(v, ".")
			organizationName := parts[0]
			accountName := parts[1]
			accountIdentifier := sdk.NewAccountIdentifier(accountName, organizationName)
			oldAllowedAccounts[i] = accountIdentifier
		}
		nad := expandStringList(n.(*schema.Set).List())
		newAllowedAccounts := make([]sdk.AccountIdentifier, len(nad))
		for i, v := range nad {
			parts := strings.Split(v, ".")
			organizationName := parts[0]
			accountName := parts[1]
			accountIdentifier := sdk.NewAccountIdentifier(accountName, organizationName)
			newAllowedAccounts[i] = accountIdentifier
		}

		var removedAccounts []sdk.AccountIdentifier
		for _, v := range oldAllowedAccounts {
			if !slices.Contains(newAllowedAccounts, v) {
				removedAccounts = append(removedAccounts, v)
			}
		}
		if len(removedAccounts) > 0 {
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Remove: &sdk.FailoverGroupRemove{
					AllowedAccounts: removedAccounts,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
				return fmt.Errorf("error removing allowed accounts for failover group %v err = %w", name, err)
			}
		}

		var addedAccounts []sdk.AccountIdentifier
		for _, v := range newAllowedAccounts {
			if !slices.Contains(oldAllowedAccounts, v) {
				addedAccounts = append(addedAccounts, v)
			}
		}

		if len(addedAccounts) > 0 {
			opts := &sdk.AlterSourceFailoverGroupOptions{
				Add: &sdk.FailoverGroupAdd{
					AllowedAccounts: addedAccounts,
				},
			}
			if err := client.FailoverGroups.AlterSource(ctx, id, opts); err != nil {
				return fmt.Errorf("error removing allowed accounts for failover group %v err = %w", name, err)
			}
		}
	}

	return ReadFailoverGroup(d, meta)
}

// DeleteFailoverGroup implements schema.DeleteFunc.
func DeleteFailoverGroup(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	name := d.Id()
	id := sdk.NewAccountObjectIdentifier(name)
	ctx := context.Background()
	err := client.FailoverGroups.Drop(ctx, id, &sdk.DropFailoverGroupOptions{IfExists: sdk.Bool(true)})
	if err != nil {
		return fmt.Errorf("error deleting failover group %v err = %w", name, err)
	}

	d.SetId("")
	return nil
}
