package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var validMaskingPoilcyPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeApply,
)

var maskingPolicyGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the masking policy on which to grant privileges.",
		ForceNew:    true,
	},
	"masking_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the masking policy on which to grant privileges immediately.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the masking policy.",
		Default:      "APPLY",
		ValidateFunc: validation.ValidatePrivilege(validMaskingPoilcyPrivileges.ToList(), true),
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
		Description: "The name of the schema containing the masking policy on which to grant privileges.",
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

// MaskingPolicyGrant returns a pointer to the resource representing a masking policy grant
func MaskingPolicyGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateMaskingPolicyGrant,
			Read:   ReadMaskingPolicyGrant,
			Delete: DeleteMaskingPolicyGrant,

			Schema: maskingPolicyGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validMaskingPoilcyPrivileges,
	}
}

// CreateMaskingPolicyGrant implements schema.CreateFunc
func CreateMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	var maskingPolicyName string
	if name, ok := d.GetOk("masking_policy_name"); ok {
		maskingPolicyName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.MaskingPolicyGrant(dbName, schemaName, maskingPolicyName)

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   maskingPolicyName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadMaskingPolicyGrant(d, meta)
}

// ReadMaskingPolicyGrant implements schema.ReadFunc
func ReadMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	maskingPolicyName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	err = d.Set("masking_policy_name", maskingPolicyName)
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

	builder := snowflake.MaskingPolicyGrant(dbName, schemaName, maskingPolicyName)

	return readGenericGrant(d, meta, maskingPolicyGrantSchema, builder, false, validMaskingPoilcyPrivileges)
}

// DeleteMaskingPolicyGrant implements schema.DeleteFunc
func DeleteMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	maskingPolicyName := grantID.ObjectName

	builder := snowflake.MaskingPolicyGrant(dbName, schemaName, maskingPolicyName)

	return deleteGenericGrant(d, meta, builder)
}
