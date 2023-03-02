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

var validSchemaPrivileges = NewPrivilegeSet(
	privilegeAddSearchOptimization,
	privilegeCreateExternalTable,
	privilegeCreateFileFormat,
	privilegeCreateFunction,
	privilegeCreateMaskingPolicy,
	privilegeCreateMaterializedView,
	privilegeCreatePipe,
	privilegeCreateProcedure,
	privilegeCreateRowAccessPolicy,
	privilegeCreateSequence,
	privilegeCreateSessionPolicy,
	privilegeCreateStage,
	privilegeCreateStream,
	privilegeCreateTable,
	privilegeCreateTag,
	privilegeCreateTask,
	privilegeCreateTemporaryTable,
	privilegeCreateView,
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeUsage,
)

var schemaGrantSchema = map[string]*schema.Schema{
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future schema. Note that if \"OWNERSHIP\" is specified, ensure that the role that terraform is using is granted access.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validSchemaPrivileges.ToList(), true),
		ForceNew:     true,
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
		Description: "Grants privilege to these shares (only valid if on_future is unset).",
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true, apply this grant on all future schemas in the given database. The schema_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"schema_name", "shares"},
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

// SchemaGrant returns a pointer to the resource representing a view grant.
func SchemaGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateSchemaGrant,
			Read:   ReadSchemaGrant,
			Delete: DeleteSchemaGrant,
			Update: UpdateSchemaGrant,

			Schema: schemaGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validSchemaPrivileges,
	}
}

// CreateSchemaGrant implements schema.CreateFunc.
func CreateSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	var schemaName string
	if _, ok := d.GetOk("schema_name"); ok {
		schemaName = d.Get("schema_name").(string)
	}
	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(databaseName)
	} else {
		builder = snowflake.SchemaGrant(databaseName, schemaName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewSchemaGrantID(databaseName, schemaName, privilege, roles, shares, withGrantOption)
	d.SetId(grantID.String())

	return ReadSchemaGrant(d, meta)
}

// UpdateSchemaGrant implements schema.UpdateFunc.
func UpdateSchemaGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := parseSchemaGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := d.Get("on_future").(bool)

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(grantID.DatabaseName)
	} else {
		builder = snowflake.SchemaGrant(grantID.DatabaseName, grantID.SchemaName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		rolesToRevoke,
		sharesToRevoke,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		grantID.WithGrantOption,
		rolesToAdd,
		sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadSchemaGrant(d, meta)
}

// ReadSchemaGrant implements schema.ReadFunc.
func ReadSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseSchemaGrantID(d.Id())
	if err != nil {
		return err
	}

	if !grantID.IsOldID {
		if err := d.Set("roles", grantID.Roles); err != nil {
			return err
		}
		if err := d.Set("shares", grantID.Shares); err != nil {
			return err
		}
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	onFuture := false
	if grantID.SchemaName == "" {
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
		builder = snowflake.FutureSchemaGrant(grantID.DatabaseName)
	} else {
		builder = snowflake.SchemaGrant(grantID.DatabaseName, grantID.SchemaName)
	}
	return readGenericGrant(d, meta, schemaGrantSchema, builder, onFuture, validSchemaPrivileges)
}

// DeleteSchemaGrant implements schema.DeleteFunc.
func DeleteSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseSchemaGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := false
	if grantID.SchemaName == "" {
		onFuture = true
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(grantID.DatabaseName)
	} else {
		builder = snowflake.SchemaGrant(grantID.DatabaseName, grantID.SchemaName)
	}
	return deleteGenericGrant(d, meta, builder)
}

type SchemaGrantID struct {
	DatabaseName    string
	SchemaName      string
	Privilege       string
	Roles           []string
	Shares          []string
	WithGrantOption bool
	IsOldID         bool
}

func NewSchemaGrantID(databaseName string, schemaName, privilege string, roles []string, shares []string, withGrantOption bool) *SchemaGrantID {
	return &SchemaGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		Privilege:       privilege,
		Roles:           roles,
		Shares:          shares,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *SchemaGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.Privilege, v.WithGrantOption, roles, shares)
}

func parseSchemaGrantID(s string) (*SchemaGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &SchemaGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			Shares:          []string{},
			WithGrantOption: idParts[5] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) != 6 {
		idParts := strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 6", len(idParts))
	}
	return &SchemaGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		Privilege:       idParts[2],
		WithGrantOption: idParts[3] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[4], ","),
		Shares:          helpers.SplitStringToSlice(idParts[5], ","),
		IsOldID:         false,
	}, nil
}
