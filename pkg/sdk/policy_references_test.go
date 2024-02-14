package sdk

import (
	"testing"
)

func TestPolicyReferencesGetForEntity(t *testing.T) {
	t.Run("validation: missing parameters", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("getForEntityPolicyReferenceOptions", "parameters"))
	})

	t.Run("validation: missing arguments", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("policyReferenceParameters", "arguments"))
	})

	t.Run("validation: missing refEntityName", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityDomain: Pointer(PolicyEntityDomainUser),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("policyReferenceFunctionArguments", "refEntityName"))
	})

	t.Run("validation: missing refEntityDomain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName: []ObjectIdentifier{NewAccountObjectIdentifierFromFullyQualifiedName("user_name")},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("policyReferenceFunctionArguments", "refEntityDomain"))
	})

	t.Run("user domain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewAccountObjectIdentifier("user_name")},
					refEntityDomain: Pointer(PolicyEntityDomainUser),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"user_name\"', REF_ENTITY_DOMAIN => 'USER'))`)
	})

	t.Run("table domain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewSchemaObjectIdentifier("db", "schema", "table")},
					refEntityDomain: Pointer(PolicyEntityDomainTable),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"db\".\"schema\".\"table\"', REF_ENTITY_DOMAIN => 'TABLE'))`)
	})

	t.Run("account domain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewAccountObjectIdentifier("account_name")},
					refEntityDomain: Pointer(PolicyEntityDomainAccount),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"account_name\"', REF_ENTITY_DOMAIN => 'ACCOUNT'))`)
	})

	t.Run("integration domain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewAccountObjectIdentifier("integration_name")},
					refEntityDomain: Pointer(PolicyEntityDomainIntegration),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"integration_name\"', REF_ENTITY_DOMAIN => 'INTEGRATION'))`)
	})

	t.Run("tag domain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewSchemaObjectIdentifier("db", "schema", "tag_name")},
					refEntityDomain: Pointer(PolicyEntityDomainTag),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"db\".\"schema\".\"tag_name\"', REF_ENTITY_DOMAIN => 'TAG'))`)
	})

	t.Run("view domain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			parameters: &policyReferenceParameters{
				arguments: &policyReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewSchemaObjectIdentifier("db", "schema", "view_name")},
					refEntityDomain: Pointer(PolicyEntityDomainView),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (REF_ENTITY_NAME => '\"db\".\"schema\".\"view_name\"', REF_ENTITY_DOMAIN => 'VIEW'))`)
	})
}
