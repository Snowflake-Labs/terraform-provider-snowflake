package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func RoleOwnershipGrants() *schema.Resource {
	return &schema.Resource{
		Create: CreateRoleOwnershipGrants,
		Read:   ReadRoleOwnershipGrants,
		Delete: DeleteRoleOwnershipGrants,
		Update: UpdateRoleOwnershipGrants,

		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The name of the role to grant ownership. Please ensure that the role that terraform is using is granted access.",
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					return snowflake.ValidateIdentifier(val)
				},
			},
			"roles": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Grants role ownership to these specified roles.",
			},
			"users": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Grants role ownership to these specified users.",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := d.Get("role_name").(string)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	users := expandStringList(d.Get("users").(*schema.Set).List())

	if len(roles) == 0 && len(users) == 0 {
		return fmt.Errorf("no users or roles specified for role grants")
	}

	grant := &grantID{
		ResourceName: roleName,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	d.SetId(dataIDInput)

	if err != nil {
		return errors.Wrap(err, "error creating role grant")
	}
	for _, role := range roles {
		err := grantRoleOwnershipToRole(db, roleName, role)
		if err != nil {
			return err
		}
	}

	for _, user := range users {
		err := grantRoleOwnershipToUser(db, roleName, user)
		if err != nil {
			return err
		}
	}

	return ReadRoleOwnershipGrants(d, meta)
}

func grantRoleOwnershipToRole(db *sql.DB, role1, role2 string) error {
	g := snowflake.RoleOwnershipGrant(role1)
	err := snowflake.Exec(db, g.Role(role2).Grant())
	return err
}

func grantRoleOwnershipToUser(db *sql.DB, role1, user string) error {
	g := snowflake.RoleOwnershipGrant(role1)
	err := snowflake.Exec(db, g.User(user).Grant())
	return err
}

func ReadRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	log.Println(d.Id())
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	roleName := grantID.ResourceName

	tfRoles := expandStringList(d.Get("roles").(*schema.Set).List())
	tfUsers := expandStringList(d.Get("users").(*schema.Set).List())

	roles := make([]string, 0)
	users := make([]string, 0)

	grants, err := readGrants(db, roleName)
	if err != nil {
		return err
	}

	for _, grant := range grants {
		switch grant.GrantedTo.String {
		case "ROLE":
			for _, tfRole := range tfRoles {
				if tfRole == grant.GranteeName.String {
					roles = append(roles, grant.GranteeName.String)
				}
			}
		case "USER":
			for _, tfUser := range tfUsers {
				if tfUser == grant.GranteeName.String {
					users = append(users, grant.GranteeName.String)
				}
			}
		default:
			return fmt.Errorf("unknown grant type %s", grant.GrantedTo.String)
		}
	}

	err = d.Set("role_name", roleName)
	if err != nil {
		return err
	}
	err = d.Set("roles", roles)
	if err != nil {
		return err
	}
	err = d.Set("users", users)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := d.Get("role_name").(string)

	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	users := expandStringList(d.Get("users").(*schema.Set).List())

	for _, role := range roles {
		err := revokeRoleFromRole(db, roleName, role)
		if err != nil {
			return err
		}
	}

	for _, user := range users {
		err := revokeRoleFromUser(db, roleName, user)
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}

func revokeRoleOwnershipFromRole(db *sql.DB, role1, role2 string) error {
	rg := snowflake.RoleOwnershipGrant(role1).Role(role2)
	err := snowflake.Exec(db, rg.Revoke())
	return err
}

func revokeRoleOwnershipFromUser(db *sql.DB, role1, user string) error {
	rg := snowflake.RoleOwnershipGrant(role1).User(user)
	err := snowflake.Exec(db, rg.Revoke())
	return err
}

func UpdateRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := d.Get("role_name").(string)

	x := func(resource string, grant func(db *sql.DB, role string, target string) error, revoke func(db *sql.DB, role string, target string) error) error {
		o, n := d.GetChange(resource)

		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove := expandStringList(os.Difference(ns).List())
		add := expandStringList(ns.Difference(os).List())

		for _, user := range remove {
			err := revoke(db, roleName, user)
			if err != nil {
				return err
			}
		}
		for _, user := range add {
			err := grant(db, roleName, user)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := x("users", grantRoleOwnershipToUser, revokeRoleOwnershipFromUser)
	if err != nil {
		return err
	}

	err = x("roles", grantRoleOwnershipToRole, revokeRoleOwnershipFromRole)
	if err != nil {
		return err
	}

	return ReadRoleOwnershipGrants(d, meta)
}
