// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowStorageIntegrationSchema represents output of SHOW query for the single StorageIntegration.
var ShowStorageIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"category": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowStorageIntegrationSchema

func StorageIntegrationToSchema(storageIntegration *sdk.StorageIntegration) map[string]any {
	storageIntegrationSchema := make(map[string]any)
	storageIntegrationSchema["name"] = storageIntegration.Name
	storageIntegrationSchema["storage_type"] = storageIntegration.StorageType
	storageIntegrationSchema["category"] = storageIntegration.Category
	storageIntegrationSchema["enabled"] = storageIntegration.Enabled
	storageIntegrationSchema["comment"] = storageIntegration.Comment
	storageIntegrationSchema["created_on"] = storageIntegration.CreatedOn.String()
	return storageIntegrationSchema
}

var _ = StorageIntegrationToSchema
