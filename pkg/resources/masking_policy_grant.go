package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		ValidateFunc: validation.StringInSlice(validMaskingPoilcyPrivileges.ToList(), true),
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
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
}

// MaskingPolicyGrant returns a pointer to the resource representing a masking policy grant.
func MaskingPolicyGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateMaskingPolicyGrant,
			Read:   ReadMaskingPolicyGrant,
			Delete: DeleteMaskingPolicyGrant,
			Update: UpdateMaskingPolicyGrant,

			Schema: maskingPolicyGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validMaskingPoilcyPrivileges,
	}
}

// CreateMaskingPolicyGrant implements schema.CreateFunc.
func CreateMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	var maskingPolicyName string
	if name, ok := d.GetOk("masking_policy_name"); ok {
		maskingPolicyName = name.(string)
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.MaskingPolicyGrant(databaseName, schemaName, maskingPolicyName)
	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewMaskingPolicyGrantID(databaseName, schemaName, maskingPolicyName, privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadMaskingPolicyGrant(d, meta)
}

// ReadMaskingPolicyGrant implements schema.ReadFunc.
func ReadMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseMaskingPolicyGrantID(d.Id())
	if err != nil {
		return err
	}
	if err := d.Set("roles", grantID.Roles); err != nil {
		return err
	}
	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	if err := d.Set("masking_policy_name", grantID.ObjectName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	builder := snowflake.MaskingPolicyGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	return readGenericGrant(d, meta, maskingPolicyGrantSchema, builder, false, validMaskingPoilcyPrivileges)
}

// DeleteMaskingPolicyGrant implements schema.DeleteFunc.
func DeleteMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseMaskingPolicyGrantID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.MaskingPolicyGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateMaskingPolicyGrant implements schema.UpdateFunc.
func UpdateMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := ParseMaskingPolicyGrantID(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.MaskingPolicyGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

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
	return ReadMaskingPolicyGrant(d, meta)
}

type MaskingPolicyGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewMaskingPolicyGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, withGrantOption bool) *MaskingPolicyGrantID {
	return &MaskingPolicyGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *MaskingPolicyGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseMaskingPolicyGrantID(s string) (*MaskingPolicyGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &MaskingPolicyGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			WithGrantOption: idParts[5] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 6 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 6 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 6", len(idParts))
	}
	return &MaskingPolicyGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
		WithGrantOption: idParts[4] == "true",
		IsOldID:         false,
	}, nil
}
