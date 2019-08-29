package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jmoiron/sqlx"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

// grant represents a generic grant of a privilege from a grant (the target) to a
// grantee. This type can be used in conjunction with github.com/jmoiron/sqlx to
// build a nice go representation of a grant
type grant struct {
	CreatedOn   time.Time `db:"created_on"`
	Privilege   string    `db:"privilege"`
	GrantType   string    `db:"granted_on"`
	GrantName   string    `db:"name"`
	GranteeType string    `db:"granted_to"`
	GranteeName string    `db:"grantee_name"`
	GrantOption bool      `db:"grant_option"`
	GrantedBy   string    `db:"granted_by"`
}

// splitGrantID takes the <db_name>|<schema_name>|<view_name>|<privilege> ID and
// returns the object name and privilege.
func splitGrantID(v string) (string, string, string, string, error) {
	arr := strings.Split(v, "|")
	if len(arr) != 4 {
		return "", "", "", "", fmt.Errorf("ID %v is invalid", v)
	}

	return arr[0], arr[1], arr[2], arr[3], nil
}

func createGenericGrant(data *schema.ResourceData, meta interface{}, builder *snowflake.GrantBuilder) error {
	db := meta.(*sql.DB)

	priv := data.Get("privilege").(string)

	roles, shares := expandRolesAndShares(data)

	if len(roles)+len(shares) == 0 {
		return fmt.Errorf("no roles or shares specified for this grant")
	}

	for _, role := range roles {
		err := DBExec(db, builder.Role(role).Grant(priv))
		if err != nil {
			return err
		}
	}

	for _, share := range shares {
		err := DBExec(db, builder.Share(share).Grant(priv))
		if err != nil {
			return err
		}
	}

	return nil
}

func d(in interface{}) {
	log.Printf("[DEBUG]%#v\n", in)
}

func readGenericGrant(data *schema.ResourceData, meta interface{}, builder *snowflake.GrantBuilder) error {
	db := meta.(*sql.DB)
	grants, err := readGenericGrants(db, builder)
	if err != nil {
		return err
	}
	priv := data.Get("privilege").(string)

	rolesIn, sharesIn := expandRolesAndShares(data)

	var roles, shares []string
	d("foo")
	for _, grant := range grants {
		// Skip if wrong privilege
		if grant.Privilege != priv {
			continue
		}
		d(grant)
		switch grant.GranteeType {
		case "ROLE":
			if !stringInSlice(grant.GranteeName, rolesIn) {
				continue
			}
			roles = append(roles, grant.GranteeName)
		case "SHARE":
			// Shares get the account appended to their name, remove this
			granteeName := StripAccountFromName(grant.GranteeName)
			if !stringInSlice(granteeName, sharesIn) {
				continue
			}

			shares = append(shares, granteeName)
		default:
			return fmt.Errorf("unknown grantee type %s", grant.GranteeType)
		}
	}

	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	err = data.Set("roles", roles)
	if err != nil {
		return err
	}

	err = data.Set("shares", shares)
	if err != nil {
		// warehouses don't use shares - check for this error
		if !strings.HasPrefix(err.Error(), "Invalid address to set") {
			return err
		}
	}

	return nil
}

func readGenericGrants(db *sql.DB, builder *snowflake.GrantBuilder) ([]*grant, error) {
	conn := sqlx.NewDb(db, "snowflake")

	stmt := builder.Show()
	rows, err := conn.Queryx(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []*grant
	for rows.Next() {
		grant := &grant{}
		err := rows.StructScan(grant)
		if err != nil {
			return nil, err
		}
		grants = append(grants, grant)
	}

	return grants, nil
}

func deleteGenericGrant(data *schema.ResourceData, meta interface{}, builder *snowflake.GrantBuilder) error {
	db := meta.(*sql.DB)

	priv := data.Get("privilege").(string)

	var roles, shares []string
	if _, ok := data.GetOk("roles"); ok {
		roles = expandStringList(data.Get("roles").(*schema.Set).List())
	}

	if _, ok := data.GetOk("shares"); ok {
		shares = expandStringList(data.Get("shares").(*schema.Set).List())
	}

	for _, role := range roles {
		err := DBExec(db, builder.Role(role).Revoke(priv))
		if err != nil {
			return err
		}
	}

	for _, share := range shares {
		err := DBExec(db, builder.Share(share).Revoke(priv))
		if err != nil {
			return err
		}
	}

	data.SetId("")
	return nil
}

func expandRolesAndShares(data *schema.ResourceData) ([]string, []string) {
	var roles, shares []string
	if _, ok := data.GetOk("roles"); ok {
		roles = expandStringList(data.Get("roles").(*schema.Set).List())
	}

	if _, ok := data.GetOk("shares"); ok {
		shares = expandStringList(data.Get("shares").(*schema.Set).List())
	}
	return roles, shares
}

func stringInSlice(v string, sl []string) bool {
	for _, s := range sl {
		if s == v {
			return true
		}
	}
	return false
}
