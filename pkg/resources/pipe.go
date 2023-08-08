package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	"integration": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies an integration for the pipe.",
	},
	"notification_channel": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Amazon Resource Name of the Amazon SQS queue for the stage named in the DEFINITION column.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the pipe.",
	},
	"error_integration": {
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

func pipeCopyStatementDiffSuppress(_, o, n string, _ *schema.ResourceData) bool {
	// standardize line endings
	o = strings.ReplaceAll(o, "\r\n", "\n")
	n = strings.ReplaceAll(n, "\r\n", "\n")

	// trim off any trailing line endings
	return strings.TrimRight(o, ";\r\n") == strings.TrimRight(n, ";\r\n")
}

type pipeID struct {
	DatabaseName string
	SchemaName   string
	PipeName     string
}

// String() takes in a pipeID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|PipeName.
func (si *pipeID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = pipeIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.PipeName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strPipeID := strings.TrimSpace(buf.String())
	return strPipeID, nil
}

// pipeIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|PipeName
// and returns a pipeID object.
func pipeIDFromString(stringID string) (*pipeID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = pipeIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
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

// CreatePipe implements schema.CreateFunc.
func CreatePipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)

	ctx := context.Background()
	objectIdentifier := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	opts := &sdk.PipeCreateOptions{}

	copyStatement := d.Get("copy_statement").(string)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	if v, ok := d.GetOk("auto_ingest"); ok && v.(bool) {
		opts.AutoIngest = sdk.Bool(true)
	}

	if v, ok := d.GetOk("aws_sns_topic_arn"); ok {
		opts.AwsSnsTopic = sdk.String(v.(string))
	}

	if v, ok := d.GetOk("integration"); ok {
		opts.Integration = sdk.String(v.(string))
	}

	if v, ok := d.GetOk("error_integration"); ok {
		opts.ErrorIntegration = sdk.String(v.(string))
	}

	err := client.Pipes.Create(ctx, objectIdentifier, copyStatement, opts)
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	return ReadPipe(d, meta)
}

// ReadPipe implements schema.ReadFunc.
func ReadPipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pipeID, err := pipeIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	name := pipeID.PipeName

	sq := snowflake.NewPipeBuilder(name, dbName, schema).Show()
	row := snowflake.QueryRow(db, sq)
	pipe, err := snowflake.ScanPipe(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] pipe (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if err := d.Set("name", pipe.Name); err != nil {
		return err
	}

	if err := d.Set("database", pipe.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("schema", pipe.SchemaName); err != nil {
		return err
	}

	if err := d.Set("copy_statement", pipe.Definition); err != nil {
		return err
	}

	if err := d.Set("owner", pipe.Owner); err != nil {
		return err
	}

	if err := d.Set("comment", pipe.Comment); err != nil {
		return err
	}

	if err := d.Set("notification_channel", pipe.NotificationChannel); err != nil {
		return err
	}

	if err := d.Set("auto_ingest", pipe.NotificationChannel != nil); err != nil {
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
	return err
}

// UpdatePipe implements schema.UpdateFunc.
func UpdatePipe(d *schema.ResourceData, meta interface{}) error {
	pipeID, err := pipeIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	builder := snowflake.NewPipeBuilder(pipe, dbName, schema)

	db := meta.(*sql.DB)
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating pipe comment on %v", d.Id())
		}
	}

	if d.HasChange("error_integration") {
		var q string
		if errorIntegration, ok := d.GetOk("error_integration"); ok {
			q = builder.ChangeErrorIntegration(errorIntegration.(string))
		} else {
			q = builder.RemoveErrorIntegration()
		}
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating pipe error_integration on %v", d.Id())
		}
	}

	return ReadPipe(d, meta)
}

// DeletePipe implements schema.DeleteFunc.
func DeletePipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pipeID, err := pipeIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.DatabaseName
	schema := pipeID.SchemaName
	pipe := pipeID.PipeName

	q := snowflake.NewPipeBuilder(pipe, dbName, schema).Drop()

	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting pipe %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
