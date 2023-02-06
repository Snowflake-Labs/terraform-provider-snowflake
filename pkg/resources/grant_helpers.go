package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/snowflakedb/gosnowflake"
)

// TerraformGrantResource augments terraform's *schema.Resource with extra context.
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
// build a nice go representation of a grant.
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

// grantID contains identifying elements that allow unique access privileges.
type grantID struct {
	ResourceName string
	SchemaName   string
	ObjectName   string
	Privilege    string
	Roles        []string
	GrantOption  bool
}

// String() takes in a grantID object and returns a pipe-delimited string:
// resourceName|schemaName|ObjectName|Privilege|Roles|GrantOption.
func (gi *grantID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = grantIDDelimiter
	grantOption := fmt.Sprintf("%v", gi.GrantOption)
	roles := strings.Join(gi.Roles, ",")
	dataIdentifiers := [][]string{{gi.ResourceName, gi.SchemaName, gi.ObjectName, gi.Privilege, roles, grantOption}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strGrantID := strings.TrimSpace(buf.String())
	return strGrantID, nil
}

// grantIDFromString() takes in a pipe-delimited string: resourceName|schemaName|ObjectName|Privilege|Roles
// and returns a grantID object.
func grantIDFromString(stringID string) (*grantID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = grantIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per grant")
	}

	// Len 1 is allowing for legacy IDs where role names are not included
	if len(lines[0]) < 1 || len(lines[0]) > 6 {
		return nil, fmt.Errorf("1 to 6 fields allowed in ID")
	}

	// Splitting string list if new ID structure, will cause issues if roles names passed are "true" or "false".
	// Checking for true/false to eliminate scenarios where it would pick up the grant option.
	// Roles will be empty list if legacy IDs are used, roles from grants are not
	// used in Read functions, just for uniqueness in IDs of resources
	roles := []string{}
	if len(lines[0]) > 4 && lines[0][4] != "true" && lines[0][4] != "false" {
		roles = strings.Split(lines[0][4], ",")
	}

	// Allowing legacy IDs to check grant option
	grantOption := false
	if len(lines[0]) == 6 && lines[0][5] == "true" {
		grantOption = true
	} else if len(lines[0]) == 5 && lines[0][4] == "true" {
		grantOption = true
	}

	schemaName := ""
	objectName := ""
	privilege := ""

	if len(lines[0]) > 3 {
		schemaName = lines[0][1]
		objectName = lines[0][2]
		privilege = lines[0][3]
	}

	grantResult := &grantID{
		ResourceName: lines[0][0],
		SchemaName:   schemaName,
		ObjectName:   objectName,
		Privilege:    privilege,
		Roles:        roles,
		GrantOption:  grantOption,
	}
	return grantResult, nil
}

// createGenericGrantRolesAndShares will create generic grants for a set of roles and shares.
func createGenericGrantRolesAndShares(
	meta interface{},
	builder snowflake.GrantBuilder,
	priv string,
	grantOption bool,
	roles []string,
	shares []string,
) error {
	db := meta.(*sql.DB)
	for _, role := range roles {
		if err := snowflake.Exec(db, builder.Role(role).Grant(priv, grantOption)); err != nil {
			return err
		}
	}

	for _, share := range shares {
		if err := snowflake.Exec(db, builder.Share(share).Grant(priv, grantOption)); err != nil {
			return err
		}
	}
	return nil
}

func createGenericGrant(d *schema.ResourceData, meta interface{}, builder snowflake.GrantBuilder) error {
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	roles, shares := expandRolesAndShares(d)

	return createGenericGrantRolesAndShares(
		meta,
		builder,
		priv,
		grantOption,
		roles,
		shares,
	)
}

func readGenericGrant(
	d *schema.ResourceData,
	meta interface{},
	grantSchema map[string]*schema.Schema,
	builder snowflake.GrantBuilder,
	futureObjects bool,
	validPrivileges PrivilegeSet,
) error {
	db := meta.(*sql.DB)
	var grants []*grant
	var err error
	if futureObjects {
		grants, err = readGenericFutureGrants(db, builder)
	} else {
		grants, err = readGenericCurrentGrants(db, builder)
	}
	if err != nil {
		// HACK HACK: If the object doesn't exist or not authorized then we can assume someone deleted it
		// We also check the error number matches
		// We set the tf id == blank and return.
		// I don't know of a better way to work around this issue
		if snowflakeErr, ok := err.(*gosnowflake.SnowflakeError); ok && //nolint:errorlint // todo: should be fixed
			snowflakeErr.Number == 2003 &&
			strings.Contains(err.Error(), "does not exist or not authorized") {
			log.Printf("[WARN] resource (%s) not found, removing from state file", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)

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

			if strings.ReplaceAll(builder.GrantType(), " ", "_") == grant.GrantType {
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

	var existingRoles *schema.Set
	if v, ok := d.GetOk("roles"); ok && v != nil {
		existingRoles = v.(*schema.Set)
	}
	multipleGrantFeatureFlag := d.Get("enable_multiple_grants").(bool)
	var roles, shares []string
	// Now see which roles have our privilege.
	for roleName, privileges := range rolePrivileges {
		if privileges.hasString(priv) {
			// CASE A: If multiple grants is not enabled (meaning this is authoritative) then we always care about what roles have privilige.
			caseA := !multipleGrantFeatureFlag
			// CASE B: If this is not authoritative, then at least continue managing whatever roles we already are managing
			caseB := multipleGrantFeatureFlag && existingRoles.Contains(roleName)
			// CASE C: If this is not authoritative and we are not managing the role, then we only care about the role if future objects is disabled. Otherwise we will get flooded with diffs.
			caseC := multipleGrantFeatureFlag && !futureObjects
			if caseA || caseB || caseC {
				roles = append(roles, roleName)
			}
		}
	}

	var existingShares *schema.Set
	if v, ok := d.GetOk("shares"); ok && v != nil {
		existingShares = v.(*schema.Set)
	}
	// Now see which shares have our privilege.
	for shareName, privileges := range sharePrivileges {
		if privileges.hasString(priv) {
			// CASE A: If multiple grants is not enabled (meaning this is authoritative) then we always care about what shares have privilige.
			caseA := !multipleGrantFeatureFlag
			// CASE B: If this is not authoritative, then at least continue managing whatever shares we already are managing
			caseB := multipleGrantFeatureFlag && existingShares.Contains(shareName)
			// CASE C: If this is not authoritative and we are not managing the share, then we only care about the share if future objects is disabled. Otherwise we will get flooded with diffs.
			caseC := multipleGrantFeatureFlag && !futureObjects
			if caseA || caseB || caseC {
				shares = append(shares, shareName)
			}
		}
	}

	if err := d.Set("privilege", priv); err != nil {
		return err
	}
	if err := d.Set("roles", roles); err != nil {
		return err
	}

	_, sharesOk := grantSchema["shares"]
	if sharesOk && !futureObjects {
		if err := d.Set("shares", shares); err != nil {
			return err
		}
	}
	if err := d.Set("with_grant_option", grantOption); err != nil {
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
		if err := rows.StructScan(currentGrant); err != nil {
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
		if err := rows.StructScan(futureGrant); err != nil {
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

// Deletes specific roles and shares from a grant
// Does not modify TF remote state.
func deleteGenericGrantRolesAndShares(
	meta interface{},
	builder snowflake.GrantBuilder,
	priv string,
	roles []string,
	shares []string,
) error {
	db := meta.(*sql.DB)

	for _, role := range roles {
		if err := snowflake.ExecMulti(db, builder.Role(role).Revoke(priv)); err != nil {
			return err
		}
	}

	for _, share := range shares {
		if err := snowflake.ExecMulti(db, builder.Share(share).Revoke(priv)); err != nil {
			return err
		}
	}
	return nil
}

func deleteGenericGrant(d *schema.ResourceData, meta interface{}, builder snowflake.GrantBuilder) error {
	priv := d.Get("privilege").(string)
	roles, shares := expandRolesAndShares(d)
	if err := deleteGenericGrantRolesAndShares(meta, builder, priv, roles, shares); err != nil {
		return err
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

// parseFunctionObjectName parses a callable object name (including procedures) into its identifier components. For example, functions and procedures.
func parseFunctionObjectName(objectIdentifier string) (string, []string) {
	nameIndex := strings.Index(objectIdentifier, `(`)
	if nameIndex == -1 {
		return "", []string{}
	}
	name := objectIdentifier[:nameIndex]
	argumentString := objectIdentifier[nameIndex+1:]

	// Backwards compatibility for functions with return_types (prior to 0.56.1).
	if strings.Contains(argumentString, ":") {
		argumentString = strings.Split(argumentString, ":")[0]
	}

	// Remove trailing ")".
	argumentString = strings.TrimRight(argumentString, `)`)
	arguments := strings.Split(argumentString, `,`)
	for i, argument := range arguments {
		arguments[i] = strings.TrimSpace(argument)
	}
	return name, arguments
}

// changeDiff calculates roles/shares to add/revoke.
func changeDiff(d *schema.ResourceData, key string) (toAdd []string, toRemove []string) {
	o, n := d.GetChange(key)
	oldSet := o.(*schema.Set)
	newSet := n.(*schema.Set)
	toAdd = expandStringList(newSet.Difference(oldSet).List())
	toRemove = expandStringList(oldSet.Difference(newSet).List())
	return
}
