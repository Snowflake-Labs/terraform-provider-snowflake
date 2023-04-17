package resources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"golang.org/x/exp/slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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

	// getting required attributes
	name := d.Get("name").(string)
	builder := snowflake.NewFailoverGroupBuilder(name)

	// if from_replica is set, then we are creating a failover group from an existing replica
	if v, ok := d.GetOk("from_replica"); ok {
		fromReplica := v.([]interface{})[0].(map[string]interface{})
		organizationName := fromReplica["organization_name"].(string)
		sourceAccountName := fromReplica["source_account_name"].(string)
		sourceFailoverGroupName := fromReplica["name"].(string)
		fullyQualifiedFailoverGroupIdentifier := fmt.Sprintf("%s.%s.%s", organizationName, sourceAccountName, sourceFailoverGroupName)
		stmt := builder.CreateFromReplica(fullyQualifiedFailoverGroupIdentifier)
		if err := snowflake.Exec(db, stmt); err != nil {
			return fmt.Errorf("error creating failover group %v err = %w", name, err)
		}
		d.SetId(name)
		return ReadFailoverGroup(d, meta)
	}

	// these two are required attributes if from_replica is not set
	if _, ok := d.GetOk("object_types"); !ok {
		return errors.New("object_types is required when not creating from a replica")
	}
	ot := d.Get("object_types").(*schema.Set).List()
	objectTypes := make([]string, len(ot))
	for i, v := range ot {
		objectTypes[i] = v.(string)
	}
	builder.WithObjectTypes(objectTypes)

	if _, ok := d.GetOk("allowed_accounts"); !ok {
		return errors.New("allowed_accounts is required when not creating from a replica")
	}
	aa := d.Get("allowed_accounts").(*schema.Set).List()
	allowedAccounts := make([]string, len(aa))
	for i, v := range aa {
		allowedAccounts[i] = v.(string)
		// validation since we cannot do that in the ValidateFunc
		parts := strings.Split(allowedAccounts[i], ".")
		if len(parts) != 2 {
			return fmt.Errorf("allowed_account %s must be of the format <org_name>.<target_account_name>", allowedAccounts[i])
		}
	}
	builder.WithAllowedAccounts(allowedAccounts)

	// setting optional attributes
	if v, ok := d.GetOk("allowed_databases"); ok {
		ad := v.(*schema.Set).List()
		allowedDatabases := make([]string, len(ad))
		for i, v := range ad {
			allowedDatabases[i] = v.(string)
		}
		builder.WithAllowedDatabases(allowedDatabases)
	}

	if v, ok := d.GetOk("allowed_shares"); ok {
		as := v.(*schema.Set).List()
		allowedShares := make([]string, len(as))
		for i, v := range as {
			allowedShares[i] = v.(string)
		}
		builder.WithAllowedShares(allowedShares)
	}

	if v, ok := d.GetOk("allowed_integration_types"); ok {
		aits := v.(*schema.Set).List()
		allowedIntegrationTypes := make([]string, len(aits))
		for i, v := range aits {
			allowedIntegrationTypes[i] = v.(string)
		}

		builder.WithAllowedIntegrationTypes(allowedIntegrationTypes)
	}

	if v, ok := d.GetOk("ignore_edition_check"); ok {
		builder.WithIgnoreEditionCheck(v.(bool))
	}

	if v, ok := d.GetOk("replication_schedule"); ok {
		replicationSchedule := v.([]interface{})[0].(map[string]interface{})
		if v, ok := replicationSchedule["cron"]; ok {
			c := v.([]interface{})
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)
				builder.WithReplicationScheduleCronExpression(cronExpression)
				if v, ok := cron["time_zone"]; ok {
					timeZone := v.(string)
					builder.WithReplicationScheduleTimeZone(timeZone)
				}
			}
		}
		if v, ok := replicationSchedule["interval"]; ok {
			interval := v.(int)
			if interval > 1 {
				builder.WithReplicationScheduleInterval(interval)
			}
		}
	}

	q := builder.Create()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error creating failover group %v err = %w", name, err)
	}

	d.SetId(name)
	return nil
}

// ReadFailoverGroup implements schema.ReadFunc.
func ReadFailoverGroup(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()

	stmt := "select current_account()"
	row := db.QueryRow(stmt)
	var accountLocator string
	if err := row.Scan(&accountLocator); err != nil {
		return fmt.Errorf("error getting current account err = %w", err)
	}

	failoverGroups, err := snowflake.ListFailoverGroups(db, accountLocator)
	if err != nil {
		return fmt.Errorf("error listing failover groups err = %w", err)
	}

	var failoverGroup snowflake.FailoverGroup
	// find the failover group we are looking for by matching the name
	for _, fg := range failoverGroups {
		if strings.EqualFold(fg.Name.String, name) && strings.EqualFold(fg.AccountLocator.String, accountLocator) {
			failoverGroup = fg
		}
	}

	if failoverGroup.Name.String == "" {
		log.Printf("[DEBUG] failover group (%v) not found when listing all failover groups in account", name)
		d.SetId("")
		return nil
	}

	if err := d.Set("name", failoverGroup.Name.String); err != nil {
		return err
	}
	// if the failover group is created from a replica, then we do not want to get the other values
	if _, ok := d.GetOk("from_replica"); ok {
		log.Printf("[DEBUG] failover group %v is created from a replica, rest of values are computed\n", name)
		return nil
	}

	replicationSchedule := failoverGroup.ReplicationSchedule.String
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

	objectTypes := helpers.StringListToList(failoverGroup.ObjectTypes.String)

	// this is basically a hack to get around the fact that the API returns the object types in a different order than what is set
	// this logic could also be put in the diff suppress function, but I think it is better to do it here.
	currentObjectTypeList := d.Get("object_types").(*schema.Set).List()
	if len(currentObjectTypeList) != len(objectTypes) {
		log.Printf("[DEBUG] object types are different, current: %v, new: %v", currentObjectTypeList, objectTypes)
		if err := d.Set("object_types", objectTypes); err != nil {
			return err
		}
	}

	for _, v := range currentObjectTypeList {
		if !slices.Contains(objectTypes, v.(string)) {
			log.Printf("[DEBUG] object types are different, current: %v, new: %v", currentObjectTypeList, objectTypes)
			if err := d.Set("object_types", objectTypes); err != nil {
				return err
			}
			break
		}
	}

	allowedIntegrationTypes := failoverGroup.AllowedIntegrationTypes.String
	if allowedIntegrationTypes != "" {
		aits := strings.Split(allowedIntegrationTypes, ",")
		var allowedIntegrationTypes []interface{}
		for _, v := range aits {
			allowedIntegrationType := strings.TrimSpace(v)
			if allowedIntegrationType == "" {
				continue
			}
			if allowedIntegrationType == "SECURITY" {
				allowedIntegrationType = "SECURITY INTEGRATIONS"
			}
			allowedIntegrationTypes = append(allowedIntegrationTypes, allowedIntegrationType)
		}
		allowedIntegrationTypesSet := schema.NewSet(schema.HashString, allowedIntegrationTypes)
		if err := d.Set("allowed_integration_types", allowedIntegrationTypesSet); err != nil {
			return err
		}
	}

	allowedAccounts := failoverGroup.AllowedAccounts.String
	if allowedAccounts != "" {
		aa := strings.Split(allowedAccounts, ",")
		var allowedAccounts []interface{}
		for _, v := range aa {
			allowedAccount := strings.TrimSpace(v)
			if allowedAccount == "" {
				continue
			}
			allowedAccounts = append(allowedAccounts, allowedAccount)
		}
		allowedAccountsSet := schema.NewSet(schema.HashString, allowedAccounts)
		if err := d.Set("allowed_accounts", allowedAccountsSet); err != nil {
			return err
		}
	}

	allowedDatabases, err := snowflake.ShowDatabasesInFailoverGroup(name, db)
	if err != nil {
		return fmt.Errorf("error listing databases in failover group %v err = %w", name, err)
	}
	if len(allowedDatabases) > 0 {
		allowedDatabasesInterface := make([]interface{}, len(allowedDatabases))
		for i, v := range allowedDatabases {
			allowedDatabasesInterface[i] = v
		}
		allowedDatabasesSet := schema.NewSet(schema.HashString, allowedDatabasesInterface)
		if err := d.Set("allowed_databases", allowedDatabasesSet); err != nil {
			return err
		}
	} else {
		if err := d.Set("allowed_databases", nil); err != nil {
			return err
		}
	}

	shares, err := snowflake.ShowSharesInFailoverGroup(name, db)
	if err != nil {
		return fmt.Errorf("error listing shares in failover group %v err = %w", name, err)
	}
	if len(shares) > 0 {
		sharesInterface := make([]interface{}, len(shares))
		for i, v := range shares {
			sharesInterface[i] = v
		}
		sharesSet := schema.NewSet(schema.HashString, sharesInterface)
		if err := d.Set("allowed_shares", sharesSet); err != nil {
			return err
		}
	}
	if err := d.Set("allowed_shares", nil); err != nil {
		return err
	}

	return nil
}

// UpdateFailoverGroup implements schema.UpdateFunc.
func UpdateFailoverGroup(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()
	builder := snowflake.NewFailoverGroupBuilder(name)

	if d.HasChange("object_types") {
		_, new := d.GetChange("object_types")
		newObjectTypes := new.(*schema.Set).List()

		var objectTypes []string
		for _, v := range newObjectTypes {
			objectTypes = append(objectTypes, v.(string))
		}
		stmt := builder.ChangeObjectTypes(objectTypes)
		if err := snowflake.Exec(db, stmt); err != nil {
			return fmt.Errorf("error updating object types for failover group %v err = %w", name, err)
		}
	}

	if d.HasChange("allowed_databases") {
		old, new := d.GetChange("allowed_databases")
		oad := old.(*schema.Set).List()
		oldAllowedDatabases := make([]string, len(oad))
		for i, v := range oad {
			oldAllowedDatabases[i] = v.(string)
		}
		nad := new.(*schema.Set).List()
		newAllowedDatabases := make([]string, len(nad))
		for i, v := range nad {
			newAllowedDatabases[i] = v.(string)
		}

		var removedDatabases []string
		for _, v := range oldAllowedDatabases {
			if !slices.Contains(newAllowedDatabases, v) {
				removedDatabases = append(removedDatabases, v)
			}
		}
		if len(removedDatabases) > 0 {
			stmt := builder.RemoveAllowedDatabases(removedDatabases)
			if err := snowflake.Exec(db, stmt); err != nil {
				return fmt.Errorf("error removing allowed databases for failover group %v err = %w", name, err)
			}
		}

		var addedDatabases []string
		for _, v := range newAllowedDatabases {
			if !slices.Contains(oldAllowedDatabases, v) {
				addedDatabases = append(addedDatabases, v)
			}
		}

		if len(addedDatabases) > 0 {
			stmt := builder.AddAllowedDatabases(addedDatabases)
			if err := snowflake.Exec(db, stmt); err != nil {
				return fmt.Errorf("error adding allowed databases for failover group %v err = %w", name, err)
			}
		}
	}

	if d.HasChange("allowed_shares") {
		old, new := d.GetChange("allowed_shares")
		oad := old.(*schema.Set).List()
		oldAllowedShares := make([]string, len(oad))
		for i, v := range oad {
			oldAllowedShares[i] = v.(string)
		}
		nad := new.(*schema.Set).List()
		newAllowedShares := make([]string, len(nad))
		for i, v := range nad {
			newAllowedShares[i] = v.(string)
		}

		var removedShares []string
		for _, v := range oldAllowedShares {
			if !slices.Contains(newAllowedShares, v) {
				removedShares = append(removedShares, v)
			}
		}
		if len(removedShares) > 0 {
			stmt := builder.RemoveAllowedShares(removedShares)
			if err := snowflake.Exec(db, stmt); err != nil {
				return fmt.Errorf("error removing allowed shares for failover group %v err = %w", name, err)
			}
		}

		var addedShares []string
		for _, v := range newAllowedShares {
			if !slices.Contains(oldAllowedShares, v) {
				addedShares = append(addedShares, v)
			}
		}

		if len(addedShares) > 0 {
			stmt := builder.AddAllowedShares(addedShares)
			if err := snowflake.Exec(db, stmt); err != nil {
				return fmt.Errorf("error adding allowed shares for failover group %v err = %w", name, err)
			}
		}
	}

	if d.HasChange("allowed_integration_types") {
		ait := d.Get("allowed_integration_types").(*schema.Set).List()
		allowedIntegrationTypes := make([]string, len(ait))
		for i, v := range ait {
			allowedIntegrationTypes[i] = v.(string)
		}
		stmt := builder.ChangeAllowedIntegrationTypes(allowedIntegrationTypes)
		if err := snowflake.Exec(db, stmt); err != nil {
			return fmt.Errorf("error updating allowed integration types for failover group %v err = %w", name, err)
		}
	}

	if d.HasChange("allowed_accounts") {
		old, new := d.GetChange("allowed_accounts")
		oad := old.(*schema.Set).List()
		oldAllowedAccounts := make([]string, len(oad))
		for i, v := range oad {
			oldAllowedAccounts[i] = v.(string)
		}
		nad := new.(*schema.Set).List()
		newAllowedAccounts := make([]string, len(nad))
		for i, v := range nad {
			newAllowedAccounts[i] = v.(string)
		}

		var removedAccounts []string
		for _, v := range oldAllowedAccounts {
			if !slices.Contains(newAllowedAccounts, v) {
				removedAccounts = append(removedAccounts, v)
			}
		}
		if len(removedAccounts) > 0 {
			stmt := builder.RemoveAllowedAccounts(removedAccounts)
			if err := snowflake.Exec(db, stmt); err != nil {
				return fmt.Errorf("error removing allowed accounts for failover group %v err = %w", name, err)
			}
		}

		var addedAccounts []string
		for _, v := range newAllowedAccounts {
			if !slices.Contains(oldAllowedAccounts, v) {
				addedAccounts = append(addedAccounts, v)
			}
		}

		if len(addedAccounts) > 0 {
			stmt := builder.AddAllowedAccounts(addedAccounts)
			if err := snowflake.Exec(db, stmt); err != nil {
				if !strings.Contains(err.Error(), "Replication already enabled for account") {
					return fmt.Errorf("error adding allowed accounts for failover group %v err = %w", name, err)
				}
			}
		}
	}

	if d.HasChange("replication_schedule") {
		_, new := d.GetChange("replication_schedule")
		replicationSchedule := new.([]interface{})[0].(map[string]interface{})
		log.Printf("[DEBUG] replicationSchedule: %v", replicationSchedule)
		log.Printf("[DEBUG] replicationSchedule[cron]: %v", replicationSchedule["cron"])
		c := replicationSchedule["cron"].([]interface{})
		if len(c) > 0 {
			if len(c) > 0 {
				cron := c[0].(map[string]interface{})
				cronExpression := cron["expression"].(string)

				timeZone := ""
				if v, ok := cron["time_zone"]; ok {
					timeZone = v.(string)
				}
				stmt := builder.ChangeReplicationCronSchedule(cronExpression, timeZone)
				if err := snowflake.Exec(db, stmt); err != nil {
					return fmt.Errorf("error updating replication cron schedule for failover group %v err = %w", name, err)
				}
			}
		} else {
			log.Printf("[DEBUG] replicationSchedule[interval]: %v", replicationSchedule["interval"])
			interval := replicationSchedule["interval"].(int)
			stmt := builder.ChangeReplicationIntervalSchedule(interval)
			if err := snowflake.Exec(db, stmt); err != nil {
				return fmt.Errorf("error updating replication interval schedule for failover group %v err = %w", name, err)
			}
		}
	}

	return ReadFailoverGroup(d, meta)
}

// DeleteFailoverGroup implements schema.DeleteFunc.
func DeleteFailoverGroup(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()
	builder := snowflake.NewFailoverGroupBuilder(name)
	stmt := builder.Drop()
	if err := snowflake.Exec(db, stmt); err != nil {
		return fmt.Errorf("error deleting file format %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
