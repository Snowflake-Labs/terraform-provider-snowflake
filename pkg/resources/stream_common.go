package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func handleStreamTimeTravel(d *schema.ResourceData) *sdk.OnStreamRequest {
	if v := d.Get(AtAttributeName).([]any); len(v) > 0 {
		return sdk.NewOnStreamRequest().WithAt(true).WithStatement(handleStreamTimeTravelStatement(v[0].(map[string]any)))
	}
	if v := d.Get(BeforeAttributeName).([]any); len(v) > 0 {
		return sdk.NewOnStreamRequest().WithBefore(true).WithStatement(handleStreamTimeTravelStatement(v[0].(map[string]any)))
	}
	return nil
}

func handleStreamTimeTravelStatement(timeTravelConfig map[string]any) sdk.OnStreamStatementRequest {
	statement := sdk.OnStreamStatementRequest{}
	if v := timeTravelConfig["timestamp"].(string); len(v) > 0 {
		statement.WithTimestamp(v)
	}
	if v := timeTravelConfig["offset"].(string); len(v) > 0 {
		statement.WithOffset(v)
	}
	if v := timeTravelConfig["statement"].(string); len(v) > 0 {
		statement.WithStatement(v)
	}
	if v := timeTravelConfig["stream"].(string); len(v) > 0 {
		statement.WithStream(v)
	}
	return statement
}

func DeleteStreamContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting stream",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}

var atSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	Description: externalChangesNotDetectedFieldDescription("This field specifies that the request is inclusive of any changes made by a statement or transaction with a timestamp equal to the specified parameter. Due to Snowflake limitations, the provider does not detect external changes on this field."),
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"timestamp": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies an exact date and time to use for Time Travel. The value must be explicitly cast to a TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ, or TIMESTAMP_TZ data type.",
				ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
			},
			"offset": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the difference in seconds from the current time to use for Time Travel, in the form -N where N can be an integer or arithmetic expression (e.g. -120 is 120 seconds, -30*60 is 1800 seconds or 30 minutes).",
				ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
			},
			"statement": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the query ID of a statement to use as the reference point for Time Travel. This parameter supports any statement of one of the following types: DML (e.g. INSERT, UPDATE, DELETE), TCL (BEGIN, COMMIT transaction), SELECT.",
				ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
			},
			"stream": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specifies the identifier (i.e. name) for an existing stream on the queried table or view. The current offset in the stream is used as the AT point in time for returning change data for the source object.",
				ExactlyOneOf:     []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
				DiffSuppressFunc: suppressIdentifierQuoting,
				ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			},
		},
	},
	ConflictsWith: []string{"before"},
}

var beforeSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	Description: externalChangesNotDetectedFieldDescription("This field specifies that the request refers to a point immediately preceding the specified parameter. This point in time is just before the statement, identified by its query ID, is completed.  Due to Snowflake limitations, the provider does not detect external changes on this field."),
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"timestamp": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies an exact date and time to use for Time Travel. The value must be explicitly cast to a TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ, or TIMESTAMP_TZ data type.",
				ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
			},
			"offset": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the difference in seconds from the current time to use for Time Travel, in the form -N where N can be an integer or arithmetic expression (e.g. -120 is 120 seconds, -30*60 is 1800 seconds or 30 minutes).",
				ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
			},
			"statement": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the query ID of a statement to use as the reference point for Time Travel. This parameter supports any statement of one of the following types: DML (e.g. INSERT, UPDATE, DELETE), TCL (BEGIN, COMMIT transaction), SELECT.",
				ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
			},
			"stream": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specifies the identifier (i.e. name) for an existing stream on the queried table or view. The current offset in the stream is used as the AT point in time for returning change data for the source object.",
				ExactlyOneOf:     []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
				DiffSuppressFunc: suppressIdentifierQuoting,
				ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			},
		},
	},
	ConflictsWith: []string{"at"},
}
