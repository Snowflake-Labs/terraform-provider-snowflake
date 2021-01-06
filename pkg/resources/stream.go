package resources

import (
	"database/sql"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var streamSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the stream; must be unique for the database and schema in which the stream is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the stream.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the stream.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stream.",
	},
	"on_table": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of the table the stream will monitor.",
	},
	"append_only": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Type of the stream that will be created.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the stream.",
	},
}

func Stream() *schema.Resource {
	return &schema.Resource{
		Create: CreateStream,
		Read:   ReadStream,
		Update: UpdateStream,
		Delete: DeleteStream,

		Schema: streamSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateStream implements schema.CreateFunc
func CreateStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	onTable := d.Get("on_table").(string)
	appendOnly := d.Get("append_only").(bool)

	builder := snowflake.Stream(name, database, schema)

	resultOnTable, err := streamOnTableIDFromString(onTable)
	if err != nil {
		return err
	}

	builder.WithOnTable(resultOnTable.DatabaseName, resultOnTable.SchemaName, resultOnTable.OnTableName)
	builder.WithAppendOnly(appendOnly)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	stmt := builder.Create()
	err = snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating stream %v", name)
	}

	streamID := &schemaScopedID{
		Database: database,
		Schema:   schema,
		Name:     name,
	}
	dataIDInput, err := streamID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadStream(d, meta)
}

// ReadStream implements schema.ReadFunc
func ReadStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := streamID.Database
	schema := streamID.Schema
	name := streamID.Name

	stmt := snowflake.Stream(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, stmt)
	stream, err := snowflake.ScanStream(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] stream (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", stream.StreamName.String)
	if err != nil {
		return err
	}

	err = d.Set("owner", stream.Owner.String)
	if err != nil {
		return err
	}

	return nil
}

// DeleteStream implements schema.DeleteFunc
func DeleteStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := streamID.Database
	schema := streamID.Schema
	streamName := streamID.Name

	q := snowflake.Stream(streamName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting stream %v", d.Id())
	}

	d.SetId("")

	return nil
}

// UpdateStream implements schema.UpdateFunc
func UpdateStream(d *schema.ResourceData, meta interface{}) error {
	streamID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := streamID.Database
	schema := streamID.Schema
	streamName := streamID.Name

	builder := snowflake.Stream(streamName, dbName, schema)

	db := meta.(*sql.DB)
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stream comment on %v", d.Id())
		}
	}

	return ReadStream(d, meta)
}
