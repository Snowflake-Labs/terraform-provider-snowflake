package datasources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var likeSchema = &schema.Schema{
	Type:        schema.TypeString,
	Optional:    true,
	Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
}

var extendedInSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	Description: "IN clause to filter the list of objects",
	MaxItems:    1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "Returns records for the entire account.",
				ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
			},
			"database": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Returns records for the current database in use or for a specified database.",
				ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			"schema": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Returns records for the current schema in use or a specified schema. Use fully qualified name.",
				ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				ValidateDiagFunc: resources.IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
			},
			"application": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Returns records for the specified application.",
				ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			"application_package": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Returns records for the specified application package.",
				ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
		},
	},
}

var startsWithSchema = &schema.Schema{
	Type:        schema.TypeString,
	Optional:    true,
	Description: "Filters the output with **case-sensitive** characters indicating the beginning of the object name.",
}

var limitFromSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	Description: "Limits the number of rows returned. If the `limit.from` is set, then the limit wll start from the first element matched by the expression. The expression is only used to match with the first element, later on the elements are not matched by the prefix, but you can enforce a certain pattern with `starts_with` or `like`.",
	MaxItems:    1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rows": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The maximum number of rows to return.",
			},
			"from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a **case-sensitive** pattern that is used to match object name. After the first match, the limit on the number of rows will be applied.",
			},
		},
	},
}

func handleLike(d *schema.ResourceData, setField **sdk.Like) {
	if likePattern, ok := d.GetOk("like"); ok {
		*setField = &sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		}
	}
}

func handleStartsWith(d *schema.ResourceData, setField **string) {
	if startsWith, ok := d.GetOk("starts_with"); ok {
		*setField = sdk.String(startsWith.(string))
	}
}

func handleLimitFrom(d *schema.ResourceData, setField **sdk.LimitFrom) {
	if v, ok := d.GetOk("limit"); ok {
		l := v.([]any)[0].(map[string]any)
		limit := &sdk.LimitFrom{}
		if v, ok := l["rows"]; ok {
			rows := v.(int)
			limit.Rows = sdk.Int(rows)
		}
		if v, ok := l["from"]; ok {
			from := v.(string)
			limit.From = sdk.String(from)
		}
		*setField = limit
	}
}

func handleExtendedIn(d *schema.ResourceData, setField **sdk.ExtendedIn) error {
	if v, ok := d.GetOk("in"); ok {
		in := v.([]any)[0].(map[string]any)
		if v, ok := in["account"]; ok && v.(bool) {
			*setField = &sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}
		}
		if v, ok := in["database"]; ok {
			if database := v.(string); database != "" {
				*setField = &sdk.ExtendedIn{In: sdk.In{Database: sdk.NewAccountObjectIdentifier(database)}}
			}
		}
		if v, ok := in["schema"]; ok {
			if schema := v.(string); schema != "" {
				schemaId, err := sdk.ParseDatabaseObjectIdentifier(schema)
				if err != nil {
					return err
				}
				*setField = &sdk.ExtendedIn{In: sdk.In{Schema: schemaId}}
			}
		}
		if v, ok := in["application"]; ok {
			if application := v.(string); application != "" {
				*setField = &sdk.ExtendedIn{Application: sdk.NewAccountObjectIdentifier(application)}
			}
		}
		if v, ok := in["application_package"]; ok {
			if applicationPackage := v.(string); applicationPackage != "" {
				*setField = &sdk.ExtendedIn{ApplicationPackage: sdk.NewAccountObjectIdentifier(applicationPackage)}
			}
		}
	}
	return nil
}
