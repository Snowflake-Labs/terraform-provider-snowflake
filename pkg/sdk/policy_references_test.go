package sdk

import (
	"strings"
	"testing"
	// "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
)

func TestPolicyReferencesGetForEntity(t *testing.T) {
	userName := NewAccountObjectIdentifierFromFullyQualifiedName("USER")

	t.Run("validation: missing refEntityName", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			tableFunction: &tableFunction{
				table: Bool(true),
				policyReferenceFunction: &policyReferenceFunction{
					functionFullyQualifiedName: Bool(true),
					arguments: &policyReferenceFunctionArguments{
						refEntityName:   nil,
						refEntityDomain: String("user"),
					},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("getForEntityPolicyReferenceOptions", "refEntityName"))
	})

	t.Run("validation: missing refEntityDomain", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			tableFunction: &tableFunction{
				table: Bool(true),
				policyReferenceFunction: &policyReferenceFunction{
					functionFullyQualifiedName: Bool(true),
					arguments: &policyReferenceFunctionArguments{
						refEntityName:   []ObjectIdentifier{userName},
						refEntityDomain: nil,
					},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("getForEntityPolicyReferenceOptions", "refEntityDomain"))
	})

	t.Run("validation: domain: user", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			tableFunction: &tableFunction{
				table: Bool(true),
				policyReferenceFunction: &policyReferenceFunction{
					functionFullyQualifiedName: Bool(true),
					arguments: &policyReferenceFunctionArguments{
						refEntityName:   []ObjectIdentifier{userName},
						refEntityDomain: String("user"),
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (ref_entity_name => '%s', ref_entity_domain => 'user'))", strings.ReplaceAll(userName.FullyQualifiedName(), `"`, `\"`))
	})

	tableName := NewSchemaObjectIdentifier("db", "schema", "table")
	t.Run("validation: domain: table", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			tableFunction: &tableFunction{
				table: Bool(true),
				policyReferenceFunction: &policyReferenceFunction{
					functionFullyQualifiedName: Bool(true),
					arguments: &policyReferenceFunctionArguments{
						refEntityName:   []ObjectIdentifier{tableName},
						refEntityDomain: String("table"),
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (ref_entity_name => '%s', ref_entity_domain => 'table'))", strings.ReplaceAll(tableName.FullyQualifiedName(), `"`, `\"`))
	})

	accountName := NewAccountObjectIdentifier("account")
	t.Run("validation: domain: table", func(t *testing.T) {
		opts := &getForEntityPolicyReferenceOptions{
			tableFunction: &tableFunction{
				table: Bool(true),
				policyReferenceFunction: &policyReferenceFunction{
					functionFullyQualifiedName: Bool(true),
					arguments: &policyReferenceFunctionArguments{
						refEntityName:   []ObjectIdentifier{accountName},
						refEntityDomain: String("account"),
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES (ref_entity_name => '%s', ref_entity_domain => 'account'))", strings.ReplaceAll(accountName.FullyQualifiedName(), `"`, `\"`))
	})
}
