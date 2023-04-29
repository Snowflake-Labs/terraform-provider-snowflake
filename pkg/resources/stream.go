package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	streamIDDelimiter         = '|'
	streamOnObjectIDDelimiter = '.'
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
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Name of the table the stream will monitor.",
		ExactlyOneOf: []string{"on_table", "on_view", "on_stage"},
	},
	"on_view": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Name of the view the stream will monitor.",
		ExactlyOneOf: []string{"on_table", "on_view", "on_stage"},
	},
	"on_stage": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Name of the stage the stream will monitor.",
		ExactlyOneOf: []string{"on_table", "on_view", "on_stage"},
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// Suppress diff if the stage name is the same, even if database and schema are not specified
			return strings.Trim(strings.Split(old, ".")[len(strings.Split(old, "."))-1], "\"") == strings.Trim(strings.Split(new, ".")[len(strings.Split(new, "."))-1], "\"")
		},
	},
	"append_only": {
		Type:        schema.TypeBool,
		Optional:    true,
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

type streamID struct {
	DatabaseName string
	SchemaName   string
	StreamName   string
}

type streamOnObjectID struct {
	DatabaseName string
	SchemaName   string
	Name         string
}

// String() takes in a streamID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|StreamName.
func (si *streamID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = streamIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.StreamName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strStreamID := strings.TrimSpace(buf.String())
	return strStreamID, nil
}

// streamIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|StreamName
// and returns a streamID object.
func streamIDFromString(stringID string) (*streamID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = streamIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
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

// streamOnObjectIDFromString() takes in a dot-delimited string: DatabaseName.SchemaName.TableName
// and returns a streamOnObjectID object.
func streamOnObjectIDFromString(stringID string) (*streamOnObjectID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = streamOnObjectIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		// return nil, fmt.Errorf("on table format: database_name.schema_name.target_table_name")
		return nil, fmt.Errorf("invalid format for on_table: %v , expected: <database_name.schema_name.target_table_name>", strings.Join(lines[0], "."))
	}

	streamOnTableResult := &streamOnObjectID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		Name:         lines[0][2],
	}
	return streamOnTableResult, nil
}

// CreateStream implements schema.CreateFunc.
func CreateStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	appendOnly := d.Get("append_only").(bool)
	insertOnly := d.Get("insert_only").(bool)
	showInitialRows := d.Get("show_initial_rows").(bool)

	builder := snowflake.Stream(name, database, schema)

	onTable, onTableSet := d.GetOk("on_table")
	onView, onViewSet := d.GetOk("on_view")
	onStage, onStageSet := d.GetOk("on_stage")

	switch {
	case onTableSet:
		id, err := streamOnObjectIDFromString(onTable.(string))
		if err != nil {
			return err
		}

		tq := snowflake.NewTableBuilder(id.Name, id.DatabaseName, id.SchemaName).Show()
		tableRow := snowflake.QueryRow(db, tq)

		t, err := snowflake.ScanTable(tableRow)
		if err != nil {
			return err
		}

		builder.WithExternalTable(t.IsExternal.String == "Y")
		builder.WithOnTable(t.DatabaseName.String, t.SchemaName.String, t.TableName.String)
	case onViewSet:
		id, err := streamOnObjectIDFromString(onView.(string))
		if err != nil {
			return err
		}

		tq := snowflake.NewViewBuilder(id.Name).WithDB(id.DatabaseName).WithSchema(id.SchemaName).Show()
		viewRow := snowflake.QueryRow(db, tq)

		t, err := snowflake.ScanView(viewRow)
		if err != nil {
			return err
		}

		builder.WithOnView(t.DatabaseName.String, t.SchemaName.String, t.Name.String)
	case onStageSet:
		id, err := streamOnObjectIDFromString(onStage.(string))
		if err != nil {
			return err
		}
		stageBuilder := snowflake.NewStageBuilder(id.Name, id.DatabaseName, id.SchemaName)
		sq := stageBuilder.Describe()
		d, err := snowflake.DescStage(db, sq)
		if err != nil {
			return err
		}
		if !strings.Contains(d.Directory, "ENABLE = true") {
			return fmt.Errorf("directory must be enabled on stage")
		}

		builder.WithOnStage(id.DatabaseName, id.SchemaName, id.Name)
	}

	builder.WithAppendOnly(appendOnly)
	builder.WithInsertOnly(insertOnly)
	builder.WithShowInitialRows(showInitialRows)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	stmt := builder.Create()
	if err := snowflake.Exec(db, stmt); err != nil {
		return fmt.Errorf("error creating stream %v", name)
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
	d.SetId(dataIDInput)

	return ReadStream(d, meta)
}

// ReadStream implements schema.ReadFunc.
func ReadStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	name := streamID.StreamName

	stmt := snowflake.Stream(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, stmt)
	stream, err := snowflake.ScanStream(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] stream (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", stream.StreamName.String); err != nil {
		return err
	}

	if err := d.Set("database", stream.DatabaseName.String); err != nil {
		return err
	}

	if err := d.Set("schema", stream.SchemaName.String); err != nil {
		return err
	}

	if stream.SourceType.String == "Stage" {
		if err := d.Set("on_stage", stream.TableName.String); err != nil {
			return err
		}
	} else if stream.SourceType.String == "View" {
		if err := d.Set("on_view", stream.TableName.String); err != nil {
			return err
		}
	} else {
		if err := d.Set("on_table", stream.TableName.String); err != nil {
			return err
		}
	}

	if err := d.Set("append_only", stream.Mode.String == "APPEND_ONLY"); err != nil {
		return err
	}

	if err := d.Set("insert_only", stream.Mode.String == "INSERT_ONLY"); err != nil {
		return err
	}

	if err := d.Set("show_initial_rows", stream.ShowInitialRows); err != nil {
		return err
	}

	if err := d.Set("comment", stream.Comment.String); err != nil {
		return err
	}

	if err := d.Set("owner", stream.Owner.String); err != nil {
		return err
	}

	return nil
}

// DeleteStream implements schema.DeleteFunc.
func DeleteStream(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	streamID, err := streamIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	streamName := streamID.StreamName

	q := snowflake.Stream(streamName, dbName, schema).Drop()

	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting stream %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}

// UpdateStream implements schema.UpdateFunc.
func UpdateStream(d *schema.ResourceData, meta interface{}) error {
	streamID, err := streamIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := streamID.DatabaseName
	schema := streamID.SchemaName
	streamName := streamID.StreamName

	builder := snowflake.Stream(streamName, dbName, schema)

	db := meta.(*sql.DB)
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating stream comment on %v", d.Id())
		}
	}

	return ReadStream(d, meta)
}
