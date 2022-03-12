package resources

import (
	"database/sql"
	"fmt"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func OwnershipGrants() *schema.Resource {
	return &schema.Resource{
		Create: CreateOwnershipGrants,
		Read:   ReadOwnershipGrants,
		Delete: DeleteOwnershipGrants,
		Update: UpdateOwnershipGrants,

		Schema: map[string]*schema.Schema{
			"roles": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Grants ownership to this specified role.",
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the role on which to grant ownership privileges.",
				ForceNew:    true,
			},
			"current_grants": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role.",
				Default:     "COPY",
				ForceNew:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateOwnershipGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roles := expandStringList(data.Get("roles").(*schema.Set).List())
	owner := data.Get("owner").(string)
	currentGrants := data.Get("current_grants").(string)

	if len(roles) == 0 {
		return fmt.Errorf("no roles specified for role grants")
	}

	for _, role := range roles {
		err := grantOwnershipToRole(db, owner, role, currentGrants)
		if err != nil {
			return err
		}
	}

	data.SetId(owner)
	return ReadOwnershipGrants(data, meta)
}

func grantOwnershipToRole(db *sql.DB, role1, role2 string, currentGrants string) error {
	g := snowflake.OwnershipGrant(role1)
	err := snowflake.Exec(db, g.Role(role2, currentGrants).Grant())
	return err
}

func ReadOwnershipGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	owner := data.Id()

	roles := make([]string, 0)

	grants, err := readGrants(db, owner)
	if err != nil {
		return err
	}

	for _, grant := range grants {
		switch grant.GrantedTo.String {
		case "ROLE":
			roles = append(roles, grant.GranteeName.String)
		default:
			return fmt.Errorf("unknown grant type %s", grant.GrantedTo.String)
		}
	}

	err = data.Set("roles", roles)
	if err != nil {
		return err
	}
	err = data.Set("owner", owner)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOwnershipGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roles := expandStringList(data.Get("roles").(*schema.Set).List())
	owner := data.Get("owner").(string)
	currentGrants := data.Get("current_grants").(string)

	for _, role := range roles {
		err := revokeOwnershipFromRole(db, owner, role, currentGrants)
		if err != nil {
			return err
		}
	}

	data.SetId("")
	return nil
}

func revokeOwnershipFromRole(db *sql.DB, role1, role2 string, currentGrants string) error {
	// rg := snowflake.RoleGrant(role1).Role(role2)
	og := snowflake.OwnershipGrant(role1).Role(role2, currentGrants)
	err := snowflake.Exec(db, og.Revoke())
	return err
}

func UpdateOwnershipGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	owner := data.Get("owner").(string)
	currentGrants := data.Get("current_grants").(string)

	x := func(resource string, grant func(db *sql.DB, role string, target string, currentGrants string) error, revoke func(db *sql.DB, role string, target string, currentGrants string) error) error {
		o, n := data.GetChange(resource)

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

		for _, role := range remove {
			err := revoke(db, owner, role, currentGrants)
			if err != nil {
				return err
			}
		}
		for _, role := range add {
			err := grant(db, owner, role, currentGrants)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := x("roles", grantOwnershipToRole, revokeOwnershipFromRole)
	if err != nil {
		return err
	}

	return ReadOwnershipGrants(data, meta)
}
