package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
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

	// trim off any trailing line endings and leading/trailing whitespace
	return strings.TrimSpace(strings.TrimRight(o, ";\r\n")) == strings.TrimSpace(strings.TrimRight(n, ";\r\n"))
}

// CreatePipe implements schema.CreateFunc.
func CreatePipe(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)

	ctx := context.Background()
	objectIdentifier := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	opts := &sdk.CreatePipeOptions{}

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
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	ctx := context.Background()
	pipe, err := client.Pipes.ShowByID(ctx, id)
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] pipe (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
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

	if err := d.Set("auto_ingest", pipe.NotificationChannel != ""); err != nil {
		return err
	}

	if strings.Contains(pipe.NotificationChannel, "arn:aws:sns:") {
		if err := d.Set("aws_sns_topic_arn", pipe.NotificationChannel); err != nil {
			return err
		}
	}

	if err := d.Set("error_integration", pipe.ErrorIntegration); err != nil {
		return err
	}

	return nil
}

// UpdatePipe implements schema.UpdateFunc.
func UpdatePipe(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	ctx := context.Background()

	pipeSet := &sdk.PipeSet{}
	pipeUnset := &sdk.PipeUnset{}
	var runSetStatement bool
	var runUnsetStatement bool

	if d.HasChange("comment") {
		if comment, ok := d.GetOk("comment"); ok {
			runSetStatement = true
			pipeSet.Comment = sdk.String(comment.(string))
		} else {
			runUnsetStatement = true
			pipeUnset.Comment = sdk.Bool(true)
		}
	}

	if d.HasChange("error_integration") {
		if errorIntegration, ok := d.GetOk("error_integration"); ok {
			runSetStatement = true
			pipeSet.ErrorIntegration = sdk.String(errorIntegration.(string))
		} else {
			runUnsetStatement = true
			pipeUnset.ErrorIntegration = sdk.Bool(true)
		}
	}

	if runSetStatement {
		options := &sdk.AlterPipeOptions{Set: pipeSet}
		err := client.Pipes.Alter(ctx, objectIdentifier, options)
		if err != nil {
			return fmt.Errorf("error updating pipe %v: %w", objectIdentifier.Name(), err)
		}
	}

	if runUnsetStatement {
		options := &sdk.AlterPipeOptions{Unset: pipeUnset}
		err := client.Pipes.Alter(ctx, objectIdentifier, options)
		if err != nil {
			return fmt.Errorf("error updating pipe %v: %w", objectIdentifier.Name(), err)
		}
	}

	return ReadPipe(d, meta)
}

// DeletePipe implements schema.DeleteFunc.
func DeletePipe(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Pipes.Drop(ctx, objectIdentifier, &sdk.DropPipeOptions{IfExists: sdk.Bool(true)})
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
