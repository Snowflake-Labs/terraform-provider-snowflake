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

var validSequencePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var sequenceGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future sequences on which to grant privileges.",
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
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future sequences in the given schema. When this is true and no schema_name is provided apply this grant on all future sequences in the given database. The sequence_name field must be unset in order to use on_future.",
		Default:     false,
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future sequence.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validSequencePrivileges.ToList(), true),
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
		Description: "The name of the schema containing the current or future sequences on which to grant privileges.",
		ForceNew:    true,
	},
	"sequence_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the sequence on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// SequenceGrant returns a pointer to the resource representing a sequence grant.
func SequenceGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateSequenceGrant,
			Read:   ReadSequenceGrant,
			Delete: DeleteSequenceGrant,
			Update: UpdateSequenceGrant,

			Schema: sequenceGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					grantID, err := ParseSequenceGrantID(d.Id())
					if err != nil {
						return nil, err
					}
					if err := d.Set("sequence_name", grantID.ObjectName); err != nil {
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
		ValidPrivs: validSequencePrivileges,
	}
}

// CreateSequenceGrant implements schema.CreateFunc.
func CreateSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	var sequenceName string
	if name, ok := d.GetOk("sequence_name"); ok {
		sequenceName = name.(string)
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (sequenceName == "") && !onFuture {
		return errors.New("sequence_name must be set unless on_future is true")
	}
	if (sequenceName != "") && onFuture {
		return errors.New("sequence_name must be empty if on_future is true")
	}
	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSequenceGrant(databaseName, schemaName)
	} else {
		builder = snowflake.SequenceGrant(databaseName, schemaName, sequenceName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewSequenceGrantID(databaseName, schemaName, sequenceName, privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadSequenceGrant(d, meta)
}

// ReadSequenceGrant implements schema.ReadFunc.
func ReadSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseSequenceGrantID(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	onFuture := false
	if grantID.ObjectName == "" {
		onFuture = true
	}
	if err := d.Set("sequence_name", grantID.ObjectName); err != nil {
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
		builder = snowflake.FutureSequenceGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.SequenceGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, sequenceGrantSchema, builder, onFuture, validSequencePrivileges)
}

// DeleteSequenceGrant implements schema.DeleteFunc.
func DeleteSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseSequenceGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSequenceGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.SequenceGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateSequenceGrant implements schema.UpdateFunc.
func UpdateSequenceGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := ParseSequenceGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSequenceGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.SequenceGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
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
	return ReadSequenceGrant(d, meta)
}

type SequenceGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewSequenceGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, withGrantOption bool) *SequenceGrantID {
	return &SequenceGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *SequenceGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseSequenceGrantID(s string) (*SequenceGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &SequenceGrantID{
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
	return &SequenceGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		WithGrantOption: idParts[4] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
		IsOldID:         false,
	}, nil
}
