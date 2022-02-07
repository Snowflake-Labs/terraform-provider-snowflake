package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validExternalTablePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeReferences,
	privilegeSelect,
)

var externalTableGrantSchema = map[string]*schema.Schema{
	"external_table_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the external table on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current or future external tables on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future external tables on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future external table.",
		Default:      "SELECT",
		ValidateFunc: validation.ValidatePrivilege(validExternalTablePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"shares": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares (only valid if on_future is false).",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future external tables in the given schema. When this is true and no schema_name is provided apply this grant on all future external tables in the given database. The external_table_name and shares fields must be unset in order to use on_future.",
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

// ExternalTableGrant returns a pointer to the resource representing a external table grant
func ExternalTableGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateExternalTableGrant,
			Read:   ReadExternalTableGrant,
			Delete: DeleteExternalTableGrant,

			Schema: externalTableGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validExternalTablePrivileges,
	}
}

// CreateExternalTableGrant implements schema.CreateFunc
func CreateExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	var externalTableName string
	if name, ok := d.GetOk("external_table_name"); ok {
		externalTableName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futureExternalTables := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (externalTableName == "") && !futureExternalTables {
		return errors.New("external_table_name must be set unless on_future is true.")
	}
	if (externalTableName != "") && futureExternalTables {
		return errors.New("external_table_name must be empty if on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if futureExternalTables {
		builder = snowflake.FutureExternalTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.ExternalTableGrant(dbName, schemaName, externalTableName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   externalTableName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadExternalTableGrant(d, meta)
}

// ReadExternalTableGrant implements schema.ReadFunc
func ReadExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	externalTableName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureExternalTablesEnabled := false
	if externalTableName == "" {
		futureExternalTablesEnabled = true
	}
	err = d.Set("external_table_name", externalTableName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureExternalTablesEnabled)
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
	if futureExternalTablesEnabled {
		builder = snowflake.FutureExternalTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.ExternalTableGrant(dbName, schemaName, externalTableName)
	}

	return readGenericGrant(d, meta, externalTableGrantSchema, builder, futureExternalTablesEnabled, validExternalTablePrivileges)
}

// DeleteExternalTableGrant implements schema.DeleteFunc
func DeleteExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	externalTableName := grantID.ObjectName

	futureExternalTables := (externalTableName == "")

	var builder snowflake.GrantBuilder
	if futureExternalTables {
		builder = snowflake.FutureExternalTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.ExternalTableGrant(dbName, schemaName, externalTableName)
	}
	return deleteGenericGrant(d, meta, builder)
}
