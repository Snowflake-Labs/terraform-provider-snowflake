package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeOpenflowConnectorSchema represents output of DESCRIBE query for the single OpenflowConnector.
// DESCRIBE returns the SHOW columns minus created_on, updated_on and the database and schema names, plus the
// git commit hashes and the last version fields.
var DescribeOpenflowConnectorSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"runtime": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"connector_definition": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"display_name": {
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
	"default_version": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_version_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_version_alias": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_version_location_uri": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_version_source_location_uri": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_version_git_commit_hash": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_version_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_version_alias": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_version_location_uri": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_version_source_location_uri": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_version_git_commit_hash": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"live_version_location_uri": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"connector_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = DescribeOpenflowConnectorSchema

func OpenflowConnectorDetailsToSchema(connector sdk.OpenflowConnectorDetails) map[string]any {
	connectorSchema := make(map[string]any)
	connectorSchema["name"] = connector.Name
	connectorSchema["status"] = string(connector.Status)
	connectorSchema["runtime"] = connector.Runtime
	if connector.ConnectorDefinition != nil {
		connectorSchema["connector_definition"] = connector.ConnectorDefinition
	}
	if connector.DisplayName != nil {
		connectorSchema["display_name"] = connector.DisplayName
	}
	connectorSchema["owner"] = connector.Owner
	if connector.Comment != nil {
		connectorSchema["comment"] = connector.Comment
	}
	if connector.DefaultVersion != nil {
		connectorSchema["default_version"] = connector.DefaultVersion
	}
	if connector.DefaultVersionName != nil {
		connectorSchema["default_version_name"] = connector.DefaultVersionName
	}
	if connector.DefaultVersionAlias != nil {
		connectorSchema["default_version_alias"] = connector.DefaultVersionAlias
	}
	if connector.DefaultVersionLocationUri != nil {
		connectorSchema["default_version_location_uri"] = connector.DefaultVersionLocationUri
	}
	if connector.DefaultVersionSourceLocationUri != nil {
		connectorSchema["default_version_source_location_uri"] = connector.DefaultVersionSourceLocationUri
	}
	if connector.DefaultVersionGitCommitHash != nil {
		connectorSchema["default_version_git_commit_hash"] = connector.DefaultVersionGitCommitHash
	}
	if connector.LastVersionName != nil {
		connectorSchema["last_version_name"] = connector.LastVersionName
	}
	if connector.LastVersionAlias != nil {
		connectorSchema["last_version_alias"] = connector.LastVersionAlias
	}
	if connector.LastVersionLocationUri != nil {
		connectorSchema["last_version_location_uri"] = connector.LastVersionLocationUri
	}
	if connector.LastVersionSourceLocationUri != nil {
		connectorSchema["last_version_source_location_uri"] = connector.LastVersionSourceLocationUri
	}
	if connector.LastVersionGitCommitHash != nil {
		connectorSchema["last_version_git_commit_hash"] = connector.LastVersionGitCommitHash
	}
	if connector.LiveVersionLocationUri != nil {
		connectorSchema["live_version_location_uri"] = connector.LiveVersionLocationUri
	}
	if connector.ConnectorUrl != nil {
		connectorSchema["connector_url"] = connector.ConnectorUrl
	}
	return connectorSchema
}
