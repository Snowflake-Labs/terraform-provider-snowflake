package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

/**
 * Will be removed while adding security integrations to the SDK.
 */

func samlIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.SAMLIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}

func oauthIntegration(t *testing.T, id string, params map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, resources.OAuthIntegration().Schema, params)
	r.NotNil(d)
	d.SetId(id)
	return d
}
