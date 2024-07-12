package resources

import (
	"context"
)

func v092ExternalOauthIntegrationStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	type renameField struct {
		from string
		to   string
	}
	fieldsToRename := []renameField{
		{from: "type", to: "external_oauth_type"},
		{from: "issuer", to: "external_oauth_issuer"},
		{from: "token_user_mapping_claims", to: "external_oauth_token_user_mapping_claim"},
		{from: "snowflake_user_mapping_attribute", to: "external_oauth_snowflake_user_mapping_attribute"},
		{from: "scope_mapping_attribute", to: "external_oauth_scope_mapping_attribute"},
		{from: "jws_keys_urls", to: "external_oauth_jws_keys_url"},
		{from: "rsa_public_key", to: "external_oauth_rsa_public_key"},
		{from: "rsa_public_key_2", to: "external_oauth_rsa_public_key_2"},
		{from: "blocked_roles", to: "external_oauth_blocked_roles_list"},
		{from: "allowed_roles", to: "external_oauth_allowed_roles_list"},
		{from: "audience_urls", to: "external_oauth_audience_list"},
		{from: "any_role_mode", to: "external_oauth_any_role_mode"},
		{from: "scope_delimiter", to: "external_oauth_scope_delimiter"},
	}

	for _, field := range fieldsToRename {
		rawState[field.to] = rawState[field.from]
		delete(rawState, field.from)
	}

	return rawState, nil
}
