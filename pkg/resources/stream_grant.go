package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validStreamPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeSelect,
)

var streamGrantSchema = map[string]*schema.Schema{
	"stream_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the stream on which to grant privileges immediately (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future streams on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future streams on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future stream.",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(validStreamPrivileges.toList(), true),
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
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future streams in the given schema. When this is true and no schema_name is provided apply this grant on all future streams in the given database. The stream_name field must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"stream_name"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// StreamGrant returns a pointer to the resource representing a stream grant
func StreamGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateStreamGrant,
		Read:   ReadStreamGrant,
		Delete: DeleteStreamGrant,

		Schema: streamGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateStreamGrant implements schema.CreateFunc
func CreateStreamGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		streamName string
		schemaName string
	)
	if _, ok := data.GetOk("stream_name"); ok {
		streamName = data.Get("stream_name").(string)
	}
	if _, ok := data.GetOk("schema_name"); ok {
		schemaName = data.Get("schema_name").(string)
	}
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	futureStreams := data.Get("on_future").(bool)
	grantOption := data.Get("with_grant_option").(bool)

	if (schemaName == "") && !futureStreams {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	if (streamName == "") && !futureStreams {
		return errors.New("stream_name must be set unless on_future is true.")
	}
	if (streamName != "") && futureStreams {
		return errors.New("stream_name must be empty if on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if futureStreams {
		builder = snowflake.FutureStreamGrant(dbName, schemaName)
	} else {
		builder = snowflake.StreamGrant(dbName, schemaName, streamName)
	}

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   streamName,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadStreamGrant(data, meta)
}

// ReadStreamGrant implements schema.ReadFunc
func ReadStreamGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	streamName := grantID.ObjectName
	priv := grantID.Privilege

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureStreamsEnabled := false
	if streamName == "" {
		futureStreamsEnabled = true
	}
	err = data.Set("stream_name", streamName)
	if err != nil {
		return err
	}
	err = data.Set("on_future", futureStreamsEnabled)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = data.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureStreamsEnabled {
		builder = snowflake.FutureStreamGrant(dbName, schemaName)
	} else {
		builder = snowflake.StreamGrant(dbName, schemaName, streamName)
	}

	return readGenericGrant(data, meta, streamGrantSchema, builder, futureStreamsEnabled, validStreamPrivileges)
}

// DeleteStreamGrant implements schema.DeleteFunc
func DeleteStreamGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	streamName := grantID.ObjectName

	futureStreams := (streamName == "")

	var builder snowflake.GrantBuilder
	if futureStreams {
		builder = snowflake.FutureStreamGrant(dbName, schemaName)
	} else {
		builder = snowflake.StreamGrant(dbName, schemaName, streamName)
	}
	return deleteGenericGrant(data, meta, builder)
}
