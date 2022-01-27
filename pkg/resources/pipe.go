package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const (
	pipeIDDelimiter = '|'
)

var pipeSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the pipe; must be unique for the database and schema in which the pipe is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the pipe.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the pipe.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the pipe.",
	},
	"copy_statement": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the copy statement for the pipe.",
		DiffSuppressFunc: pipeCopyStatementDiffSuppress,
	},
	"auto_ingest": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		ForceNew:    true,
		Description: "Specifies a auto_ingest param for the pipe.",
	},
	"aws_sns_topic_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the Amazon Resource Name (ARN) for the SNS topic for your S3 bucket.",
	},
	"integration": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies an integration for the pipe.",
	},
	"notification_channel": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Amazon Resource Name of the Amazon SQS queue for the stage named in the DEFINITION column.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the pipe.",
	},
	"error_integration": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the name of the notification integration used for error notifications.",
	},
}

func Pipe() *schema.Resource {
	return &schema.Resource{
		Create: CreatePipe,
		Read:   ReadPipe,
		Update: UpdatePipe,
		Delete: DeletePipe,

		Schema: pipeSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func pipeCopyStatementDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// standardise line endings
	old = strings.ReplaceAll(old, "\r\n", "\n")
	new = strings.ReplaceAll(new, "\r\n", "\n")

	// trim off any trailing line endings
	return strings.TrimRight(old, ";\r\n") == strings.TrimRight(new, ";\r\n")
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
func CreatePipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)

	builder := snowflake.Pipe(name, database, schema)

	// Set optionals
	if v, ok := d.GetOk("copy_statement"); ok {
		builder.WithCopyStatement(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("auto_ingest"); ok && v.(bool) {
		builder.WithAutoIngest()
	}

	if v, ok := d.GetOk("aws_sns_topic_arn"); ok {
		builder.WithAwsSnsTopicArn(v.(string))
	}

	if v, ok := d.GetOk("integration"); ok {
		builder.WithIntegration(v.(string))
	}

	if v, ok := d.GetOk("error_integration"); ok {
		builder.WithErrorIntegration((v.(string)))
	}

	q := builder.Create()

	err := snowflake.Exec(db, q)
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
	d.SetId(dataIDInput)

	return ReadPipe(d, meta)
}

// ReadPipe implements schema.ReadFunc
func ReadPipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pipeID, err := pipeIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	name := pipeID.PipeName

	sq := snowflake.Pipe(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, sq)
	pipe, err := snowflake.ScanPipe(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] pipe (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", pipe.Name)
	if err != nil {
		return err
	}

	err = d.Set("database", pipe.DatabaseName)
	if err != nil {
		return err
	}

	err = d.Set("schema", pipe.SchemaName)
	if err != nil {
		return err
	}

	err = d.Set("copy_statement", pipe.Definition)
	if err != nil {
		return err
	}

	err = d.Set("owner", pipe.Owner)
	if err != nil {
		return err
	}

	err = d.Set("comment", pipe.Comment)
	if err != nil {
		return err
	}

	err = d.Set("notification_channel", pipe.NotificationChannel)
	if err != nil {
		return err
	}

	err = d.Set("auto_ingest", pipe.NotificationChannel != nil)
	if err != nil {
		return err
	}

	if pipe.NotificationChannel != nil && strings.Contains(*pipe.NotificationChannel, "arn:aws:sns:") {
		err = d.Set("aws_sns_topic_arn", pipe.NotificationChannel)
		return err
	}

	// The "DESCRIBE PIPE ..." command returns the string "null" for error_integration
	if pipe.ErrorIntegration.String == "null" {
		pipe.ErrorIntegration.Valid = false
		pipe.ErrorIntegration.String = ""
	}
	err = d.Set("error_integration", pipe.ErrorIntegration.String)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePipe implements schema.UpdateFunc
func UpdatePipe(d *schema.ResourceData, meta interface{}) error {
	pipeID, err := pipeIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	builder := snowflake.Pipe(pipe, dbName, schema)

	db := meta.(*sql.DB)
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating pipe comment on %v", d.Id())
		}
	}

	if d.HasChange("error_integration") {
		var q string
		if errorIntegration, ok := d.GetOk("error_integration"); ok {
			q = builder.ChangeErrorIntegration(errorIntegration.(string))
		} else {
			q = builder.RemoveErrorIntegration()
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating pipe error_integration on %v", d.Id())
		}
	}

	return ReadPipe(d, meta)
}

// DeletePipe implements schema.DeleteFunc
func DeletePipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pipeID, err := pipeIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	q := snowflake.Pipe(pipe, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting pipe %v", d.Id())
	}

	d.SetId("")

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
