//go:build !account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func schemaObjectShowByIDWrapper[T any](showByIdFn func(context.Context, sdk.SchemaObjectIdentifier) (*T, error)) func(context.Context, sdk.SchemaObjectIdentifier) error {
	return func(ctx context.Context, id sdk.SchemaObjectIdentifier) error {
		_, err := showByIdFn(ctx, id)
		return err
	}
}

func schemaObjectDropWrapper[T any](dropFn func(ctx context.Context, req T) error, ctx context.Context, req T) func() error {
	return func() error {
		return dropFn(ctx, req)
	}
}

func schemaObjectWithArgumentsShowByIDWrapper[T any](showByIdFn func(context.Context, sdk.SchemaObjectIdentifierWithArguments) (*T, error)) func(context.Context, sdk.SchemaObjectIdentifierWithArguments) error {
	return func(ctx context.Context, id sdk.SchemaObjectIdentifierWithArguments) error {
		_, err := showByIdFn(ctx, id)
		return err
	}
}

func TestInt_ShowSchemaObjectInNonExistingDatabase(t *testing.T) {
	doesNotExistOrNotAuthorized := sdk.ErrObjectNotExistOrAuthorized.Error() // Database '\"non-existing-database\"' does not exist or not authorized
	doesNotExistOrOperationCannotBePerformed := sdk.ErrDoesNotExistOrOperationCannotBePerformed.Error()

	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifier) error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tables.ShowByID)},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).DynamicTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).CortexSearchServices.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).ExternalTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).EventTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Views.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaterializedViews.ShowByID)},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Sequences.ShowByID)},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streams.ShowByID)},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tasks.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaskingPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).RowAccessPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tags.ShowByID)},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Secrets.ShowByID)},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Stages.ShowByID)},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).FileFormats.ShowByID)},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Pipes.ShowByID)},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Alerts.ShowByID)},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streamlits.ShowByID)},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).NetworkRules.ShowByID)},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).AuthenticationPolicies.ShowByID)},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifier("non-existing-database", "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifierWithArguments) error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Functions.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalFunction, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).ExternalFunctions.ShowByID)},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Procedures.ShowByID)},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifierWithArguments("non-existing-database", "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}
}

func TestInt_ShowSchemaObjectInNonExistingSchema(t *testing.T) {
	doesNotExistOrNotAuthorized := sdk.ErrObjectNotExistOrAuthorized.Error() // Schema '\"non-existing-schema\"' does not exist or not authorized
	doesNotExistOrOperationCannotBePerformed := sdk.ErrDoesNotExistOrOperationCannotBePerformed.Error()

	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifier) error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tables.ShowByID)},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).DynamicTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).CortexSearchServices.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).ExternalTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).EventTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Views.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaterializedViews.ShowByID)},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Sequences.ShowByID)},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streams.ShowByID)},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tasks.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaskingPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).RowAccessPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tags.ShowByID)},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Secrets.ShowByID)},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Stages.ShowByID)},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).FileFormats.ShowByID)},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Pipes.ShowByID)},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Alerts.ShowByID)},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streamlits.ShowByID)},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: doesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).NetworkRules.ShowByID)},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).AuthenticationPolicies.ShowByID)},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifier(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifierWithArguments) error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Functions.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalFunction, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).ExternalFunctions.ShowByID)},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: doesNotExistOrNotAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Procedures.ShowByID)},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifierWithArguments(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}
}

func TestInt_DropSchemaObjectInNonExistingDatabase(t *testing.T) {
	doesNotExistOrNotAuthorized := sdk.ErrObjectNotExistOrAuthorized.Error() // Database '\"non-existing-database\"' does not exist or not authorized

	id := sdk.NewSchemaObjectIdentifier("non-existing-database", "non-existing-schema", "non-existing-schema-object")
	idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments("non-existing-database", "non-existing-schema", "non-existing-schema-object", sdk.DataTypeInt)
	ctx := context.Background()

	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		DropFn      func() error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tables.Drop, ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).DynamicTables.Drop, ctx, sdk.NewDropDynamicTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).CortexSearchServices.Drop, ctx, sdk.NewDropCortexSearchServiceRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).ExternalTables.Drop, ctx, sdk.NewDropExternalTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).EventTables.Drop, ctx, sdk.NewDropEventTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Views.Drop, ctx, sdk.NewDropViewRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).MaterializedViews.Drop, ctx, sdk.NewDropMaterializedViewRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Sequences.Drop, ctx, sdk.NewDropSequenceRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streams.Drop, ctx, sdk.NewDropStreamRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tasks.Drop, ctx, sdk.NewDropTaskRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).MaskingPolicies.Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).RowAccessPolicies.Drop, ctx, sdk.NewDropRowAccessPolicyRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tags.Drop, ctx, sdk.NewDropTagRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Secrets.Drop, ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Stages.Drop, ctx, sdk.NewDropStageRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).FileFormats.Drop(ctx, id, &sdk.DropFileFormatOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).Pipes.Drop(ctx, id, &sdk.DropPipeOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).Alerts.Drop(ctx, id, &sdk.DropAlertOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streamlits.Drop, ctx, sdk.NewDropStreamlitRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).NetworkRules.Drop, ctx, sdk.NewDropNetworkRuleRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).AuthenticationPolicies.Drop, ctx, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(true))},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := tt.DropFn()
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		DropFn      func() error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Functions.Drop, ctx, sdk.NewDropFunctionRequest(idWithArguments).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Procedures.Drop, ctx, sdk.NewDropProcedureRequest(idWithArguments).WithIfExists(true))},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := tt.DropFn()
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}
}

func TestInt_DropSchemaObjectInNonExistingSchema(t *testing.T) {
	doesNotExistOrNotAuthorized := sdk.ErrObjectNotExistOrAuthorized.Error() // Schema '\"non-existing-schema\"' does not exist or not authorized

	ctx := context.Background()
	id := sdk.NewSchemaObjectIdentifier(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object")
	idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object", sdk.DataTypeInt)

	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		DropFn      func() error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tables.Drop, ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).DynamicTables.Drop, ctx, sdk.NewDropDynamicTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).CortexSearchServices.Drop, ctx, sdk.NewDropCortexSearchServiceRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).ExternalTables.Drop, ctx, sdk.NewDropExternalTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).EventTables.Drop, ctx, sdk.NewDropEventTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Views.Drop, ctx, sdk.NewDropViewRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).MaterializedViews.Drop, ctx, sdk.NewDropMaterializedViewRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Sequences.Drop, ctx, sdk.NewDropSequenceRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streams.Drop, ctx, sdk.NewDropStreamRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tasks.Drop, ctx, sdk.NewDropTaskRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).MaskingPolicies.Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).RowAccessPolicies.Drop, ctx, sdk.NewDropRowAccessPolicyRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tags.Drop, ctx, sdk.NewDropTagRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Secrets.Drop, ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Stages.Drop, ctx, sdk.NewDropStageRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).FileFormats.Drop(ctx, id, &sdk.DropFileFormatOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).Pipes.Drop(ctx, id, &sdk.DropPipeOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: func() error {
			return testClient(t).Alerts.Drop(ctx, id, &sdk.DropAlertOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streamlits.Drop, ctx, sdk.NewDropStreamlitRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).NetworkRules.Drop, ctx, sdk.NewDropNetworkRuleRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).AuthenticationPolicies.Drop, ctx, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(true))},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := tt.DropFn()
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr string
		DropFn      func() error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Functions.Drop, ctx, sdk.NewDropFunctionRequest(idWithArguments).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: doesNotExistOrNotAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Procedures.Drop, ctx, sdk.NewDropProcedureRequest(idWithArguments).WithIfExists(true))},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := tt.DropFn()
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}
}
