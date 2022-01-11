package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestSAMLIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.SAMLIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSAMLIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":            "test_saml_integration",
		"enabled":         true,
		"saml2_issuer":    "test_issuer",
		"saml2_sso_url":   "https://testsamlissuer.com",
		"saml2_provider":  "CUSTOM",
		"saml2_x509_cert": "MIICdummybase64certificate",
	}
	d := schema.TestResourceDataRaw(t, resources.SAMLIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE SECURITY INTEGRATION "test_saml_integration" TYPE=SAML2 SAML2_ISSUER='test_issuer' SAML2_PROVIDER='CUSTOM' SAML2_SSO_URL='https://testsamlissuer.com' SAML2_X509_CERT='MIICdummybase64certificate' ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadSAMLIntegration(mock)

		err := resources.CreateSAMLIntegration(d, db)
		r.NoError(err)
	})
}

func TestSAMLIntegrationRead(t *testing.T) {
	r := require.New(t)

	d := samlIntegration(t, "test_saml_integration", map[string]interface{}{"name": "test_saml_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadSAMLIntegration(mock)

		err := resources.ReadSAMLIntegration(d, db)
		r.NoError(err)
	})
}

func TestSAMLIntegrationDelete(t *testing.T) {
	r := require.New(t)

	d := samlIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP SECURITY INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteSAMLIntegration(d, db)
		r.NoError(err)
	})
}

func expectReadSAMLIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on"},
	).AddRow("test_saml_integration", "SAML2", "SECURITY", true, "now")
	mock.ExpectQuery(`^SHOW SECURITY INTEGRATIONS LIKE 'test_saml_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("SAML2_X509_CERT", "String", "MIICdummybase64certificate", nil).
		AddRow("SAML2_PROVIDER", "String", "CUSTOM", nil).
		AddRow("SAML2_ENABLE_SP_INITIATED", "Boolean", false, false).
		AddRow("SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", "String", "MyLabel", nil).
		AddRow("SAML2_SSO_URL", "String", "https://testsamlissuer.com", nil).
		AddRow("SAML2_ISSUER", "String", "test_issuer", nil).
		AddRow("SAML2_SNOWFLAKE_X509_CERT", "String", "MIICdummybase64certificate", nil).
		AddRow("SAML2_REQUESTED_NAMEID_FORMAT", "String", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress").
		AddRow("SAML2_FORCE_AUTHN", "Boolean", false, false).
		AddRow("SAML2_POST_LOGOUT_REDIRECT_URL", "String", "https://myredirecturl.com", nil).
		AddRow("SAML2_SIGN_REQUEST", "Boolean", false, false).
		AddRow("SAML2_SNOWFLAKE_ACS_URL", "String", "https://myinstance.my-region-1.snowflakecomputing.com/fed/login", nil).
		AddRow("SAML2_SNOWFLAKE_ISSUER_URL", "String", "https://myinstance.my-region-1.snowflakecomputing.com", nil).
		AddRow("SAML2_SNOWFLAKE_METADATA", "String", "<md:EntityDescriptor...>", nil).
		AddRow("SAML2_DIGEST_METHODS_USED", "http://www.w3.org/2001/04/xmlenc#sha256", "CUSTOM", nil).
		AddRow("SAML2_SIGNATURE_METHODS_USED", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256", "CUSTOM", nil).
		AddRow("COMMENT", "String", "Some Comment", nil)

	mock.ExpectQuery(`DESCRIBE SECURITY INTEGRATION "test_saml_integration"$`).WillReturnRows(descRows)
}
