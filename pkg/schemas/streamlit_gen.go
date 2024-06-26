// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowStreamlitSchema represents output of SHOW query for the single Streamlit.
var ShowStreamlitSchema = map[string]*schema.Schema{
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
	"title": {
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
	"query_warehouse": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"url_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowStreamlitSchema

func StreamlitToSchema(streamlit *sdk.Streamlit) map[string]any {
	streamlitSchema := make(map[string]any)
	streamlitSchema["created_on"] = streamlit.CreatedOn
	streamlitSchema["name"] = streamlit.Name
	streamlitSchema["database_name"] = streamlit.DatabaseName
	streamlitSchema["schema_name"] = streamlit.SchemaName
	streamlitSchema["title"] = streamlit.Title
	streamlitSchema["owner"] = streamlit.Owner
	streamlitSchema["comment"] = streamlit.Comment
	streamlitSchema["query_warehouse"] = streamlit.QueryWarehouse
	streamlitSchema["url_id"] = streamlit.UrlId
	streamlitSchema["owner_role_type"] = streamlit.OwnerRoleType
	return streamlitSchema
}

var _ = StreamlitToSchema
