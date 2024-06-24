package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DatabaseShowSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_transient": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_default": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_current": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"origin": {
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
	"options": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"retention_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"resource_group": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func DatabaseToSchema(database sdk.Database) map[string]any {
	return map[string]any{
		"created_on":      database.CreatedOn.String(),
		"name":            database.Name,
		"kind":            database.Kind,
		"is_transient":    database.Transient,
		"is_default":      database.IsDefault,
		"is_current":      database.IsCurrent,
		"origin":          database.Origin,
		"owner":           database.Owner,
		"comment":         database.Comment,
		"options":         database.Options,
		"retention_time":  database.RetentionTime,
		"resource_group":  database.ResourceGroup,
		"owner_role_type": database.OwnerRoleType,
	}
}

var DatabaseDescribeSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func DatabaseDescriptionToSchema(description sdk.DatabaseDetails) []map[string]any {
	result := make([]map[string]any, len(description.Rows))
	for i, row := range description.Rows {
		result[i] = map[string]any{
			"created_on": row.CreatedOn.String(),
			"name":       row.Name,
			"kind":       row.Kind,
		}
	}
	return result
}
