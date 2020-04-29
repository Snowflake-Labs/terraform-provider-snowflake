package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

const (
	pipeIDDelimiter = '|'
)

var pipeSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the pipe; must be unique for the database and schema in which the pipe is created.",
	},
	"schema": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the pipe.",
	},
	"database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the pipe.",
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the pipe.",
	},
	"copy_statement": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the copy statement for the pipe.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if strings.TrimSuffix(old, "\n") == strings.TrimSuffix(new, "\n") {
				return true
			}
			return false
		},
	},
	"auto_ingest": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		ForceNew:    true,
		Description: "Specifies a auto_ingest param for the pipe.",
	},
	"notification_channel": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Amazon Resource Name of the Amazon SQS queue for the stage named in the DEFINITION column.",
	},
	"owner": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the pipe.",
	},
}

func Pipe() *schema.Resource {
	return &schema.Resource{
		Create: CreatePipe,
		Read:   ReadPipe,
		Update: UpdatePipe,
		Delete: DeletePipe,
		Exists: PipeExists,

		Schema: pipeSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

type pipeID struct {
	DatabaseName string
	SchemaName   string
	PipeName     string
}

//String() takes in a pipeID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|PipeName
func (si *pipeID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = pipeIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.PipeName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strPipeID := strings.TrimSpace(buf.String())
	return strPipeID, nil
}

// pipeIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|PipeName
// and returns a pipeID object
func pipeIDFromString(stringID string) (*pipeID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = pipeIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per pipe")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	pipeResult := &pipeID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		PipeName:     lines[0][2],
	}
	return pipeResult, nil
}

// CreatePipe implements schema.CreateFunc
func CreatePipe(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	schema := data.Get("schema").(string)
	name := data.Get("name").(string)

	builder := snowflake.Pipe(name, database, schema)

	// Set optionals
	if v, ok := data.GetOk("copy_statement"); ok {
		builder.WithCopyStatement(v.(string))
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("auto_ingest"); ok && v.(bool) {
		builder.WithAutoIngest()
	}

	q := builder.Create()

	err := DBExec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating pipe %v", name)
	}

	pipeID := &pipeID{
		DatabaseName: database,
		SchemaName:   schema,
		PipeName:     name,
	}
	dataIDInput, err := pipeID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadPipe(data, meta)
}

// ReadPipe implements schema.ReadFunc
func ReadPipe(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pipeID, err := pipeIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	sq := snowflake.Pipe(pipe, dbName, schema).Show()
	pipeShow, err := showPipe(db, sq)
	if err != nil {
		return err
	}

	err = data.Set("name", pipe)
	if err != nil {
		return err
	}

	err = data.Set("database", dbName)
	if err != nil {
		return err
	}

	err = data.Set("schema", schema)
	if err != nil {
		return err
	}

	err = data.Set("copy_statement", pipeShow.definition)
	if err != nil {
		return err
	}

	err = data.Set("owner", pipeShow.owner)
	if err != nil {
		return err
	}

	err = data.Set("comment", pipeShow.comment)
	if err != nil {
		return err
	}

	err = data.Set("notification_channel", pipeShow.notificationChannel)
	if err != nil {
		return err
	}

	err = data.Set("auto_ingest", pipeShow.notificationChannel != "")
	if err != nil {
		return err
	}

	return nil
}

// UpdatePipe implements schema.UpdateFunc
func UpdatePipe(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	pipeID, err := pipeIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	builder := snowflake.Pipe(pipe, dbName, schema)

	db := meta.(*sql.DB)
	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")
		q := builder.ChangeComment(comment.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating pipe comment on %v", data.Id())
		}

		data.SetPartial("comment")
	}

	return ReadPipe(data, meta)
}

// DeletePipe implements schema.DeleteFunc
func DeletePipe(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pipeID, err := pipeIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	q := snowflake.Pipe(pipe, dbName, schema).Drop()

	err = DBExec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting pipe %v", data.Id())
	}

	data.SetId("")

	return nil
}

// PipeExists implements schema.ExistsFunc
func PipeExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	pipeID, err := pipeIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	q := snowflake.Pipe(pipe, dbName, schema).Show()
	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

type showPipeResult struct {
	createdOn           string
	name                string
	databaseName        string
	schemaName          string
	definition          string
	owner               string
	notificationChannel string
	comment             string
}

func showPipe(db *sql.DB, query string) (showPipeResult, error) {
	var r showPipeResult
	row := db.QueryRow(query)
	err := row.Scan(
		&r.createdOn,
		&r.name,
		&r.databaseName,
		&r.schemaName,
		&r.definition,
		&r.owner,
		&r.notificationChannel,
		&r.comment,
	)
	if err != nil {
		return r, err
	}
	return r, nil
}
