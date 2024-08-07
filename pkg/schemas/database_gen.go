// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowDatabaseSchema represents output of SHOW query for the single Database.
var ShowDatabaseSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
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
	"dropped_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"transient": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowDatabaseSchema

func DatabaseToSchema(database *sdk.Database) map[string]any {
	databaseSchema := make(map[string]any)
	databaseSchema["created_on"] = database.CreatedOn.String()
	databaseSchema["name"] = database.Name
	databaseSchema["is_default"] = database.IsDefault
	databaseSchema["is_current"] = database.IsCurrent
	if database.Origin != nil {
		databaseSchema["origin"] = database.Origin.FullyQualifiedName()
	}
	databaseSchema["owner"] = database.Owner
	databaseSchema["comment"] = database.Comment
	databaseSchema["options"] = database.Options
	databaseSchema["retention_time"] = database.RetentionTime
	databaseSchema["resource_group"] = database.ResourceGroup
	databaseSchema["dropped_on"] = database.DroppedOn.String()
	databaseSchema["transient"] = database.Transient
	databaseSchema["kind"] = database.Kind
	databaseSchema["owner_role_type"] = database.OwnerRoleType
	return databaseSchema
}

var _ = DatabaseToSchema
