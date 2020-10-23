package resources

import (
	"bytes"
	"database/sql"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"

	"encoding/csv"
	"fmt"
	"strings"
)

const (
	streamIDDelimiter        = '|'
	streamOnTableIDDelimiter = '.'
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
			State: schema.ImportStatePassthrough,
		},
	}
}

type streamID struct {
	DatabaseName string
	SchemaName   string
	StreamName   string
}

type streamOnTableID struct {
	DatabaseName string
	SchemaName   string
	OnTableName  string
}

//String() takes in a streamID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|StreamName
func (si *streamID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = streamIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.StreamName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strStreamID := strings.TrimSpace(buf.String())
	return strStreamID, nil
}

// streamIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|StreamName
// and returns a streamID object
func streamIDFromString(stringID string) (*streamID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = streamIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	streamResult := &streamID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		StreamName:   lines[0][2],
	}
	return streamResult, nil
}

// streamOnTableIDFromString() takes in a dot-delimited string: DatabaseName.SchemaName.TableName
// and returns a streamOnTableID object
func streamOnTableIDFromString(stringID string) (*streamOnTableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = streamOnTableIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		//return nil, fmt.Errorf("on table format: database_name.schema_name.target_table_name")
		return nil, fmt.Errorf("invalid format for on_table: %v , expected: <database_name.schema_name.target_table_name>", strings.Join(lines[0], "."))
	}

	streamOnTableResult := &streamOnTableID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		OnTableName:  lines[0][2],
	}
	return streamOnTableResult, nil
}

// CreateStream implements schema.CreateFunc
func CreateStream(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	schema := data.Get("schema").(string)
	name := data.Get("name").(string)
	onTable := data.Get("on_table").(string)
	appendOnly := data.Get("append_only").(bool)

	builder := snowflake.Stream(name, database, schema)

	resultOnTable, err := streamOnTableIDFromString(onTable)
	if err != nil {
		return err
	}

	builder.WithOnTable(resultOnTable.DatabaseName, resultOnTable.SchemaName, resultOnTable.OnTableName)
	builder.WithAppendOnly(appendOnly)

	// Set optionals
	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	stmt := builder.Create()
	err = snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating stream %v", name)
	}

	streamID := &streamID{
		DatabaseName: database,
		SchemaName:   schema,
		StreamName:   name,
	}
	dataIDInput, err := streamID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadStream(data, meta)
}

// ReadStream implements schema.ReadFunc
func ReadStream(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	name := streamID.StreamName

	stmt := snowflake.Stream(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, stmt)
	stream, err := snowflake.ScanStream(row)
	if err != nil {
		return err
	}

	err = data.Set("name", stream.StreamName.String)
	if err != nil {
		return err
	}

	err = data.Set("owner", stream.Owner.String)
	if err != nil {
		return err
	}

	return nil
}

// DeleteStream implements schema.DeleteFunc
func DeleteStream(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	streamName := streamID.StreamName

	q := snowflake.Stream(streamName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting stream %v", data.Id())
	}

	data.SetId("")

	return nil
}

// UpdateStream implements schema.UpdateFunc
func UpdateStream(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	streamID, err := streamIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	streamName := streamID.StreamName

	builder := snowflake.Stream(streamName, dbName, schema)

	db := meta.(*sql.DB)
	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stream comment on %v", data.Id())
		}

		data.SetPartial("comment")
	}

	return ReadStream(data, meta)
}
