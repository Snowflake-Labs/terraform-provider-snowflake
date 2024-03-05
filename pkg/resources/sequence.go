package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sequenceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the name for the sequence.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies a comment for the sequence.",
	},
	"increment": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     1,
		Description: "The amount the sequence will increase by each time it is used",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the sequence. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the sequence. Don't use the | character.",
		ForceNew:    true,
	},
	"next_value": {
		Type:        schema.TypeInt,
		Description: "The increment sequence interval.",
		Computed:    true,
		ForceNew:    true,
	},
	"ordering": {
		Type:        schema.TypeString,
		Description: "The ordering of the sequence. Either ORDER or NOORDER. Default is ORDER.",
		Optional:    true,
		Default:     "ORDER",
		ValidateDiagFunc: StringInSlice(
			[]string{
				string(sdk.ValuesBehaviorNoOrder),
				string(sdk.ValuesBehaviorOrder),
			}, false),
	},
	"fully_qualified_name": {
		Type:        schema.TypeString,
		Description: "The fully qualified name of the sequence.",
		Computed:    true,
	},
}

func Sequence() *schema.Resource {
	return &schema.Resource{
		Create: CreateSequence,
		Read:   ReadSequence,
		Delete: DeleteSequence,
		Update: UpdateSequence,

		Schema: sequenceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateSequence(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)
	req := sdk.NewCreateSequenceRequest(id)

	if v, ok := d.GetOk("increment"); ok {
		req.WithIncrement(sdk.Int(v.(int)))
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("ordering"); ok {
		req.WithValuesBehavior(sdk.ValuesBehaviorPointer(sdk.ValuesBehavior(v.(string))))
	}
	err := client.Sequences.Create(ctx, req)
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(database, schema, name))

	return ReadSequence(d, meta)
}

func ReadSequence(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	seq, err := client.Sequences.ShowByID(ctx, id)
	if err != nil {
		return err
	}

	if err := d.Set("name", seq.Name); err != nil {
		return err
	}

	if err := d.Set("schema", seq.SchemaName); err != nil {
		return err
	}

	if err := d.Set("database", seq.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("comment", seq.Comment); err != nil {
		return err
	}

	if err := d.Set("increment", seq.Interval); err != nil {
		return err
	}

	if err := d.Set("next_value", seq.NextValue); err != nil {
		return err
	}
	if seq.Ordered {
		if err := d.Set("ordering", "ORDER"); err != nil {
			return err
		}
	} else {
		if err := d.Set("ordering", "NOORDER"); err != nil {
			return err
		}
	}

	if err := d.Set("fully_qualified_name", id.FullyQualifiedName()); err != nil {
		return err
	}
	return nil
}

func UpdateSequence(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("comment") {
		req := sdk.NewAlterSequenceRequest(id)
		req.WithSet(sdk.NewSequenceSetRequest().WithComment(sdk.String(d.Get("comment").(string))))
		if err := client.Sequences.Alter(ctx, req); err != nil {
			return err
		}
	}

	if d.HasChange("increment") {
		req := sdk.NewAlterSequenceRequest(id)
		req.WithSetIncrement(sdk.Int(d.Get("increment").(int)))
		if err := client.Sequences.Alter(ctx, req); err != nil {
			return err
		}
	}

	if d.HasChange("ordering") {
		req := sdk.NewAlterSequenceRequest(id)
		req.WithSet(sdk.NewSequenceSetRequest().WithValuesBehavior(sdk.ValuesBehaviorPointer(sdk.ValuesBehavior(d.Get("ordering").(string)))))
		if err := client.Sequences.Alter(ctx, req); err != nil {
			return err
		}
	}
	return ReadSequence(d, meta)
}

func DeleteSequence(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Sequences.Drop(ctx, sdk.NewDropSequenceRequest(id).WithIfExists(sdk.Bool(true)))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
