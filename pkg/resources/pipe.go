package resources

import (
	"database/sql"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
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
		Description: "Specifies the Amazon Resource Name (ARN) for the SNS topic for your S3 bucket.",
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

	q := builder.Create()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating pipe %v", name)
	}

	pipeID := &schemaScopedID{
		Database: database,
		Schema:   schema,
		Name:     name,
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
	pipeID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.Database
	schema := pipeID.Schema
	name := pipeID.Name

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

	err = d.Set("auto_ingest", pipe.NotificationChannel != "")
	if err != nil {
		return err
	}

	if strings.Contains(pipe.NotificationChannel, "arn:aws:sns:") {
		err = d.Set("aws_sns_topic_arn", pipe.NotificationChannel)
		return err
	}

	return nil
}

// UpdatePipe implements schema.UpdateFunc
func UpdatePipe(d *schema.ResourceData, meta interface{}) error {
	pipeID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.Database
	schema := pipeID.Schema
	pipe := pipeID.Name

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

	return ReadPipe(d, meta)
}

// DeletePipe implements schema.DeleteFunc
func DeletePipe(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	pipeID, err := idFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := pipeID.Database
	schema := pipeID.Schema
	pipe := pipeID.Name

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
	pipeID, err := idFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := pipeID.Database
	schema := pipeID.Schema
	pipe := pipeID.Name

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
