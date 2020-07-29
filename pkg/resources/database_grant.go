package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var ValidDatabasePrivileges = newPrivilegeSet(
	privilegeAll,
	privilegeCreateSchema,
	privilegeImportedPrivileges,
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeReferenceUsage,
	privilegeUsage,
)

var databaseGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the database.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(ValidDatabasePrivileges.toList(), true),
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

	grant := &grantID{
		ResourceName: dbName,
		Privilege:    priv,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadDatabaseGrant(data, meta)
}

// ReadDatabaseGrant implements schema.ReadFunc
func ReadDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	err = data.Set("database_name", grantID.ResourceName)
	if err != nil {
		return err
	}
	err = data.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}

	// IMPORTED PRIVILEGES is not a real resource, so we can't actually verify
	// that it is still there. Just exit for now
	if grantID.Privilege == "IMPORTED PRIVILEGES" {
		return nil
	}

	builder := snowflake.DatabaseGrant(grantID.ResourceName)

	return readGenericGrant(data, meta, builder, false, ValidDatabasePrivileges)
}

// DeleteDatabaseGrant implements schema.DeleteFunc
func DeleteDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	dbName := data.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(dbName)

	return deleteGenericGrant(data, meta, builder)
}
