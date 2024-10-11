package testint

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func schemaObjectShowByIDWrapper[T any](showByIdFn func(context.Context, sdk.SchemaObjectIdentifier) (*T, error)) func(context.Context, sdk.SchemaObjectIdentifier) error {
	return func(ctx context.Context, id sdk.SchemaObjectIdentifier) error {
		_, err := showByIdFn(ctx, id)
		return err
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
	doesNotExistOrOperationCannotBePerformed := "Object does not exist, or operation cannot be performed"

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
	doesNotExistOrOperationCannotBePerformed := "Object does not exist, or operation cannot be performed"

	database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

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
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifier(database.ID().Name(), "non-existing-schema", "non-existing-schema-object"))
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
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifierWithArguments(database.ID().Name(), "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorContains(t, err, tt.ExpectedErr)
		})
	}
}
