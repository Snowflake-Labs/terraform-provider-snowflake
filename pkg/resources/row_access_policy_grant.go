package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var validRowAccessPoilcyPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeApply,
)

var rowAccessPolicyGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the row access policy on which to grant privileges.",
		ForceNew:    true,
	},
	"row_access_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the row access policy on which to grant privileges immediately.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the row access policy.",
		Default:      "APPLY",
		ValidateFunc: validation.ValidatePrivilege(validRowAccessPoilcyPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the row access policy on which to grant privileges.",
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

// RowAccessPolicyGrant returns a pointer to the resource representing a row access policy grant
func RowAccessPolicyGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateRowAccessPolicyGrant,
			Read:   ReadRowAccessPolicyGrant,
			Delete: DeleteRowAccessPolicyGrant,

			Schema: rowAccessPolicyGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validRowAccessPoilcyPrivileges,
	}
}

// CreateRowAccessPolicyGrant implements schema.CreateFunc
func CreateRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	var rowAccessPolicyName string
	if name, ok := d.GetOk("row_access_policy_name"); ok {
		rowAccessPolicyName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.RowAccessPolicyGrant(dbName, schemaName, rowAccessPolicyName)

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   rowAccessPolicyName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadRowAccessPolicyGrant(d, meta)
}

// ReadRowAccessPolicyGrant implements schema.ReadFunc
func ReadRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	rowAccessPolicyName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	err = d.Set("row_access_policy_name", rowAccessPolicyName)
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

	builder := snowflake.RowAccessPolicyGrant(dbName, schemaName, rowAccessPolicyName)

	return readGenericGrant(d, meta, rowAccessPolicyGrantSchema, builder, false, validRowAccessPoilcyPrivileges)
}

// DeleteRowAccessPolicyGrant implements schema.DeleteFunc
func DeleteRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	rowAccessPolicyName := grantID.ObjectName

	builder := snowflake.RowAccessPolicyGrant(dbName, schemaName, rowAccessPolicyName)

	return deleteGenericGrant(d, meta, builder)
}
