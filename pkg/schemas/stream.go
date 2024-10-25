package schemas

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeStreamSchema represents output of SHOW query for the single Stream.
var DescribeStreamSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"table_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"source_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"base_tables": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"stale": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"mode": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"stale_after": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"invalid_reason": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowStreamSchema

// TODO(SNOW-1733130): Remove the logic for stage handling. Use schema object identifiers in the schema.
func StreamDescriptionToSchema(stream sdk.Stream) map[string]any {
	streamSchema := make(map[string]any)
	streamSchema["created_on"] = stream.CreatedOn.String()
	streamSchema["name"] = stream.Name
	streamSchema["database_name"] = stream.DatabaseName
	streamSchema["schema_name"] = stream.SchemaName
	if stream.Owner != nil {
		streamSchema["owner"] = stream.Owner
	}
	if stream.Comment != nil {
		streamSchema["comment"] = stream.Comment
	}
	if stream.TableName != nil {
		if stream.SourceType != nil && *stream.SourceType == sdk.StreamSourceTypeStage {
			streamSchema["table_name"] = *stream.TableName
		} else {
			tableId, err := sdk.ParseSchemaObjectIdentifier(*stream.TableName)
			if err != nil {
				log.Printf("[DEBUG] could not parse table ID: %v", err)
			} else {
				streamSchema["table_name"] = tableId.FullyQualifiedName()
			}
		}
	}
	if stream.SourceType != nil {
		streamSchema["source_type"] = stream.SourceType
	}
	if stream.BaseTables != nil {
		if stream.SourceType != nil && *stream.SourceType == sdk.StreamSourceTypeStage {
			streamSchema["base_tables"] = stream.BaseTables
		} else {
			streamSchema["base_tables"] = collections.Map(stream.BaseTables, func(s string) string {
				id, err := sdk.ParseSchemaObjectIdentifier(s)
				if err != nil {
					log.Printf("[DEBUG] could not parse base table ID: %v", err)
					return ""
				}
				return id.FullyQualifiedName()
			})
		}
	}
	if stream.Type != nil {
		streamSchema["type"] = stream.Type
	}
	streamSchema["stale"] = stream.Stale
	if stream.Mode != nil {
		streamSchema["mode"] = stream.Mode
	}
	if stream.StaleAfter != nil {
		streamSchema["stale_after"] = stream.StaleAfter.String()
	}
	if stream.InvalidReason != nil {
		streamSchema["invalid_reason"] = stream.InvalidReason
	}
	if stream.OwnerRoleType != nil {
		streamSchema["owner_role_type"] = stream.OwnerRoleType
	}
	return streamSchema
}

var _ = StreamToSchema
