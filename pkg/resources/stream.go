package resources

import (
	"bytes"
	"database/sql"
<<<<<<< HEAD

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"

	"encoding/csv"
	"fmt"
	"strings"
=======
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
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
<<<<<<< HEAD
=======
		ForceNew:    true,
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
		Description: "Name of the table the stream will monitor.",
	},
	"append_only": {
		Type:        schema.TypeBool,
		Optional:    true,
<<<<<<< HEAD
		Default:     false,
		Description: "Type of the stream that will be created.",
	},
=======
		ForceNew:    true,
		Default:     false,
		Description: "Type of the stream that will be created.",
	},
	"insert_only": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Create an insert only stream type.",
	},
	"show_initial_rows": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Specifies whether to return all existing rows in the source table as row inserts the first time the stream is consumed.",
	},
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
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
<<<<<<< HEAD
			State: schema.ImportStatePassthrough,
=======
			StateContext: schema.ImportStatePassthroughContext,
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
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
<<<<<<< HEAD
func CreateStream(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	schema := data.Get("schema").(string)
	name := data.Get("name").(string)
	onTable := data.Get("on_table").(string)
	appendOnly := data.Get("append_only").(bool)
=======
func CreateStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	onTable := d.Get("on_table").(string)
	appendOnly := d.Get("append_only").(bool)
	insertOnly := d.Get("insert_only").(bool)
	showInitialRows := d.Get("show_initial_rows").(bool)
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c

	builder := snowflake.Stream(name, database, schema)

	resultOnTable, err := streamOnTableIDFromString(onTable)
	if err != nil {
		return err
	}

	builder.WithOnTable(resultOnTable.DatabaseName, resultOnTable.SchemaName, resultOnTable.OnTableName)
	builder.WithAppendOnly(appendOnly)
<<<<<<< HEAD

	// Set optionals
	if v, ok := data.GetOk("comment"); ok {
=======
	builder.WithInsertOnly(insertOnly)
	builder.WithShowInitialRows(showInitialRows)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
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
<<<<<<< HEAD
	data.SetId(dataIDInput)

	return ReadStream(data, meta)
}

// ReadStream implements schema.ReadFunc
func ReadStream(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(data.Id())
=======
	d.SetId(dataIDInput)

	return ReadStream(d, meta)
}

// ReadStream implements schema.ReadFunc
func ReadStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(d.Id())
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	name := streamID.StreamName

	stmt := snowflake.Stream(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, stmt)
	stream, err := snowflake.ScanStream(row)
<<<<<<< HEAD
=======
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

	err = d.Set("database", stream.DatabaseName.String)
	if err != nil {
		return err
	}

	err = d.Set("schema", stream.SchemaName.String)
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	if err != nil {
		return err
	}

<<<<<<< HEAD
	err = data.Set("name", stream.StreamName.String)
=======
	err = d.Set("on_table", stream.TableName.String)
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	if err != nil {
		return err
	}

<<<<<<< HEAD
	err = data.Set("owner", stream.Owner.String)
=======
	err = d.Set("append_only", stream.Mode.String == "APPEND_ONLY")
	if err != nil {
		return err
	}

	err = d.Set("insert_only", stream.InsertOnly)
	if err != nil {
		return err
	}

	err = d.Set("show_initial_rows", stream.ShowInitialRows)
	if err != nil {
		return err
	}

	err = d.Set("comment", stream.Comment.String)
	if err != nil {
		return err
	}

	err = d.Set("owner", stream.Owner.String)
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	if err != nil {
		return err
	}

	return nil
}

// DeleteStream implements schema.DeleteFunc
<<<<<<< HEAD
func DeleteStream(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(data.Id())
=======
func DeleteStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(d.Id())
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	streamName := streamID.StreamName

	q := snowflake.Stream(streamName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
<<<<<<< HEAD
		return errors.Wrapf(err, "error deleting stream %v", data.Id())
	}

	data.SetId("")
=======
		return errors.Wrapf(err, "error deleting stream %v", d.Id())
	}

	d.SetId("")
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c

	return nil
}

// UpdateStream implements schema.UpdateFunc
<<<<<<< HEAD
func UpdateStream(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	streamID, err := streamIDFromString(data.Id())
=======
func UpdateStream(d *schema.ResourceData, meta interface{}) error {
	streamID, err := streamIDFromString(d.Id())
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	streamName := streamID.StreamName

	builder := snowflake.Stream(streamName, dbName, schema)

	db := meta.(*sql.DB)
<<<<<<< HEAD
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
=======
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stream comment on %v", d.Id())
		}
	}

	return ReadStream(d, meta)
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
}
