package resources

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/jmoiron/sqlx"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

// grant represents a generic grant of a privilge from a grant (the target) to a
// grantee. This type can be used in conjunction with github.com/jmoiron/sqlx to
// build a nice go replresentation of a grant
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

var validDatabasePrivileges = []string{"USAGE", "REFERENCE_USAGE"}

var databaseGrantSchema = map[string]*schema.Schema{
	"database_name": &schema.Schema{
		Type: schema.TypeString,
		Required: true,
		Description: "The name of the database on which to grant privileges.",
		ForceNew: true,
	},
	"privilege": &schema.Schema{
		Type: schema.TypeString,
		Optional: true,
		Description: "The privilege to grant to the database.",
		Default: "USAGE",
		ValidateFunc: validation.StringInSlice(validDatabasePrivileges, true),
		ForceNew: true,
	},
	"roles": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew: true,
	},
	"shares": &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{Type: schema.TypeString},
		Optional: true,
		Description: "Grants privilege to these shares.",
		ForceNew: true,
	},
}

// DatabaseGrant returns a pointer to the resource representing a database grant
func DatabaseGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabaseGrant,
		Read: ReadDatabaseGrant,
		Delete: DeleteDatabaseGrant,

		Schema: databaseGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateDatabaseGrant implements schema.CreateFunc
func CreateDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)

	roles := expandStringList(data.Get("roles").(*schema.Set).List())
	shares := expandStringList(data.Get("shares").(*schema.Set).List())

	if len(roles) + len(shares) == 0 {
		return fmt.Errorf("no roles or shares specified for database grants")
	}

	builder := snowflake.DatabaseGrant(dbName)
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

	// ID format is <database_name>_<privilege>
	data.SetId(fmt.Sprintf("%v_%v", dbName, priv))
	return ReadDatabaseGrant(data, meta)
}

// ReadDatabaseGrant implements schema.ReadFunc
func ReadDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, priv, err := splitID(data.Id())
	if err != nil {
		return err
	}

	grants, err := readDatabaseGrants(db, dbName)

	var roles []string
	var shares []string

	for _, grant := range grants {
		if grant.Privilege != priv {
			continue
		}

		switch grant.GranteeType {
		case "ROLE":
			roles = append(roles, grant.GranteeName)
		case "SHARE":
			// Shares get the account appended to their name, remove this
			shares = append(shares, StripAccountFromName(grant.GranteeName))
		default:
			return fmt.Errorf("unknown grantee type %s", grant.GranteeType)
		}
	}

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// splitID takes the <database_name>_<privilege> ID and returns the database
// name and privilege. Privileges never have underscores so we'll just split on 
// the final underscore
func splitID(v string) (string, string, error) {
	i := strings.LastIndex(v, "_")
	if i == -1 {
		return "", "", fmt.Errorf("ID %v is invalid", v)
	}

	return v[:i], v[i+1:], nil
}

func readDatabaseGrants(db *sql.DB, n string) ([]*grant, error) {
	conn := sqlx.NewDb(db, "snowflake")

	stmt := snowflake.DatabaseGrant(n).Show()
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

// DeleteDatabaseGrant implements schema.DeleteFunc
func DeleteDatabaseGrant(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	dbName, priv, err := splitID(data.Id())
	if err != nil {
		return err
	}

	roles := expandStringList(data.Get("roles").(*schema.Set).List())
	shares := expandStringList(data.Get("shares").(*schema.Set).List())

	builder := snowflake.DatabaseGrant(dbName)
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
