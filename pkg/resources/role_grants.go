package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jmoiron/sqlx"
)

func RoleGrants() *schema.Resource {
	return &schema.Resource{
		Create: CreateRoleGrants,
		Read:   ReadRoleGrants,
		Delete: DeleteRoleGrants,
		Update: UpdateRoleGrants,

		Schema: map[string]*schema.Schema{
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
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Grants role to this specified role.",
			},
			"users": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Grants role to this specified user.",
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateRoleGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := data.Get("role_name").(string)
	roles := expandStringList(data.Get("roles").([]interface{}))
	users := expandStringList(data.Get("users").([]interface{}))

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
	data.SetId(roleName)
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

type roleGrant struct {
	CreatedOn   sql.RawBytes   `db:"created_on"`
	Role        sql.NullString `db:"role"`
	GrantedTo   sql.NullString `db:"granted_to"`
	GranteeName sql.NullString `db:"grantee_name"`
	Grantedby   sql.NullString `db:"granted_by"`
}

func ReadRoleGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := data.Id()

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

	err = data.Set("role_name", roleName)
	if err != nil {
		return err
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

func readGrants(db *sql.DB, roleName string) ([]*roleGrant, error) {
	sdb := sqlx.NewDb(db, "snowflake")

	stmt := fmt.Sprintf(`SHOW GRANTS OF ROLE "%s"`, roleName)
	rows, err := sdb.Queryx(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grants := make([]*roleGrant, 0)
	for rows.Next() {
		g := &roleGrant{}
		err = rows.StructScan(g)
		if err != nil {
			return nil, err
		}
		grants = append(grants, g)

	}

	for _, g := range grants {
		if g.GranteeName.Valid {
			s := g.GranteeName.String
			s = strings.TrimPrefix(s, `"`)
			s = strings.TrimSuffix(s, `"`)
			g.GranteeName = sql.NullString{String: s}
		}
	}

	return grants, nil
}

func DeleteRoleGrants(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := data.Get("role_name").(string)

	roles := expandStringList(data.Get("roles").([]interface{}))
	users := expandStringList(data.Get("users").([]interface{}))

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

	log.Printf("[DEBUG] Updating role_grants for %s", roleName)

	x := func(resource string, grant func(db *sql.DB, role string, target string) error, revoke func(db *sql.DB, role string, target string) error) error {
		o, n := data.GetChange(resource)

		if o == nil {
			o = []interface{}{}
		}
		if n == nil {
			n = []interface{}{}
		}

		oldEntities := createStringSet(o.([]interface{}))
		newEntities := createStringSet(n.([]interface{}))

		remove := oldEntities.Difference(newEntities).List()
		add := newEntities.Difference(oldEntities).List()

		log.Printf("[DEBUG] role_grants would remove %v from role %s\n", remove, roleName)
		log.Printf("[DEBUG] role_grants would add %v from role %s\n", add, roleName)

		for _, user := range remove {
			log.Printf("[DEBUG] Removing resource %s grant %s from role %s\n", resource, user, roleName)
			err := revoke(db, roleName, user)
			if err != nil {
				return err
			}
		}
		for _, user := range add {
			log.Printf("[DEBUG] Adding resource %s grant %s from role %s\n", resource, user, roleName)
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
