package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validTablePrivileges = NewPrivilegeSet(
	privilegeSelect,
	privilegeInsert,
	privilegeUpdate,
	privilegeDelete,
	privilegeTruncate,
	privilegeReferences,
)

var tableGrantSchema = map[string]*schema.Schema{
	"table_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the table on which to grant privileges immediately (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future tables on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future tables on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future table.",
		Default:      privilegeSelect.String(),
		ValidateFunc: validation.StringInSlice(validTablePrivileges.ToList(), true),
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
		Description: "Grants privilege to these shares (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future tables in the given schema. When this is true and no schema_name is provided apply this grant on all future tables in the given database. The table_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"table_name", "shares"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// TableGrant returns a pointer to the resource representing a Table grant
func TableGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateTableGrant,
			Read:   ReadTableGrant,
			Delete: DeleteTableGrant,

			Schema: tableGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validTablePrivileges,
	}
}

// CreateTableGrant implements schema.CreateFunc
func CreateTableGrant(d *schema.ResourceData, meta interface{}) error {
	var (
		tableName  string
		schemaName string
	)
	if _, ok := d.GetOk("table_name"); ok {
		tableName = d.Get("table_name").(string)
	} else {
		tableName = ""
	}
	if _, ok := d.GetOk("schema_name"); ok {
		schemaName = d.Get("schema_name").(string)
	} else {
		schemaName = ""
	}
	dbName := d.Get("database_name").(string)
	priv := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)

	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	if (tableName == "") && !onFuture {
		return errors.New("table_name must be set unless on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.TableGrant(dbName, schemaName, tableName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	// table_name is empty when on_future = true
	grantID := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	if !onFuture {
		grantID.ObjectName = tableName
	}

	dataIDInput, err := grantID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)
	return ReadTableGrant(d, meta)
}

// ReadTableGrant implements schema.ReadFunc
func ReadTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	tableName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	onFuture := false
	if tableName == "" {
		onFuture = true
	}
	err = d.Set("table_name", tableName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", onFuture)
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
	if onFuture {
		builder = snowflake.FutureTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.TableGrant(dbName, schemaName, tableName)
	}

	return readGenericGrant(d, meta, tableGrantSchema, builder, onFuture, validTablePrivileges)
}

// DeleteTableGrant implements schema.DeleteFunc
func DeleteTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	tableName := grantID.ObjectName
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	onFuture := false
	if tableName == "" {
		onFuture = true
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.TableGrant(dbName, schemaName, tableName)
	}
	return deleteGenericGrant(d, meta, builder)
}
