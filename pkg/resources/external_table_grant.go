package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validExternalTablePrivileges = newPrivilegeSet(
	privilegeSelect,
)

var externalTableGrantSchema = map[string]*schema.Schema{
	"external_table_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the externalTable on which to grant privileges immediately (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future externalTables on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future externalTables on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future externalTable.",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(validExternalTablePrivileges.toList(), true),
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
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future externalTables in the given schema. When this is true and no schema_name is provided apply this grant on all future externalTables in the given database. The external_table_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"external_table_name", "shares"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// ExternalTableGrant returns a pointer to the resource representing a ExternalTable grant
func ExternalTableGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateExternalTableGrant,
		Read:   ReadExternalTableGrant,
		Delete: DeleteExternalTableGrant,

		Schema: externalTableGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateExternalTableGrant implements schema.CreateFunc
func CreateExternalTableGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		externalTableName string
		schemaName        string
	)
	if _, ok := data.GetOk("external_table_name"); ok {
		externalTableName = data.Get("external_table_name").(string)
	} else {
		externalTableName = ""
	}
	if _, ok := data.GetOk("schema_name"); ok {
		schemaName = data.Get("schema_name").(string)
	} else {
		schemaName = ""
	}
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	onFuture := data.Get("on_future").(bool)
	grantOption := data.Get("with_grant_option").(bool)

	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	if (externalTableName == "") && !onFuture {
		return errors.New("external_table_name must be set unless on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureExternalTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.ExternalTableGrant(dbName, schemaName, externalTableName)
	}

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	// external_table_name is empty when on_future = true
	grantID := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	if !onFuture {
		grantID.ObjectName = externalTableName
	}

	dataIDInput, err := grantID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)
	return ReadExternalTableGrant(data, meta)
}

// ReadExternalTableGrant implements schema.ReadFunc
func ReadExternalTableGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	externalTableName := grantID.ObjectName
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
	if externalTableName == "" {
		onFuture = true
	}
	err = data.Set("external_table_name", externalTableName)
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
	err = data.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureExternalTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.ExternalTableGrant(dbName, schemaName, externalTableName)
	}

	return readGenericGrant(data, meta, builder, onFuture, validExternalTablePrivileges)
}

// DeleteExternalTableGrant implements schema.DeleteFunc
func DeleteExternalTableGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}

	externalTableName := grantID.ObjectName
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	onFuture := false
	if externalTableName == "" {
		onFuture = true
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureExternalTableGrant(dbName, schemaName)
	} else {
		builder = snowflake.ExternalTableGrant(dbName, schemaName, externalTableName)
	}
	return deleteGenericGrant(data, meta, builder)
}
