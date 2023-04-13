package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validTaskPrivileges = NewPrivilegeSet(
	privilegeMonitor,
	privilegeOperate,
	privilegeOwnership,
)

var taskGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future tasks on which to grant privileges.",
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future tasks in the given schema. When this is true and no schema_name is provided apply this grant on all future tasks in the given database. The task_name field must be unset in order to use on_future.",
		Default:     false,
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future task.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validTaskPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Grants privilege to these roles.",
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future tasks on which to grant privileges.",
		ForceNew:    true,
	},
	"task_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the task on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// TaskGrant returns a pointer to the resource representing a task grant.
func TaskGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateTaskGrant,
			Read:   ReadTaskGrant,
			Delete: DeleteTaskGrant,
			Update: UpdateTaskGrant,

			Schema: taskGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 7 {
						return nil, fmt.Errorf("unexpected format of ID (%v), expected database_name|schema_name|task_name|privilege|with_grant_option|on_future|roles", d.Id())
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if err := d.Set("task_name", parts[2]); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", parts[3]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[4])); err != nil {
						return nil, err
					}
					if err := d.Set("on_future", helpers.StringToBool(parts[5])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[6])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validTaskPrivileges,
	}
}

// CreateTaskGrant implements schema.CreateFunc.
func CreateTaskGrant(d *schema.ResourceData, meta interface{}) error {
	var taskName string
	if name, ok := d.GetOk("task_name"); ok {
		taskName = name.(string)
	}
	if err := d.Set("task_name", taskName); err != nil {
		return err
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (taskName == "") && !onFuture {
		return errors.New("task_name must be set unless on_future is true")
	}
	if (taskName != "") && onFuture {
		return errors.New("task_name must be empty if on_future is true")
	}
	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTaskGrant(databaseName, schemaName)
	} else {
		builder = snowflake.TaskGrant(databaseName, schemaName, taskName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}
	grantID := helpers.SnowflakeID(databaseName, schemaName, taskName, privilege, withGrantOption, onFuture, roles)
	d.SetId(grantID)

	return ReadTaskGrant(d, meta)
}

// ReadTaskGrant implements schema.ReadFunc.
func ReadTaskGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	taskName := d.Get("task_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTaskGrant(databaseName, schemaName)
	} else {
		builder = snowflake.TaskGrant(databaseName, schemaName, taskName)
	}
	// TODO
	onAll := false

	err := readGenericGrant(d, meta, taskGrantSchema, builder, onFuture, onAll, validTaskPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.SnowflakeID(databaseName, schemaName, taskName, privilege, withGrantOption, onFuture, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteTaskGrant implements schema.DeleteFunc.
func DeleteTaskGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	taskName := d.Get("task_name").(string)
	onFuture := d.Get("on_future").(bool)

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTaskGrant(databaseName, schemaName)
	} else {
		builder = snowflake.TaskGrant(databaseName, schemaName, taskName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateTaskGrant implements schema.UpdateFunc.
func UpdateTaskGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update and we're done
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}

	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	taskName := d.Get("task_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTaskGrant(databaseName, schemaName)
	} else {
		builder = snowflake.TaskGrant(databaseName, schemaName, taskName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, privilege, withGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadTaskGrant(d, meta)
}
