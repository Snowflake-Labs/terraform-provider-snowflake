package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validDatabasePrivileges = []string{
	"ALL", "CREATE SCHEMA", "IMPORTED PRIVILEGES", "MODIFY", "MONITOR",
	"OWNERSHIP", "REFERENCE_USAGE", "USAGE",
}

var databaseGrantSchema = map[string]*schema.Schema{
	"database_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the database.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validDatabasePrivileges, true),
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
		Description: "Grants privilege to these shares.",
		ForceNew:    true,
	},
}

// DatabaseGrant returns a pointer to the resource representing a database grant
func DatabaseGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabaseGrant,
		Read:   ReadDatabaseGrant,
		Delete: DeleteDatabaseGrant,

		Schema: databaseGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateDatabaseGrant implements schema.CreateFunc
func CreateDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	dbName := data.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(dbName)
	priv := data.Get("privilege").(string)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	// ID format is <db_name>|||<privilege>
	dataIdentifiers := make([][]string, 1)
	dataIdentifiers[0] = make([]string, 2)
	dataIdentifiers[0][0] = dbName
	dataIdentifiers[0][1] = priv
	grantID, err := createGrantID(dataIdentifiers)

	if err != nil {
		return err
	}

	data.SetId(grantID)

	return ReadDatabaseGrant(data, meta)
}

// ReadDatabaseGrant implements schema.ReadFunc
func ReadDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	grantIDArray, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}
	dbName, priv := grantIDArray[0], grantIDArray[1]

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	// IMPORTED PRIVILEGES is not a real resource, so we can't actually verify
	// that it is still there. Just exit for now
	if priv == "IMPORTED PRIVILEGES" {
		return nil
	}

	builder := snowflake.DatabaseGrant(dbName)

	return readGenericGrant(data, meta, builder, false)
}

// DeleteDatabaseGrant implements schema.DeleteFunc
func DeleteDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	dbName := data.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(dbName)

	return deleteGenericGrant(data, meta, builder)
}
