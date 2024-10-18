// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type SecretWithClientCredentialsResourceAssert struct {
	*assert.ResourceAssert
}

func SecretWithClientCredentialsResource(t *testing.T, name string) *SecretWithClientCredentialsResourceAssert {
	t.Helper()

	return &SecretWithClientCredentialsResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedSecretWithClientCredentialsResource(t *testing.T, id string) *SecretWithClientCredentialsResourceAssert {
	t.Helper()

	return &SecretWithClientCredentialsResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (s *SecretWithClientCredentialsResourceAssert) HasApiAuthenticationString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("api_authentication", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasCommentString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("comment", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasDatabaseString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("database", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasFullyQualifiedNameString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNameString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("name", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopesString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSchemaString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("schema", expected))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasSecretTypeString(expected string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("secret_type", expected))
	return s
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (s *SecretWithClientCredentialsResourceAssert) HasNoApiAuthentication() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("api_authentication"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoComment() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("comment"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoDatabase() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("database"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoFullyQualifiedName() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoName() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("name"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoOauthScopes() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("oauth_scopes"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoSchema() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("schema"))
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasNoSecretType() *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueNotSet("secret_type"))
	return s
}