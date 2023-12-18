package resources

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strings"
)

// TODO: Add unit tests for marshaling / unmarshalling

type DatabaseRoleGrantKind string

const (
	OnDatabaseDatabaseRoleGrantKind     DatabaseRoleGrantKind = "OnDatabase"
	OnSchemaDatabaseRoleGrantKind       DatabaseRoleGrantKind = "OnSchema"
	OnSchemaObjectDatabaseRoleGrantKind DatabaseRoleGrantKind = "OnSchemaObject"
)

// TODO: Move to the shareable file between this and grant_priv_to_role.go file
type OnSchemaGrantKind string

const (
	OnSchemaSchemaGrantKind                  OnSchemaGrantKind = "OnSchema"
	OnAllSchemasInDatabaseSchemaGrantKind    OnSchemaGrantKind = "OnAllSchemasInDatabase"
	OnFutureSchemasInDatabaseSchemaGrantKind OnSchemaGrantKind = "OnFutureSchemasInDatabase"
)

// TODO: Move to the shareable file between this and grant_priv_to_role.go file
type OnSchemaObjectGrantKind string

const (
	OnObjectSchemaObjectGrantKind OnSchemaObjectGrantKind = "OnObject"
	OnAllSchemaObjectGrantKind    OnSchemaObjectGrantKind = "OnAll"
	OnFutureSchemaObjectGrantKind OnSchemaObjectGrantKind = "OnFuture"
)

type GrantPrivilegesToDatabaseRoleId struct {
	DatabaseRoleName sdk.AccountObjectIdentifier
	WithGrantOption  bool
	Privileges       []string
	Kind             DatabaseRoleGrantKind
	Data             any
}

type OnDatabaseGrantData struct {
	DatabaseName sdk.AccountObjectIdentifier
}

type OnSchemaGrantData struct {
	Kind         OnSchemaGrantKind
	SchemaName   *sdk.DatabaseObjectIdentifier
	DatabaseName *sdk.AccountObjectIdentifier
}

type OnSchemaObjectGrantData struct {
	Kind          OnSchemaObjectGrantKind
	Object        *sdk.Object
	OnAllOrFuture *BulkOperationGrantData
}

type BulkOperationGrantKind string

const (
	InDatabaseBulkOperationGrantKind BulkOperationGrantKind = "InDatabase"
	InSchemaBulkOperationGrantKind   BulkOperationGrantKind = "InSchema"
)

type BulkOperationGrantData struct {
	ObjectNamePlural sdk.PluralObjectType
	Kind             *BulkOperationGrantKind
	Database         *sdk.AccountObjectIdentifier
	Schema           *sdk.DatabaseObjectIdentifier
}

// TODO: Describe how to put a right identifier in the documentation (so the users will be able to use it in the import)
func ParseGrantPrivilegesToDatabaseRoleId(id string) (GrantPrivilegesToDatabaseRoleId, error) {
	var databaseRoleId GrantPrivilegesToDatabaseRoleId

	parts := strings.Split(id, helpers.IDDelimiter)
	if len(parts) < 5 {
		return databaseRoleId, sdk.NewError(`database role identifier should hold at least 4 parts "<database_role_name>|<with_grant_option>|<privileges>|<grant_type>|<grant_data>"`)
	}

	databaseRoleId.DatabaseRoleName = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
	databaseRoleId.WithGrantOption = parts[1] == "true"
	privileges := strings.Split(parts[2], ",")
	if len(privileges) == 1 && privileges[0] == "" {
		privileges = []string{}
	}
	// TODO: All privileges
	databaseRoleId.Privileges = privileges
	databaseRoleId.Kind = DatabaseRoleGrantKind(parts[3])

	switch databaseRoleId.Kind {
	case OnDatabaseDatabaseRoleGrantKind:
		databaseRoleId.Data = OnDatabaseGrantData{
			DatabaseName: sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[4]),
		}
	case OnSchemaDatabaseRoleGrantKind:
		if len(parts) < 6 {
			return databaseRoleId, sdk.NewError(`database role identifier should hold at least 6 parts "<database_role_name>|<with_grant_option>|<privileges>|<grant_type>|<grant_on_schema_type>|<on_schema_grant_data>..."`)
		}
		onSchemaGrantData := OnSchemaGrantData{
			Kind: OnSchemaGrantKind(parts[4]),
		}
		switch onSchemaGrantData.Kind {
		case OnSchemaSchemaGrantKind:
			onSchemaGrantData.SchemaName = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[5]))
		case OnAllSchemasInDatabaseSchemaGrantKind, OnFutureSchemasInDatabaseSchemaGrantKind:
			onSchemaGrantData.DatabaseName = sdk.Pointer(sdk.NewAccountObjectIdentifier(parts[5]))
		default:
			return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaGrantKind: %s", onSchemaGrantData.Kind))
		}
		databaseRoleId.Data = onSchemaGrantData
	case OnSchemaObjectDatabaseRoleGrantKind:
		if len(parts) < 6 {
			return databaseRoleId, sdk.NewError(`database role identifier should hold at least 6 parts "<database_role_name>|<with_grant_option>|<privileges>|<grant_type>|<grant_on_schema_object_type>|<on_schema_object_grant_data>..."`)
		}
		onSchemaObjectGrantData := OnSchemaObjectGrantData{
			Kind: OnSchemaObjectGrantKind(parts[4]),
		}
		switch onSchemaObjectGrantData.Kind {
		case OnObjectSchemaObjectGrantKind:
			if len(parts) != 7 {
				return databaseRoleId, sdk.NewError(`database role identifier should hold 7 parts "<database_role_name>|<with_grant_option>|<privileges>|OnSchemaObject|OnObject|<object_type>|<object_name>"`)
			}
			onSchemaObjectGrantData.Object = &sdk.Object{
				ObjectType: sdk.ObjectType(parts[5]),
				Name:       sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(parts[6]),
			}
		case OnAllSchemaObjectGrantKind, OnFutureSchemaObjectGrantKind:
			bulkOperationGrantData := &BulkOperationGrantData{
				ObjectNamePlural: sdk.PluralObjectType(parts[5]),
			}
			if len(parts) > 6 {
				if len(parts) != 8 {
					return databaseRoleId, sdk.NewError(`database role identifier should hold 8 parts "<database_role_name>|<with_grant_option>|<privileges>|OnSchemaObject|On[All or Future]|<object_type_plural>|In[Database or Schema]|<identifier>"`)
				}
				bulkOperationGrantData.Kind = sdk.Pointer(BulkOperationGrantKind(parts[6]))
				switch *bulkOperationGrantData.Kind {
				case InDatabaseBulkOperationGrantKind:
					bulkOperationGrantData.Database = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[7]))
				case InSchemaBulkOperationGrantKind:
					bulkOperationGrantData.Schema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[7]))
				default:
					return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid BulkOperationGrantKind: %s", *bulkOperationGrantData.Kind))
				}
			}
			onSchemaObjectGrantData.OnAllOrFuture = bulkOperationGrantData
		default:
			return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaObjectGrantKind: %s", onSchemaObjectGrantData.Kind))
		}
		databaseRoleId.Data = onSchemaObjectGrantData
	default:
		return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid DatabaseRoleGrantKind: %s", databaseRoleId.Kind))
	}

	return databaseRoleId, nil
}
