package resources

import (
	"database/sql"
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
	allObjects bool,
	_ PrivilegeSet,
) error {
	db := meta.(*sql.DB)
	var grants []*grant
	var err error

	priv := d.Get("privilege").(string)

	if priv == "ALL PRIVILEGES" {
		// When running e.g. GRANT ALL PRIVILEGES ON TABLE <table_name> TO ROLE <role_name>..., then Snowflake creates a grant for each individual permission
		// There is no way to attribute existing grants to a GRANT ALL PRIVILEGES grant. Thus they cannot be checked. However they can still be revoked.
		return nil
	}
	switch {
	case futureObjects:
		grants, err = readGenericFutureGrants(db, builder)
	case allObjects:
		// When running e.g. GRANT SELECT ON ALL TABLES IN ..., then Snowflake creates a grant for each individual existing table.
		// There is no way to attribute existing table grants to a GRANT SELECT ON ALL TABLES grant. Thus they cannot be checked (or removed).
		return nil
	default:
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
			// Strip account name from grantee name
			s := grant.GranteeName
			granteeNameStrippedAccount := s[strings.LastIndex(s, ".")+1:]
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

	existingRoles := schema.NewSet(schema.HashString, []interface{}{})
	if v, ok := d.GetOk("roles"); ok && v != nil {
		existingRoles = v.(*schema.Set)
	}
	multipleGrantFeatureFlag := d.Get("enable_multiple_grants").(bool)
	var roles, shares []string
	// Now see which roles have our privilege.
	for roleName, privileges := range rolePrivileges {
		if privileges.hasString(priv) {
			// CASE A: Whatever role we were already managing, continue to do so.
			caseA := existingRoles.Contains(roleName)
			// CASE B : If multiple grants is not enabled (meaning this is an authoritative resource) then we care about what roles have privilege unless on_future is enabled in which case we don't care (because we will get flooded with diffs)
			caseB := !multipleGrantFeatureFlag && !futureObjects
			if caseA || caseB {
				roles = append(roles, roleName)
			}
		}
	}

	existingShares := schema.NewSet(schema.HashString, []interface{}{})
	if v, ok := d.GetOk("shares"); ok && v != nil {
		existingShares = v.(*schema.Set)
	}
	// Now see which shares have our privilege.
	for shareName, privileges := range sharePrivileges {
		if privileges.hasString(priv) {
			// CASE A: Whatever share we were already managing, continue to do so.
			caseA := existingShares.Contains(shareName)
			// CASE B : If multiple grants is not enabled (meaning this is an authoritative resource) then we care about what shares have privilege unless on_future is enabled in which case we don't care (because we will get flooded with diffs)
			caseB := !multipleGrantFeatureFlag && !futureObjects
			if caseA || caseB {
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
	reversionRole string,
	roles []string,
	shares []string,
) error {
	db := meta.(*sql.DB)

	for _, role := range roles {
		executable := builder.Role(role).Revoke(priv)
		if priv == "OWNERSHIP" {
			executable = builder.Role(role).RevokeOwnership(reversionRole)
		}
		if err := snowflake.ExecMulti(db, executable); err != nil {
			return err
		}
	}

	for _, share := range shares {
		executable := builder.Share(share).Revoke(priv)
		if priv == "OWNERSHIP" {
			executable = builder.Share(share).RevokeOwnership(reversionRole)
		}
		if err := snowflake.ExecMulti(db, executable); err != nil {
			return err
		}
	}
	return nil
}

func deleteGenericGrant(d *schema.ResourceData, meta interface{}, builder snowflake.GrantBuilder) error {
	priv := d.Get("privilege").(string)
	rr := d.Get("revert_ownership_to_role_name")
	var reversionRole string
	if rr != nil {
		reversionRole = rr.(string)
	}
	roles, shares := expandRolesAndShares(d)
	if err := deleteGenericGrantRolesAndShares(meta, builder, priv, reversionRole, roles, shares); err != nil {
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

// changeDiff calculates roles/shares to add/revoke.
func changeDiff(d *schema.ResourceData, key string) (toAdd []string, toRemove []string) {
	o, n := d.GetChange(key)
	oldSet := o.(*schema.Set)
	newSet := n.(*schema.Set)
	toAdd = expandStringList(newSet.Difference(oldSet).List())
	toRemove = expandStringList(oldSet.Difference(newSet).List())
	return
}
