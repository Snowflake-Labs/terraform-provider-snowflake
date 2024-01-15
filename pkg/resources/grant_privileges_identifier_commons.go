package resources

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strings"
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
		parts = append(parts, d.OnAllOrFuture.ObjectNamePlural.String())
		parts = append(parts, string(d.OnAllOrFuture.Kind))
		switch d.OnAllOrFuture.Kind {
		case InDatabaseBulkOperationGrantKind:
			parts = append(parts, d.OnAllOrFuture.Database.FullyQualifiedName())
		case InSchemaBulkOperationGrantKind:
			parts = append(parts, d.OnAllOrFuture.Schema.FullyQualifiedName())
		}
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
	pluralObjectType := sdk.PluralObjectType(allOrFuture["object_type_plural"].(string))
	grantOnSchemaObjectIn := &sdk.GrantOnSchemaObjectIn{
		PluralObjectType: pluralObjectType,
	}

	if inDatabase, ok := allOrFuture["in_database"].(string); ok && len(inDatabase) > 0 {
		grantOnSchemaObjectIn.InDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(inDatabase))
	}

	if inSchema, ok := allOrFuture["in_schema"].(string); ok && len(inSchema) > 0 {
		grantOnSchemaObjectIn.InSchema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(inSchema))
	}

	return grantOnSchemaObjectIn
}
