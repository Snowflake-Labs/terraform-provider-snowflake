package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validSequencePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var sequenceGrantSchema = map[string]*schema.Schema{
	"sequence_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the sequence on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current or future sequences on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future sequences on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future sequence.",
		Default:      "USAGE",
		ValidateFunc: validation.ValidatePrivilege(validSequencePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future sequences in the given schema. When this is true and no schema_name is provided apply this grant on all future sequences in the given database. The sequence_name field must be unset in order to use on_future.",
		Default:     false,
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

// SequenceGrant returns a pointer to the resource representing a sequence grant
func SequenceGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateSequenceGrant,
			Read:   ReadSequenceGrant,
			Delete: DeleteSequenceGrant,

			Schema: sequenceGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validSequencePrivileges,
	}
}

// CreateSequenceGrant implements schema.CreateFunc
func CreateSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	var sequenceName string
	if name, ok := d.GetOk("sequence_name"); ok {
		sequenceName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futureSequences := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (sequenceName == "") && !futureSequences {
		return errors.New("sequence_name must be set unless on_future is true.")
	}
	if (sequenceName != "") && futureSequences {
		return errors.New("sequence_name must be empty if on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if futureSequences {
		builder = snowflake.FutureSequenceGrant(dbName, schemaName)
	} else {
		builder = snowflake.SequenceGrant(dbName, schemaName, sequenceName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   sequenceName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadSequenceGrant(d, meta)
}

// ReadSequenceGrant implements schema.ReadFunc
func ReadSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	sequenceName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureSequencesEnabled := false
	if sequenceName == "" {
		futureSequencesEnabled = true
	}
	err = d.Set("sequence_name", sequenceName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureSequencesEnabled)
	if err != nil {
		return err
	}
	err = d.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureSequencesEnabled {
		builder = snowflake.FutureSequenceGrant(dbName, schemaName)
	} else {
		builder = snowflake.SequenceGrant(dbName, schemaName, sequenceName)
	}

	return readGenericGrant(d, meta, sequenceGrantSchema, builder, futureSequencesEnabled, validSequencePrivileges)
}

// DeleteSequenceGrant implements schema.DeleteFunc
func DeleteSequenceGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	sequenceName := grantID.ObjectName

	futureSequences := (sequenceName == "")

	var builder snowflake.GrantBuilder
	if futureSequences {
		builder = snowflake.FutureSequenceGrant(dbName, schemaName)
	} else {
		builder = snowflake.SequenceGrant(dbName, schemaName, sequenceName)
	}
	return deleteGenericGrant(d, meta, builder)
}
