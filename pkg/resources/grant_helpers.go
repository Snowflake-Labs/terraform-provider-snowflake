package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
)

// TerraformGrantResource augments terraform's *schema.Resource with extra context
type TerraformGrantResource struct {
	Resource   *schema.Resource
	ValidPrivs PrivilegeSet
}

type TerraformGrantResources map[string]*TerraformGrantResource

func (t TerraformGrantResources) GetTfSchemas() map[string]*schema.Resource {
	out := map[string]*schema.Resource{}
	for name, grant := range t {
		out[name] = grant.Resource
	}
	return out
}

const (
	grantIDDelimiter = '|'
)

// currentGrant represents a generic grant of a privilege from a grant (the target) to a
// grantee. This type can be used in conjunction with github.com/jmoiron/sqlx to
// build a nice go representation of a grant
type currentGrant struct {
	CreatedOn   time.Time `db:"created_on"`
	Privilege   string    `db:"privilege"`
	GrantType   string    `db:"granted_on"`
	GrantName   string    `db:"name"`
	GranteeType string    `db:"granted_to"`
	GranteeName string    `db:"grantee_name"`
	GrantOption bool      `db:"grant_option"`
	GrantedBy   string    `db:"granted_by"`
}

// futureGrant represents the columns in the response from `SHOW FUTURE GRANTS
// IN SCHEMA...` and can be used in conjunction with sqlx.
type futureGrant struct {
	CreatedOn   time.Time `db:"created_on"`
	Privilege   string    `db:"privilege"`
	GrantType   string    `db:"grant_on"`
	GrantName   string    `db:"name"`
	GranteeType string    `db:"grant_to"`
	GranteeName string    `db:"grantee_name"`
	GrantOption bool      `db:"grant_option"`
}

// grant is simply the least common denominator of fields in currentGrant and
// futureGrant.
type grant struct {
	CreatedOn   time.Time
	Privilege   string
	GrantType   string
	GrantName   string
	GranteeType string
	GranteeName string
	GrantOption bool
}

// grantID contains identifying elements that allow unique access privileges
type grantID struct {
	ResourceName string
	SchemaName   string
	ObjectName   string
	Privilege    string
	GrantOption  bool
}

// String() takes in a grantID object and returns a pipe-delimited string:
// resourceName|schemaName|ObjectName|Privilege|GrantOption
func (gi *grantID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = grantIDDelimiter
	grantOption := fmt.Sprintf("%v", gi.GrantOption)
	dataIdentifiers := [][]string{{gi.ResourceName, gi.SchemaName, gi.ObjectName, gi.Privilege, grantOption}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strGrantID := strings.TrimSpace(buf.String())
	return strGrantID, nil
}

// grantIDFromString() takes in a pipe-delimited string: resourceName|schemaName|ObjectName|Privilege
// and returns a grantID object
func grantIDFromString(stringID string) (*grantID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = grantIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per grant")
	}
	if len(lines[0]) != 4 && len(lines[0]) != 5 {
		return nil, fmt.Errorf("4 or 5 fields allowed")
	}

	grantOption := false
	if len(lines[0]) == 5 && lines[0][4] == "true" {
		grantOption = true
	}

	grantResult := &grantID{
		ResourceName: lines[0][0],
		SchemaName:   lines[0][1],
		ObjectName:   lines[0][2],
		Privilege:    lines[0][3],
		GrantOption:  grantOption,
	}
	return grantResult, nil
}

func createGenericGrant(d *schema.ResourceData, meta interface{}, builder snowflake.GrantBuilder) error {
	db := meta.(*sql.DB)

	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)

	roles, shares := expandRolesAndShares(d)

	if len(roles)+len(shares) == 0 {
		return fmt.Errorf("no roles or shares specified for this grant")
	}

	for _, role := range roles {
		err := snowflake.Exec(db, builder.Role(role).Grant(priv, grantOption))
		if err != nil {
			return err
		}
	}

	for _, share := range shares {
		err := snowflake.Exec(db, builder.Share(share).Grant(priv, grantOption))
		if err != nil {
			return err
		}
	}

	return nil
}

func readGenericGrant(
	d *schema.ResourceData,
	meta interface{},
	schema map[string]*schema.Schema,
	builder snowflake.GrantBuilder,
	futureObjects bool,
	validPrivileges PrivilegeSet) error {
	db := meta.(*sql.DB)
	var grants []*grant
	var err error
	if futureObjects {
		grants, err = readGenericFutureGrants(db, builder)
	} else {
		grants, err = readGenericCurrentGrants(db, builder)
	}
	if err != nil {
		return err
	}
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)

	// This is the only way how I can test that this function is reading VIEW grants or TABLE grants
	// is checking what kind of builder we have. If it is future grant, then I double check if the
	// privilegeSet has only one member - SELECT - then it is a VIEW, if it has 6 members and contains
	// Truncate then it must be Table
	futureGrantOnViews := false
	futureGrantOnTables := false
	if reflect.TypeOf(builder) == reflect.TypeOf(&snowflake.FutureGrantBuilder{}) {
		if _, ok := validPrivileges[privilegeSelect]; ok && len(validPrivileges) == 1 {
			futureGrantOnViews = true
		}
		if _, ok := validPrivileges[privilegeTruncate]; ok && len(validPrivileges) == 6 {
			futureGrantOnTables = true
		}
	}

	// Map of roles to privileges
	rolePrivileges := map[string]PrivilegeSet{}
	sharePrivileges := map[string]PrivilegeSet{}

	// List of all grants for each schema_database
	for _, grant := range grants {
		switch grant.GranteeType {
		case "ROLE":
			roleName := grant.GranteeName
			// Find set of privileges
			privileges, ok := rolePrivileges[roleName]
			if !ok {
				// If not there, create an empty set
				privileges = PrivilegeSet{}
			}
			// Add privilege to the set but consider valid privileges only
			// for VIEW in ReadViewGrant
			// and for non-VIEW in ReadTableGrant
			if futureGrantOnViews || futureGrantOnTables {
				if (futureGrantOnViews && grant.GrantType == "VIEW") || (futureGrantOnTables && grant.GrantType == "TABLE") {
					privileges.addString(grant.Privilege)
				}
			} else { // Other grants
				privileges.addString(grant.Privilege)
			}
			// Reassign set back
			rolePrivileges[roleName] = privileges
		case "SHARE":
			granteeNameStrippedAccount := StripAccountFromName(grant.GranteeName)
			// Find set of privileges
			privileges, ok := sharePrivileges[granteeNameStrippedAccount]
			if !ok {
				// If not there, create an empty set
				privileges = PrivilegeSet{}
			}
			// Add privilege to the set
			privileges.addString(grant.Privilege)
			// Reassign set back
			sharePrivileges[granteeNameStrippedAccount] = privileges
		default:
			return fmt.Errorf("unknown grantee type %s", grant.GranteeType)
		}
	}

	var roles, shares []string
	// Now see which roles have our privilege
	for roleName, privileges := range rolePrivileges {
		// Where priv is not all so it should match exactly
		if privileges.hasString(priv) {
			roles = append(roles, roleName)
		}
	}

	// Now see which shares have our privilege
	for shareName, privileges := range sharePrivileges {
		// Where priv is not all so it should match exactly
		if privileges.hasString(priv) {
			shares = append(shares, shareName)
		}
	}

	err = d.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = d.Set("roles", roles)
	if err != nil {
		return err
	}

	_, sharesOk := schema["shares"]
	if sharesOk && !futureObjects {
		err = d.Set("shares", shares)
		if err != nil {
			return err
		}
	}
	err = d.Set("with_grant_option", grantOption)
	if err != nil {
		return err
	}
	return nil
}

func readGenericCurrentGrants(db *sql.DB, builder snowflake.GrantBuilder) ([]*grant, error) {
	stmt := builder.Show()
	rows, err := snowflake.Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []*grant
	for rows.Next() {
		currentGrant := &currentGrant{}
		err := rows.StructScan(currentGrant)
		if err != nil {
			return nil, err
		}
		if currentGrant.GrantedBy == "" {
			// If GrantedBy is empty string, terraform can't
			// manage the grant because the grant is a default
			// grant seeded by Snowflake.
			continue
		}

		grant := &grant{
			CreatedOn:   currentGrant.CreatedOn,
			Privilege:   currentGrant.Privilege,
			GrantType:   currentGrant.GrantType,
			GrantName:   currentGrant.GrantName,
			GranteeType: currentGrant.GranteeType,
			GranteeName: currentGrant.GranteeName,
			GrantOption: currentGrant.GrantOption,
		}
		grants = append(grants, grant)
	}

	return grants, nil
}

func readGenericFutureGrants(db *sql.DB, builder snowflake.GrantBuilder) ([]*grant, error) {
	conn := sqlx.NewDb(db, "snowflake")

	stmt := builder.Show()
	rows, err := conn.Queryx(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []*grant
	for rows.Next() {
		futureGrant := &futureGrant{}
		err := rows.StructScan(futureGrant)
		if err != nil {
			return nil, err
		}
		grant := &grant{
			CreatedOn:   futureGrant.CreatedOn,
			Privilege:   futureGrant.Privilege,
			GrantType:   futureGrant.GrantType,
			GrantName:   futureGrant.GrantName,
			GranteeType: futureGrant.GranteeType,
			GranteeName: futureGrant.GranteeName,
			GrantOption: futureGrant.GrantOption,
		}
		grants = append(grants, grant)
	}

	return grants, nil
}

func deleteGenericGrant(d *schema.ResourceData, meta interface{}, builder snowflake.GrantBuilder) error {
	db := meta.(*sql.DB)

	priv := d.Get("privilege").(string)

	var roles, shares []string
	if _, ok := d.GetOk("roles"); ok {
		roles = expandStringList(d.Get("roles").(*schema.Set).List())
	}

	if _, ok := d.GetOk("shares"); ok {
		shares = expandStringList(d.Get("shares").(*schema.Set).List())
	}

	for _, role := range roles {
		err := snowflake.Exec(db, builder.Role(role).Revoke(priv))
		if err != nil {
			return err
		}
	}

	for _, share := range shares {
		err := snowflake.Exec(db, builder.Share(share).Revoke(priv))
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}

func expandRolesAndShares(d *schema.ResourceData) ([]string, []string) {
	var roles, shares []string
	if _, ok := d.GetOk("roles"); ok {
		roles = expandStringList(d.Get("roles").(*schema.Set).List())
	}

	if _, ok := d.GetOk("shares"); ok {
		shares = expandStringList(d.Get("shares").(*schema.Set).List())
	}
	return roles, shares
}
