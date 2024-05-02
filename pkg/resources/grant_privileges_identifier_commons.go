package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type OnSchemaGrantKind string

const (
	OnSchemaSchemaGrantKind                  OnSchemaGrantKind = "OnSchema"
	OnAllSchemasInDatabaseSchemaGrantKind    OnSchemaGrantKind = "OnAllSchemasInDatabase"
	OnFutureSchemasInDatabaseSchemaGrantKind OnSchemaGrantKind = "OnFutureSchemasInDatabase"
)

type OnSchemaObjectGrantKind string

const (
	OnObjectSchemaObjectGrantKind OnSchemaObjectGrantKind = "OnObject"
	OnAllSchemaObjectGrantKind    OnSchemaObjectGrantKind = "OnAll"
	OnFutureSchemaObjectGrantKind OnSchemaObjectGrantKind = "OnFuture"
)

type OnSchemaGrantData struct {
	Kind         OnSchemaGrantKind
	SchemaName   *sdk.DatabaseObjectIdentifier
	DatabaseName *sdk.AccountObjectIdentifier
}

func (d *OnSchemaGrantData) String() string {
	var parts []string
	parts = append(parts, string(d.Kind))
	switch d.Kind {
	case OnSchemaSchemaGrantKind:
		parts = append(parts, d.SchemaName.FullyQualifiedName())
	case OnAllSchemasInDatabaseSchemaGrantKind, OnFutureSchemasInDatabaseSchemaGrantKind:
		parts = append(parts, d.DatabaseName.FullyQualifiedName())
	}
	return strings.Join(parts, helpers.IDDelimiter)
}

type OnSchemaObjectGrantData struct {
	Kind          OnSchemaObjectGrantKind
	Object        *sdk.Object
	OnAllOrFuture *BulkOperationGrantData
}

func (d *OnSchemaObjectGrantData) String() string {
	var parts []string
	parts = append(parts, string(d.Kind))
	switch d.Kind {
	case OnObjectSchemaObjectGrantKind:
		parts = append(parts, fmt.Sprintf("%s|%s", d.Object.ObjectType, d.Object.Name.FullyQualifiedName()))
	case OnAllSchemaObjectGrantKind, OnFutureSchemaObjectGrantKind:
		parts = append(parts, d.OnAllOrFuture.String())
	}
	return strings.Join(parts, helpers.IDDelimiter)
}

type BulkOperationGrantKind string

const (
	InDatabaseBulkOperationGrantKind BulkOperationGrantKind = "InDatabase"
	InSchemaBulkOperationGrantKind   BulkOperationGrantKind = "InSchema"
)

type BulkOperationGrantData struct {
	ObjectNamePlural sdk.PluralObjectType
	Kind             BulkOperationGrantKind
	Database         *sdk.AccountObjectIdentifier
	Schema           *sdk.DatabaseObjectIdentifier
}

func (d *BulkOperationGrantData) String() string {
	var parts []string
	parts = append(parts, d.ObjectNamePlural.String())
	parts = append(parts, string(d.Kind))
	switch d.Kind {
	case InDatabaseBulkOperationGrantKind:
		parts = append(parts, d.Database.FullyQualifiedName())
	case InSchemaBulkOperationGrantKind:
		parts = append(parts, d.Schema.FullyQualifiedName())
	}
	return strings.Join(parts, helpers.IDDelimiter)
}

func getBulkOperationGrantData(in *sdk.GrantOnSchemaObjectIn) *BulkOperationGrantData {
	bulkOperationGrantData := &BulkOperationGrantData{
		ObjectNamePlural: in.PluralObjectType,
	}

	if in.InDatabase != nil {
		bulkOperationGrantData.Kind = InDatabaseBulkOperationGrantKind
		bulkOperationGrantData.Database = in.InDatabase
	}

	if in.InSchema != nil {
		bulkOperationGrantData.Kind = InSchemaBulkOperationGrantKind
		bulkOperationGrantData.Schema = in.InSchema
	}

	return bulkOperationGrantData
}

func getGrantOnSchemaObjectIn(allOrFuture map[string]any) *sdk.GrantOnSchemaObjectIn {
	grantOnSchemaObjectIn := &sdk.GrantOnSchemaObjectIn{
		PluralObjectType: sdk.PluralObjectType(strings.ToUpper(allOrFuture["object_type_plural"].(string))),
	}

	if inDatabase, ok := allOrFuture["in_database"].(string); ok && len(inDatabase) > 0 {
		grantOnSchemaObjectIn.InDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(inDatabase))
	}

	if inSchema, ok := allOrFuture["in_schema"].(string); ok && len(inSchema) > 0 {
		grantOnSchemaObjectIn.InSchema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(inSchema))
	}

	return grantOnSchemaObjectIn
}
