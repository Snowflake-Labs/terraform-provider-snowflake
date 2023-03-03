package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validDatabasePrivileges = NewPrivilegeSet(
	privilegeCreateSchema,
	privilegeImportedPrivileges,
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeReferenceUsage,
	privilegeUsage,
)

var databaseGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the database.",
		Default:      "USAGE",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(validDatabasePrivileges.ToList(), true),
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"shares": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares.",
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
	},
}

// DatabaseGrant returns a pointer to the resource representing a database grant.
func DatabaseGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateDatabaseGrant,
			Read:   ReadDatabaseGrant,
			Delete: DeleteDatabaseGrant,
			Update: UpdateDatabaseGrant,

			Schema: databaseGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validDatabasePrivileges,
	}
}

// CreateDatabaseGrant implements schema.CreateFunc.
func CreateDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(databaseName)
	if err := createGenericGrant(d, meta, builder); err != nil {
		return fmt.Errorf("error creating database grant err = %w", err)
	}

	privilege := d.Get("privilege").(string)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())
	withGrantOption := d.Get("with_grant_option").(bool)
	grantID := NewDatabaseGrantID(databaseName, privilege, roles, shares, withGrantOption)

	d.SetId(grantID.String())

	return ReadDatabaseGrant(d, meta)
}

// ReadDatabaseGrant implements schema.ReadFunc.
func ReadDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseDatabaseGrantID(d.Id())
	if err != nil {
		return err
	}
	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("roles", grantID.Roles); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}
	if !grantID.IsOldID {
		if err := d.Set("shares", grantID.Shares); err != nil {
			return err
		}
	}

	builder := snowflake.DatabaseGrant(grantID.DatabaseName)
	return readGenericGrant(d, meta, databaseGrantSchema, builder, false, validDatabasePrivileges)
}

// DeleteDatabaseGrant implements schema.DeleteFunc.
func DeleteDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(databaseName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateDatabaseGrant implements schema.UpdateFunc.
func UpdateDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update and we're done
	if !d.HasChanges("roles", "shares") {
		return nil
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}
	sharesToAdd := []string{}
	sharesToRevoke := []string{}
	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}
	if d.HasChange("shares") {
		sharesToAdd, sharesToRevoke = changeDiff(d, "shares")
	}

	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	// create the builder
	builder := snowflake.DatabaseGrant(databaseName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		rolesToRevoke,
		sharesToRevoke,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		withGrantOption,
		rolesToAdd,
		sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadDatabaseGrant(d, meta)
}

type DatabaseGrantID struct {
	DatabaseName    string
	Privilege       string
	Roles           []string
	Shares          []string
	WithGrantOption bool
	IsOldID         bool
}

func NewDatabaseGrantID(databaseName string, privilege string, roles []string, shares []string, withGrantOption bool) *DatabaseGrantID {
	return &DatabaseGrantID{
		DatabaseName:    databaseName,
		Privilege:       privilege,
		Roles:           roles,
		Shares:          shares,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *DatabaseGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v", v.DatabaseName, v.Privilege, v.WithGrantOption, roles, shares)
}

func ParseDatabaseGrantID(s string) (*DatabaseGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &DatabaseGrantID{
			DatabaseName:    idParts[0],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			Shares:          []string{},
			WithGrantOption: idParts[5] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 5 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 5 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 5", len(idParts))
	}
	return &DatabaseGrantID{
		DatabaseName:    idParts[0],
		Privilege:       idParts[1],
		WithGrantOption: idParts[2] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[3], ","),
		Shares:          helpers.SplitStringToSlice(idParts[4], ","),
		IsOldID:         false,
	}, nil
}
