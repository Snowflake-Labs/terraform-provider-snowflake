package sdk

import "testing"

func TestDataMetricFunctionReferences_GetForEntity(t *testing.T) {
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *GetForEntityDataMetricFunctionReferenceOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: missing parameters", func(t *testing.T) {
		opts := &GetForEntityDataMetricFunctionReferenceOptions{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GetForEntityDataMetricFunctionReferenceOptions", "parameters"))
	})

	t.Run("validation: missing arguments", func(t *testing.T) {
		opts := &GetForEntityDataMetricFunctionReferenceOptions{
			parameters: &dataMetricFunctionReferenceParameters{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("dataMetricFunctionReferenceParameters", "arguments"))
	})

	t.Run("validation: missing refEntityName", func(t *testing.T) {
		opts := &GetForEntityDataMetricFunctionReferenceOptions{
			parameters: &dataMetricFunctionReferenceParameters{
				arguments: &dataMetricFunctionReferenceFunctionArguments{
					refEntityDomain: Pointer(DataMetricFuncionRefEntityDomainView),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("dataMetricFunctionReferenceFunctionArguments", "refEntityName"))
	})

	t.Run("validation: missing refEntityDomain", func(t *testing.T) {
		opts := &GetForEntityDataMetricFunctionReferenceOptions{
			parameters: &dataMetricFunctionReferenceParameters{
				arguments: &dataMetricFunctionReferenceFunctionArguments{
					refEntityName: []ObjectIdentifier{NewSchemaObjectIdentifier("a", "b", "c")},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("dataMetricFunctionReferenceFunctionArguments", "refEntityDomain"))
	})

	t.Run("view domain", func(t *testing.T) {
		opts := &GetForEntityDataMetricFunctionReferenceOptions{
			parameters: &dataMetricFunctionReferenceParameters{
				arguments: &dataMetricFunctionReferenceFunctionArguments{
					refEntityName:   []ObjectIdentifier{NewSchemaObjectIdentifier("a", "b", "c")},
					refEntityDomain: Pointer(DataMetricFuncionRefEntityDomainView),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SELECT * FROM TABLE (SNOWFLAKE.INFORMATION_SCHEMA.DATA_METRIC_FUNCTION_REFERENCES (REF_ENTITY_NAME => '\"a\".\"b\".\"c\"', REF_ENTITY_DOMAIN => 'VIEW'))`)
	})
}
