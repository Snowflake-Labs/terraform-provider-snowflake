package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validExternalTablePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeReferences,
	privilegeSelect,
)

var externalTableGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future external tables on which to grant privileges.",
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"external_table_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the external table on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future external tables in the given schema. When this is true and no schema_name is provided apply this grant on all future external tables in the given database. The external_table_name and shares fields must be unset in order to use on_future.",
		Default:     false,
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future external table.",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(validExternalTablePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Grants privilege to these roles.",
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future external tables on which to grant privileges.",
		ForceNew:    true,
	},
	"shares": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares (only valid if on_future is false).",
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// ExternalTableGrant returns a pointer to the resource representing a external table grant.
func ExternalTableGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateExternalTableGrant,
			Read:   ReadExternalTableGrant,
			Delete: DeleteExternalTableGrant,
			Update: UpdateExternalTableGrant,

			Schema: externalTableGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validExternalTablePrivileges,
	}
}

// CreateExternalTableGrant implements schema.CreateFunc.
func CreateExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	var externalTableName string
	if name, ok := d.GetOk("external_table_name"); ok {
		externalTableName = name.(string)
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	if (externalTableName == "") && !onFuture {
		return errors.New("external_table_name must be set unless on_future is true")
	}
	if (externalTableName != "") && onFuture {
		return errors.New("external_table_name must be empty if on_future is true")
	}
	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureExternalTableGrant(databaseName, schemaName)
	} else {
		builder = snowflake.ExternalTableGrant(databaseName, schemaName, externalTableName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewExternalTableGrantID(databaseName, schemaName, externalTableName, privilege, roles, shares, withGrantOption)
	d.SetId(grantID.String())

	return ReadExternalTableGrant(d, meta)
}

// ReadExternalTableGrant implements schema.ReadFunc.
func ReadExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseExternalTableGrantID(d.Id())
	if err != nil {
		return err
	}

	if !grantID.IsOldID {
		fmt.Printf("[DEBUG] id: %v\n", d.Id())
		fmt.Printf("[DEBUG] reading external table grant: %v\n", grantID)
		fmt.Printf("[DEBUG] reading external table grant shares: %v\n", grantID.Shares)
		fmt.Printf("[DEBUG] len(external table grant shares): %v\n", len(grantID.Shares))
		if err := d.Set("shares", grantID.Shares); err != nil {
			return err
		}
	}

	if err := d.Set("roles", grantID.Roles); err != nil {
		return err
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	if err := d.Set("external_table_name", grantID.ObjectName); err != nil {
		return err
	}

	onFuture := false
	if grantID.ObjectName == "" {
		onFuture = true
	}
	if err := d.Set("on_future", onFuture); err != nil {
		return err
	}

	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureExternalTableGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ExternalTableGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, externalTableGrantSchema, builder, onFuture, validExternalTablePrivileges)
}

// DeleteExternalTableGrant implements schema.DeleteFunc.
func DeleteExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseExternalTableGrantID(d.Id())
	if err != nil {
		return err
	}
	onFuture := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureExternalTableGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ExternalTableGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateExternalTableGrant implements schema.UpdateFunc.
func UpdateExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := ParseExternalTableGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureExternalTableGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ExternalTableGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, sharesToRevoke,
	); err != nil {
		return err
	}
	// then add

	if err := createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.WithGrantOption, rolesToAdd, sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadExternalTableGrant(d, meta)
}

type ExternalTableGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	Shares          []string
	WithGrantOption bool
	IsOldID         bool
}

func NewExternalTableGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, shares []string, withGrantOption bool) *ExternalTableGrantID {
	return &ExternalTableGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		Shares:          shares,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *ExternalTableGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles, shares)
}

func ParseExternalTableGrantID(s string) (*ExternalTableGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &ExternalTableGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			Shares:          []string{},
			WithGrantOption: idParts[5] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 7 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 7 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 7", len(idParts))
	}
	return &ExternalTableGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		WithGrantOption: idParts[4] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
		Shares:          helpers.SplitStringToSlice(idParts[6], ","),
		IsOldID:         false,
	}, nil
}
