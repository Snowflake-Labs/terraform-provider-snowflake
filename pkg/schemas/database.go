package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowDatabaseSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_transient": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func DatabaseToSchema(database *sdk.Database) map[string]any {
	databaseSchema := make(map[string]any)
	databaseSchema["name"] = database.Name
	databaseSchema["is_transient"] = database.Transient
	databaseSchema["comment"] = database.Comment
	return databaseSchema
}
