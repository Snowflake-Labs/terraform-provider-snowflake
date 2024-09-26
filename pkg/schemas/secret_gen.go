// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowSecretSchema represents output of SHOW query for the single Secret.
var ShowSecretSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
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
	"secret_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_scopes": {
		Type:     schema.TypeInvalid,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowSecretSchema

func SecretToSchema(secret *sdk.Secret) map[string]any {
	secretSchema := make(map[string]any)
	secretSchema["created_on"] = secret.CreatedOn.String()
	secretSchema["name"] = secret.Name
	secretSchema["schema_name"] = secret.SchemaName
	secretSchema["database_name"] = secret.DatabaseName
	secretSchema["owner"] = secret.Owner
	if secret.Comment != nil {
		secretSchema["comment"] = secret.Comment
	}
	secretSchema["secret_type"] = secret.SecretType
	secretSchema["oauth_scopes"] = secret.OauthScopes
	secretSchema["owner_role_type"] = secret.OwnerRoleType
	return secretSchema
}

var _ = SecretToSchema
