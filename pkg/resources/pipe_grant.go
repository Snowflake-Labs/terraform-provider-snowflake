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

var validPipePrivileges = NewPrivilegeSet(
	privilegeMonitor,
	privilegeOperate,
	privilegeOwnership,
)

var pipeGrantSchema = map[string]*schema.Schema{
	"pipe_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the pipe on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future pipes on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future pipes on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future pipe.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validPipePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future pipes in the given schema. When this is true and no schema_name is provided apply this grant on all future pipes in the given database. The pipe_name field must be unset in order to use on_future.",
		Default:     false,
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
		ForceNew:    true,
	},
}

// PipeGrant returns a pointer to the resource representing a pipe grant.
func PipeGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreatePipeGrant,
			Read:   ReadPipeGrant,
			Delete: DeletePipeGrant,
			Update: UpdatePipeGrant,

			Schema: pipeGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validPipePrivileges,
	}
}

// CreatePipeGrant implements schema.CreateFunc.
func CreatePipeGrant(d *schema.ResourceData, meta interface{}) error {
	var pipeName string
	if name, ok := d.GetOk("pipe_name"); ok {
		pipeName = name.(string)
	}
	var schemaName string
	if name, ok := d.GetOk("schema_name"); ok {
		schemaName = name.(string)
	}

	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}
	if (pipeName == "") && !onFuture {
		return errors.New("pipe_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FuturePipeGrant(databaseName, schemaName)
	} else {
		builder = snowflake.PipeGrant(databaseName, schemaName, pipeName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewPipeGrantID(databaseName, schemaName, pipeName, privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadPipeGrant(d, meta)
}

// ReadPipeGrant implements schema.ReadFunc.
func ReadPipeGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parsePipeGrantID(d.Id())
	if err != nil {
		return err
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

	onFuture := (grantID.ObjectName == "")

	if err := d.Set("pipe_name", grantID.ObjectName); err != nil {
		return err
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
		builder = snowflake.FuturePipeGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.PipeGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, pipeGrantSchema, builder, onFuture, validPipePrivileges)
}

// DeletePipeGrant implements schema.DeleteFunc.
func DeletePipeGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parsePipeGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FuturePipeGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.PipeGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdatePipeGrant implements schema.UpdateFunc.
func UpdatePipeGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update and we're done
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}

	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}

	grantID, err := parsePipeGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FuturePipeGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.PipeGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
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
	return ReadPipeGrant(d, meta)
}

type PipeGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewPipeGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, withGrantOption bool) *PipeGrantID {
	return &PipeGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *PipeGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func parsePipeGrantID(s string) (*PipeGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &PipeGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			WithGrantOption: idParts[5] == "true",
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
	return &PipeGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		WithGrantOption: idParts[4] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
		IsOldID:         false,
	}, nil
}
