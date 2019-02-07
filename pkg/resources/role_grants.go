package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jmoiron/sqlx"
)

func RoleGrants() *schema.Resource {
	return &schema.Resource{
		Create: CreateRoleGrants,
		Read:   ReadRoleGrants,
		Delete: DeleteRoleGrants,
		Update: UpdateRoleGrants,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"role_name": &schema.Schema{
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The name of the role we are granting.",
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					return snowflake.ValidateIdentifier(val)
				},
			},
			"roles": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Grants role to this specified role.",
			},
			"users": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Grants role to this specified user.",
			},
		},
		// TODO
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
	}
}

func CreateRoleGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	roleName := data.Get("role_name").(string)
	roles := expandStringList(data.Get("roles").(*schema.Set).List())
	users := expandStringList(data.Get("users").(*schema.Set).List())

	if len(roles) == 0 && len(users) == 0 {
		return fmt.Errorf("no users or roles specified for role grants")
	}

	for _, role := range roles {
		err := grantRoleToRole(db, roleName, role)
		if err != nil {
			return err
		}
	}

	for _, user := range users {
		err := grantRoleToUser(db, roleName, user)
		if err != nil {
			return err
		}
	}
	data.SetId(name)
	return ReadRoleGrants(data, meta)
}

func grantRoleToRole(db *sql.DB, role1, role2 string) error {
	g := snowflake.RoleGrant(role1)
	err := DBExec(db, g.Role(role2).Grant())
	return err
}

func grantRoleToUser(db *sql.DB, role1, user string) error {
	g := snowflake.RoleGrant(role1)
	err := DBExec(db, g.User(user).Grant())
	return err
}

// Trying out this new way of reading rows. If it works well, will move to pkg/snowflake and apply to other
//  types too.
type grant struct {
	CreatedOn   sql.RawBytes   `db:"created_on"`
	Role        sql.NullString `db:"role"`
	GrantedTo   sql.NullString `db:"granted_to"`
	GranteeName sql.NullString `db:"grantee_name"`
	Grantedby   sql.NullString `db:"granted_by"`
}

func ReadRoleGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := data.Get("role_name").(string)

	roles := make([]string, 0)
	users := make([]string, 0)

	grants, err := readGrants(db, roleName)
	if err != nil {
		return err
	}

	for _, grant := range grants {
		switch grant.GrantedTo.String {
		case "ROLE":
			roles = append(roles, grant.GranteeName.String)
		case "USER":
			users = append(users, grant.GranteeName.String)
		default:
			return fmt.Errorf("unknown grant type %s", grant.GrantedTo.String)
		}
	}

	err = data.Set("roles", roles)
	if err != nil {
		return err
	}
	err = data.Set("users", users)
	if err != nil {
		return err
	}

	return nil
}

func readGrants(db *sql.DB, roleName string) ([]*grant, error) {
	sdb := sqlx.NewDb(db, "snowflake")

	stmt := fmt.Sprintf(`SHOW GRANTS OF ROLE "%s"`, roleName)
	log.Printf("[DEBUG] stmt %s", stmt)
	rows, err := sdb.Queryx(stmt)

	if err != nil {
		return nil, err
	}

	grants := make([]*grant, 0)
	for rows.Next() {
		g := &grant{}
		err = rows.StructScan(g)
		if err != nil {
			return nil, err
		}
		grants = append(grants, g)

	}
	return grants, nil
}

func DeleteRoleGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := data.Get("role_name").(string)

	roles := expandStringList(data.Get("roles").(*schema.Set).List())
	users := expandStringList(data.Get("users").(*schema.Set).List())

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

	data.SetId("")
	return nil
}

func revokeRoleFromRole(db *sql.DB, role1, role2 string) error {
	rg := snowflake.RoleGrant(role1).Role(role2)
	err := DBExec(db, rg.Revoke())
	return err
}

func revokeRoleFromUser(db *sql.DB, role1, user string) error {
	rg := snowflake.RoleGrant(role1).User(user)
	err := DBExec(db, rg.Revoke())
	return err
}

func UpdateRoleGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := data.Get("role_name").(string)

	x := func(resource string, grant func(db *sql.DB, role string, target string) error, revoke func(db *sql.DB, role string, target string) error) error {
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

	err := x("users", grantRoleToUser, revokeRoleFromUser)
	if err != nil {
		return err
	}

	err = x("roles", grantRoleToRole, revokeRoleFromRole)
	if err != nil {
		return err
	}

	return ReadRoleGrants(data, meta)
}

// borrowed from https://github.com/terraform-providers/terraform-provider-aws/blob/master/aws/structure.go#L924:6
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
