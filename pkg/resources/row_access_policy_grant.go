package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		ValidateFunc: validation.StringInSlice(validRowAccessPoilcyPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
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
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
}

// RowAccessPolicyGrant returns a pointer to the resource representing a row access policy grant.
func RowAccessPolicyGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateRowAccessPolicyGrant,
			Read:   ReadRowAccessPolicyGrant,
			Delete: DeleteRowAccessPolicyGrant,
			Update: UpdateRowAccessPolicyGrant,

			Schema: rowAccessPolicyGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validRowAccessPoilcyPrivileges,
	}
}

// CreateRowAccessPolicyGrant implements schema.CreateFunc.
func CreateRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	var rowAccessPolicyName string
	if name, ok := d.GetOk("row_access_policy_name"); ok {
		rowAccessPolicyName = name.(string)
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.RowAccessPolicyGrant(databaseName, schemaName, rowAccessPolicyName)

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewRowAccessPolicyGrantID(databaseName, schemaName, rowAccessPolicyName, privilege, roles, grantOption)

	d.SetId(grantID.String())

	return ReadRowAccessPolicyGrant(d, meta)
}

// ReadRowAccessPolicyGrant implements schema.ReadFunc.
func ReadRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseRowAccessPolicyGrantID(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	if err := d.Set("row_access_policy_name", grantID.ObjectName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	builder := snowflake.RowAccessPolicyGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	return readGenericGrant(d, meta, rowAccessPolicyGrantSchema, builder, false, validRowAccessPoilcyPrivileges)
}

// DeleteRowAccessPolicyGrant implements schema.DeleteFunc.
func DeleteRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseRowAccessPolicyGrantID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.RowAccessPolicyGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateRowAccessPolicyGrant implements schema.UpdateFunc.
func UpdateRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := parseRowAccessPolicyGrantID(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.RowAccessPolicyGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.WithGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadRowAccessPolicyGrant(d, meta)
}

type RowAccessPolicyGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID		 bool
}

func NewRowAccessPolicyGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, withGrantOption bool) *RowAccessPolicyGrantID {
	return &RowAccessPolicyGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID: false,
	}
}

func (v *RowAccessPolicyGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v❄️%v❄️%v❄️%v❄️%v❄️%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, roles, v.WithGrantOption)
}

func parseRowAccessPolicyGrantID(s string) (*RowAccessPolicyGrantID, error) {
	// is this an old ID format?
	if !strings.Contains(s, "❄️") {
		idParts := strings.Split(s, "|")
		return &RowAccessPolicyGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:          []string{},
			WithGrantOption: idParts[4] == "true",
			IsOldID: true,
		}, nil
	}
	idParts := strings.Split(s, "❄️")
	if len(idParts) != 6 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 6", len(idParts))
	}
	return &RowAccessPolicyGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		Roles:           strings.Split(idParts[4], ","),
		WithGrantOption: idParts[5] == "true",
		IsOldID: false,
	}, nil
}
