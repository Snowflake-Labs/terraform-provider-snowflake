package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validSequencePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
	privilegeAllPrivileges,
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
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future sequences in the given schema. When this is true and no schema_name is provided apply this grant on all future sequences in the given database. The sequence_name field must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"sequence_name"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all sequences in the given schema. When this is true and no schema_name is provided apply this grant on all sequences in the given database. The sequence_name field must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"sequence_name"},
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future sequence. To grant all privileges, use the value `ALL PRIVILEGES`",
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
	"revert_ownership_to_role_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the role to revert ownership to on destroy. Has no effect unless `privilege` is set to `OWNERSHIP`",
		Default:     "",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			additionalCharsToIgnoreValidation := []string{".", " ", ":", "(", ")"}
			return sdk.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

// SequenceGrant returns a pointer to the resource representing a sequence grant.
func SequenceGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateSequenceGrant,
			Read:               ReadSequenceGrant,
			Delete:             DeleteSequenceGrant,
			Update:             UpdateSequenceGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             sequenceGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 8 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected database_name|schema_name|sequence_name|privilege|with_grant_option|on_future|roles", d.Id())
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("sequence_name", parts[2]); err != nil {
							return nil, err
						}
					}
					if err := d.Set("privilege", parts[3]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[4])); err != nil {
						return nil, err
					}
					if err := d.Set("on_future", helpers.StringToBool(parts[5])); err != nil {
						return nil, err
					}
					if err := d.Set("on_all", helpers.StringToBool(parts[6])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[7])); err != nil {
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
	sequenceName := d.Get("sequence_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (sequenceName == "") && !onFuture && !onAll {
		return errors.New("sequence_name must be set unless on_future or on_all is true")
	}
	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSequenceGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllSequenceGrant(databaseName, schemaName)
	default:
		builder = snowflake.SequenceGrant(databaseName, schemaName, sequenceName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, sequenceName, privilege, withGrantOption, onFuture, onAll, roles)
	d.SetId(grantID)

	return ReadSequenceGrant(d, meta)
}

// ReadSequenceGrant implements schema.ReadFunc.
func ReadSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	sequenceName := d.Get("sequence_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSequenceGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllSequenceGrant(databaseName, schemaName)
	default:
		builder = snowflake.SequenceGrant(databaseName, schemaName, sequenceName)
	}

	err := readGenericGrant(d, meta, sequenceGrantSchema, builder, onFuture, onAll, validSequencePrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, sequenceName, privilege, withGrantOption, onFuture, onAll, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return err
}

// DeleteSequenceGrant implements schema.DeleteFunc.
func DeleteSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	sequenceName := d.Get("sequence_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSequenceGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllSequenceGrant(databaseName, schemaName)
	default:
		builder = snowflake.SequenceGrant(databaseName, schemaName, sequenceName)
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

	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	sequenceName := d.Get("sequence_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSequenceGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllSequenceGrant(databaseName, schemaName)
	default:
		builder = snowflake.SequenceGrant(databaseName, schemaName, sequenceName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, reversionRole, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, privilege, withGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadSequenceGrant(d, meta)
}
