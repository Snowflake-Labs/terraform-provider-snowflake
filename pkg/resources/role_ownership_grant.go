package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmoiron/sqlx"
)

func RoleOwnershipGrants() *schema.Resource {
	return &schema.Resource{
		Create: CreateRoleOwnershipGrants,
		Read:   ReadRoleOwnershipGrants,
		Delete: DeleteRoleOwnershipGrants,
		Update: UpdateRoleOwnershipGrants,

		Schema: map[string]*schema.Schema{
			"on_role_name": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The name of the role ownership is granted on.",
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					return snowflake.ValidateIdentifier(val)
				},
			},
			"to_role_name": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The name of the role to grant ownership. Please ensure that the role that terraform is using is granted access.",
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					return snowflake.ValidateIdentifier(val)
				},
			},
			"current_grants": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role.",
				Default:     "COPY",
				ValidateFunc: validation.StringInSlice([]string{
					"COPY",
					"REVOKE",
				}, true),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	onRoleName := d.Get("on_role_name").(string)
	toRoleName := d.Get("to_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	g := snowflake.RoleOwnershipGrant(onRoleName)
	err := snowflake.Exec(db, g.Role(toRoleName, currentGrants).Grant())
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf(`%s|%s|%s`, onRoleName, toRoleName, currentGrants))

	return ReadRoleOwnershipGrants(d, meta)
}

func ReadRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	log.Println(d.Id())
	onRoleName := strings.Split(d.Id(), "|")[0]
	toRoleName := strings.Split(d.Id(), "|")[1]
	currentGrants := strings.Split(d.Id(), "|")[2]

	err := readOwnershipGrants(db, onRoleName)
	if err != nil {
		return err
	}

	err = d.Set("on_role_name", onRoleName)
	if err != nil {
		return err
	}

	err = d.Set("to_role_name", toRoleName)
	if err != nil {
		return err
	}

	err = d.Set("current_grants", currentGrants)
	if err != nil {
		return err
	}

	return nil
}

type roleOwnershipGrant struct {
	CreatedOn       sql.RawBytes   `db:"created_on"`
	Name            sql.NullString `db:"name"`
	IsDefault       sql.NullString `db:"is_default"`
	IsCurrent       sql.NullString `db:"is_current"`
	IsInherited     sql.NullString `db:"is_inherited"`
	AssignedToUsers sql.NullString `db:"assigned_to_users"`
	GrantedToRoles  sql.NullString `db:"granted_to_roles"`
	GrantedRoles    sql.NullString `db:"granted_roles"`
	Owner           sql.NullString `db:"owner"`
	Comment         sql.NullString `db:"comment"`
}

func readOwnershipGrants(db *sql.DB, onRoleName string) error {
	sdb := sqlx.NewDb(db, "snowflake")

	stmt := fmt.Sprintf(`SHOW ROLES LIKE '%s'`, onRoleName)

	row := sdb.QueryRowx(stmt)

	g := &roleOwnershipGrant{}
	err := row.StructScan(g)
	if err != nil {
		return err
	}

	if g.Owner.Valid {
		s := g.Owner.String
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		g.Owner = sql.NullString{String: s}
	}

	return nil
}

func DeleteRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	onRoleName := d.Get("on_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	g := snowflake.RoleOwnershipGrant(onRoleName, currentGrants)
	err := snowflake.Exec(db, g.Role("ACCOUNTADMIN").Revoke())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func UpdateRoleOwnershipGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	onRoleName := d.Get("on_role_name").(string)
	toRoleName := d.Get("to_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	g := snowflake.RoleOwnershipGrant(onRoleName, currentGrants)
	err := snowflake.Exec(db, g.Role(toRoleName).Revoke())
	if err != nil {
		return err
	}

	return ReadRoleOwnershipGrants(d, meta)
}
