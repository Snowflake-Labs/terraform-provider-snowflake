package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validTablePrivileges = newPrivilegeSet(
	privilegeSelect,
	privilegeInsert,
	privilegeUpdate,
	privilegeDelete,
	privilegeTruncate,
	privilegeReferences,
)

var tableGrantSchema = map[string]*schema.Schema{
	"table_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the table on which to grant privileges immediately (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"schema_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future tables on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future tables on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future table.",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(validTablePrivileges.toList(), true),
		ForceNew:     true,
	},
	"roles": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"shares": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"on_future": &schema.Schema{
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future tables in the given schema. When this is true and no schema_name is provided apply this grant on all future tables in the given database. The table_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"table_name", "shares"},
	},
}

// TableGrant returns a pointer to the resource representing a Table grant
func TableGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateTableGrant,
		Read:   ReadTableGrant,
		Delete: DeleteTableGrant,

		Schema: tableGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateTableGrant implements schema.CreateFunc
func CreateTableGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		tableName  string
		schemaName string
	)
	if _, ok := data.GetOk("table_name"); ok {
		tableName = data.Get("table_name").(string)
	} else {
		tableName = ""
	}
	if _, ok := data.GetOk("schema_name"); ok {
		schemaName = data.Get("schema_name").(string)
	} else {
		schemaName = ""
	}
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	onFuture := data.Get("on_future").(bool)

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

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	// table_name is empty when on_future = true
	grantID := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		Privilege:    priv,
	}
	if !onFuture {
		grantID.ObjectName = tableName
	}

	dataIDInput, err := grantID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)
	return ReadTableGrant(data, meta)
}

// ReadTableGrant implements schema.ReadFunc
func ReadTableGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	tableName := grantID.ObjectName
	priv := grantID.Privilege
	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	onFuture := false
	if tableName == "" {
		onFuture = true
	}
	err = data.Set("table_name", tableName)
	if err != nil {
		return err
	}
	err = data.Set("on_future", onFuture)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.TableGrant(dbName, schemaName, tableName)
	}

	return readGenericGrant(data, meta, builder, onFuture, validTablePrivileges)
}

// DeleteTableGrant implements schema.DeleteFunc
func DeleteTableGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
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
	return deleteGenericGrant(data, meta, builder)
}
