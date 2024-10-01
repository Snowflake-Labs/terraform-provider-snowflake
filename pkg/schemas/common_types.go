package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ParameterListSchema represents Snowflake parameter object.
var ParameterListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: ShowParameterSchema,
	},
}

// DescribePropertyListSchema represents Snowflake property object returned by DESCRIBE query.
var DescribePropertyListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: ShowSecurityIntegrationPropertySchema,
	},
}

var FullyQualifiedNameSchema = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).",
}

var AtSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	Description: "This field specifies that the request is inclusive of any changes made by a statement or transaction with a timestamp equal to the specified parameter. Due to Snowflake limitations, the provider does not detect external changes on this field.",
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the identifier (i.e. name) for an existing stream on the queried table or view. The current offset in the stream is used as the AT point in time for returning change data for the source object.",
				ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
			},
		},
	},
	ConflictsWith: []string{"before"},
}

var BeforeSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	Description: "This field specifies that the request refers to a point immediately preceding the specified parameter. This point in time is just before the statement, identified by its query ID, is completed.  Due to Snowflake limitations, the provider does not detect external changes on this field.",
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the identifier (i.e. name) for an existing stream on the queried table or view. The current offset in the stream is used as the AT point in time for returning change data for the source object.",
				ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
			},
		},
	},
	ConflictsWith: []string{"at"},
}
