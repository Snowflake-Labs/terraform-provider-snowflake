package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validStagePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
	// These privileges are only valid for internal stages
	privilegeRead,
	privilegeWrite,
)

var stageGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current stage on which to grant privileges.",
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future stages in the given schema. When this is true and no schema_name is provided apply this grant on all future stages in the given database. The stage_name field must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"stage_name"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all stages in the given schema. When this is true and no schema_name is provided apply this grant on all stages in the given database. The stage_name field must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"stage_name"},
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the stage.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validStagePrivileges.ToList(), true),
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
		Description: "The name of the schema containing the current stage on which to grant privileges.",
		ForceNew:    true,
	},
	"stage_name": {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The name of the stage on which to grant privilege (only valid if on_future is false).",
		ForceNew:      true,
		ConflictsWith: []string{"on_future"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// StageGrant returns a pointer to the resource representing a stage grant.
func StageGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateStageGrant,
			Read:   ReadStageGrant,
			Delete: DeleteStageGrant,
			Update: UpdateStageGrant,

			Schema: stageGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					grantID, err := ParseStageGrantID(d.Id())
					if err != nil {
						return nil, err
					}
					if err := d.Set("stage_name", grantID.ObjectName); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", grantID.SchemaName); err != nil {
						return nil, err
					}
					if err := d.Set("database_name", grantID.DatabaseName); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", grantID.Privilege); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
						return nil, err
					}
					if err := d.Set("roles", grantID.Roles); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validStagePrivileges,
	}
}

// CreateStageGrant implements schema.CreateFunc.
func CreateStageGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	if onFuture && onAll {
		return errors.New("on_future and on_all cannot both be true")
	}
	privilege := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)

	var schemaName string
	if name, ok := d.GetOk("schema_name"); ok {
		schemaName = name.(string)
	}
	var stageName string
	if name, ok := d.GetOk("stage_name"); ok {
		stageName = name.(string)
	}

	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}
	if (stageName == "") && !onFuture && !onAll {
		return errors.New("stage_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllStageGrant(databaseName, schemaName)
	default:
		builder = snowflake.StageGrant(databaseName, schemaName, stageName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	grantID := NewStageGrantID(databaseName, schemaName, stageName, privilege, roles, grantOption)
	d.SetId(grantID.String())

	return ReadStageGrant(d, meta)
}

// ReadStageGrant implements schema.ReadFunc.
func ReadStageGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseStageGrantID(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	if err := d.Set("stage_name", grantID.ObjectName); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(grantID.DatabaseName, grantID.SchemaName)
	case onAll:
		builder = snowflake.AllStageGrant(grantID.DatabaseName, grantID.SchemaName)
	default:
		builder = snowflake.StageGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, stageGrantSchema, builder, onFuture, onAll, validStagePrivileges)
}

// UpdateStageGrant implements schema.UpdateFunc.
func UpdateStageGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update, and we're done
	if !d.HasChanges("roles") {
		return nil
	}

	grantID, err := ParseStageGrantID(d.Id())
	if err != nil {
		return err
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}

	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}

	var builder snowflake.GrantBuilder
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(grantID.DatabaseName, grantID.SchemaName)
	case onAll:
		builder = snowflake.AllStageGrant(grantID.DatabaseName, grantID.SchemaName)
	default:
		builder = snowflake.StageGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.WithGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadStageGrant(d, meta)
}

// DeleteStageGrant implements schema.DeleteFunc.
func DeleteStageGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseStageGrantID(d.Id())
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(grantID.DatabaseName, grantID.SchemaName)
	case onAll:
		builder = snowflake.AllStageGrant(grantID.DatabaseName, grantID.SchemaName)
	default:
		builder = snowflake.StageGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return deleteGenericGrant(d, meta, builder)
}

type StageGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewStageGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, withGrantOption bool) *StageGrantID {
	return &StageGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *StageGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseStageGrantID(s string) (*StageGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		var roles []string
		var withGrantOption bool
		if len(idParts) == 6 {
			withGrantOption = idParts[5] == "true"
			roles = helpers.SplitStringToSlice(idParts[4], ",")
		} else {
			withGrantOption = idParts[4] == "true"
		}
		return &StageGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:           roles,
			WithGrantOption: withGrantOption,
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 6 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 6 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 6", len(idParts))
	}
	return &StageGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		WithGrantOption: idParts[4] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
		IsOldID:         false,
	}, nil
}
