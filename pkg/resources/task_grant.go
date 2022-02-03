package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validTaskPrivileges = NewPrivilegeSet(
	privilegeMonitor,
	privilegeOperate,
	privilegeOwnership,
)

var taskGrantSchema = map[string]*schema.Schema{
	"task_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the task on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current or future tasks on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future tasks on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future task.",
		Default:      "USAGE",
		ValidateFunc: validation.ValidatePrivilege(validTaskPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future tasks in the given schema. When this is true and no schema_name is provided apply this grant on all future tasks in the given database. The task_name field must be unset in order to use on_future.",
		Default:     false,
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

// TaskGrant returns a pointer to the resource representing a task grant
func TaskGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateTaskGrant,
			Read:   ReadTaskGrant,
			Delete: DeleteTaskGrant,

			Schema: taskGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validTaskPrivileges,
	}
}

// CreateTaskGrant implements schema.CreateFunc
func CreateTaskGrant(d *schema.ResourceData, meta interface{}) error {
	var taskName string
	if name, ok := d.GetOk("task_name"); ok {
		taskName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futureTasks := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (taskName == "") && !futureTasks {
		return errors.New("task_name must be set unless on_future is true.")
	}
	if (taskName != "") && futureTasks {
		return errors.New("task_name must be empty if on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if futureTasks {
		builder = snowflake.FutureTaskGrant(dbName, schemaName)
	} else {
		builder = snowflake.TaskGrant(dbName, schemaName, taskName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   taskName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadTaskGrant(d, meta)
}

// ReadTaskGrant implements schema.ReadFunc
func ReadTaskGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	taskName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureTasksEnabled := false
	if taskName == "" {
		futureTasksEnabled = true
	}
	err = d.Set("task_name", taskName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureTasksEnabled)
	if err != nil {
		return err
	}
	err = d.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureTasksEnabled {
		builder = snowflake.FutureTaskGrant(dbName, schemaName)
	} else {
		builder = snowflake.TaskGrant(dbName, schemaName, taskName)
	}

	return readGenericGrant(d, meta, taskGrantSchema, builder, futureTasksEnabled, validTaskPrivileges)
}

// DeleteTaskGrant implements schema.DeleteFunc
func DeleteTaskGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	taskName := grantID.ObjectName

	futureTasks := (taskName == "")

	var builder snowflake.GrantBuilder
	if futureTasks {
		builder = snowflake.FutureTaskGrant(dbName, schemaName)
	} else {
		builder = snowflake.TaskGrant(dbName, schemaName, taskName)
	}
	return deleteGenericGrant(d, meta, builder)
}
