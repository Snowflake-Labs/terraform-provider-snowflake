package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validIntegrationPrivileges = NewPrivilegeSet(
	privilegeUsage,
	privilegeOwnership,
)

var integrationGrantSchema = map[string]*schema.Schema{
	"integration_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the integration; must be unique for your account.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the integration.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validIntegrationPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
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

// IntegrationGrant returns a pointer to the resource representing a integration grant.
func IntegrationGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateIntegrationGrant,
			Read:   ReadIntegrationGrant,
			Delete: DeleteIntegrationGrant,
			Update: UpdateIntegrationGrant,

			Schema: integrationGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validIntegrationPrivileges,
	}
}

// CreateIntegrationGrant implements schema.CreateFunc.
func CreateIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	integrationName := d.Get("integration_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.IntegrationGrant(integrationName)
	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewIntegrationGrantID(integrationName, privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadIntegrationGrant(d, meta)
}

// ReadIntegrationGrant implements schema.ReadFunc.
func ReadIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseIntegrationGrantID(d.Id())
	if err != nil {
		return err
	}
	if err := d.Set("integration_name", grantID.ObjectName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	builder := snowflake.IntegrationGrant(grantID.ObjectName)

	return readGenericGrant(d, meta, integrationGrantSchema, builder, false, validIntegrationPrivileges)
}

// DeleteIntegrationGrant implements schema.DeleteFunc.
func DeleteIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseIntegrationGrantID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.IntegrationGrant(grantID.ObjectName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateIntegrationGrant implements schema.UpdateFunc.
func UpdateIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := ParseIntegrationGrantID(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.IntegrationGrant(grantID.ObjectName)

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
	return ReadIntegrationGrant(d, meta)
}

type IntegrationGrantID struct {
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewIntegrationGrantID(objectName string, privilege string, roles []string, withGrantOption bool) *IntegrationGrantID {
	return &IntegrationGrantID{
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *IntegrationGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v", v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseIntegrationGrantID(s string) (*IntegrationGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		var roles []string
		var withGrantOption bool
		if len(idParts) == 6 {
			withGrantOption = idParts[5] == "true"
			roles = helpers.SplitStringToSlice(idParts[4], ",")
		} else {
			withGrantOption = idParts[4] == "true"
		}
		return &IntegrationGrantID{
			ObjectName:      idParts[0],
			Privilege:       idParts[3],
			Roles:           roles,
			WithGrantOption: withGrantOption,
			IsOldID:         true,
		}, nil
	}
	idParts := helpers.SplitStringToSlice(s, "|")
	if len(idParts) < 4 {
		idParts = helpers.SplitStringToSlice(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 4 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 4", len(idParts))
	}
	return &IntegrationGrantID{
		ObjectName:      idParts[0],
		Privilege:       idParts[1],
		WithGrantOption: idParts[2] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[3], ","),
		IsOldID:         false,
	}, nil
}
