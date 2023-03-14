package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validUserPrivileges = NewPrivilegeSet(
	privilegeMonitor,
	privilegeOwnership,
)

var userGrantSchema = map[string]*schema.Schema{
	"user_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the user on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The privilege to grant on the user.",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(validUserPrivileges.ToList(), true),
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
	},
}

// UserGrant returns a pointer to the resource representing a user grant.
func UserGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateUserGrant,
			Read:   ReadUserGrant,
			Delete: DeleteUserGrant,
			Update: UpdateUserGrant,

			Schema: userGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					grantID, err := ParseUserGrantID(d.Id())
					if err != nil {
						return nil, err
					}
					if err := d.Set("user_name", grantID.ObjectName); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", grantID.Privilege); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
						return nil, err
					}
					if err := d.Set("roles", grantID.Roles); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validUserPrivileges,
	}
}

// CreateUserGrant implements schema.CreateFunc.
func CreateUserGrant(d *schema.ResourceData, meta interface{}) error {
	userName := d.Get("user_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	builder := snowflake.UserGrant(userName)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewUserGrantID(userName, privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadUserGrant(d, meta)
}

// ReadUserGrant implements schema.ReadFunc.
func ReadUserGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseUserGrantID(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("user_name", grantID.ObjectName); err != nil {
		return err
	}

	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}

	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	builder := snowflake.UserGrant(grantID.ObjectName)

	return readGenericGrant(d, meta, userGrantSchema, builder, false, validUserPrivileges)
}

// DeleteUserGrant implements schema.DeleteFunc.
func DeleteUserGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseUserGrantID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.UserGrant(grantID.ObjectName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateUserGrant implements schema.UpdateFunc.
func UpdateUserGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update is roles. if nothing changed,
	// nothing to update and we're done.
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd, rolesToRevoke := changeDiff(d, "roles")

	grantID, err := ParseUserGrantID(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.UserGrant(grantID.ObjectName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		rolesToRevoke,
		nil,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		grantID.WithGrantOption,
		rolesToAdd,
		nil,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadUserGrant(d, meta)
}

type UserGrantID struct {
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewUserGrantID(objectName string, privilege string, roles []string, withGrantOption bool) *UserGrantID {
	return &UserGrantID{
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
	}
}

func (v *UserGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v", v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseUserGrantID(s string) (*UserGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &UserGrantID{
			ObjectName:      idParts[0],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			WithGrantOption: idParts[5] == "true",
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
	return &UserGrantID{
		ObjectName:      idParts[0],
		Privilege:       idParts[1],
		WithGrantOption: idParts[2] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[3], ","),
		IsOldID:         false,
	}, nil
}
