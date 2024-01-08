package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type DatabaseRoleGrantKind string

const (
	OnDatabaseDatabaseRoleGrantKind     DatabaseRoleGrantKind = "OnDatabase"
	OnSchemaDatabaseRoleGrantKind       DatabaseRoleGrantKind = "OnSchema"
	OnSchemaObjectDatabaseRoleGrantKind DatabaseRoleGrantKind = "OnSchemaObject"
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

type GrantPrivilegesToDatabaseRoleId struct {
	DatabaseRoleName sdk.DatabaseObjectIdentifier
	WithGrantOption  bool
	AlwaysApply      bool
	AllPrivileges    bool
	Privileges       []string
	Kind             DatabaseRoleGrantKind
	Data             fmt.Stringer
}

func (g *GrantPrivilegesToDatabaseRoleId) String() string {
	var parts []string
	parts = append(parts, g.DatabaseRoleName.FullyQualifiedName())
	parts = append(parts, strconv.FormatBool(g.WithGrantOption))
	parts = append(parts, strconv.FormatBool(g.AlwaysApply))
	if g.AllPrivileges {
		parts = append(parts, "ALL")
	} else {
		parts = append(parts, strings.Join(g.Privileges, ","))
	}
	parts = append(parts, string(g.Kind))
	parts = append(parts, g.Data.String())
	return strings.Join(parts, helpers.IDDelimiter)
}

type OnDatabaseGrantData struct {
	DatabaseName sdk.AccountObjectIdentifier
}

func (d *OnDatabaseGrantData) String() string {
	return d.DatabaseName.FullyQualifiedName()
}

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

func ParseGrantPrivilegesToDatabaseRoleId(id string) (GrantPrivilegesToDatabaseRoleId, error) {
	var databaseRoleId GrantPrivilegesToDatabaseRoleId

	parts := strings.Split(id, helpers.IDDelimiter)
	if len(parts) < 6 {
		return databaseRoleId, sdk.NewError(`database role identifier should hold at least 5 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_data>"`)
	}

	// TODO: Identifier parsing should be replaced with better version introduced in SNOW-999049.
	// Right now, it's same as sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName, but with error handling.
	databaseRoleNameParts := strings.Split(parts[0], ".")
	if len(databaseRoleNameParts) == 0 ||
		(len(databaseRoleNameParts) == 1 && databaseRoleNameParts[0] == "") ||
		(len(databaseRoleNameParts) == 2 && databaseRoleNameParts[1] == "") ||
		len(databaseRoleNameParts) > 2 {
		return databaseRoleId, sdk.NewError(fmt.Sprintf(`invalid DatabaseRoleName value: %s, should be a fully qualified name of database object <database_name>.<name>`, parts[0]))
	}
	databaseRoleId.DatabaseRoleName = sdk.NewDatabaseObjectIdentifier(
		strings.Trim(databaseRoleNameParts[0], `"`),
		strings.Trim(databaseRoleNameParts[1], `"`),
	)

	if parts[1] != "false" && parts[1] != "true" {
		return databaseRoleId, sdk.NewError(fmt.Sprintf(`invalid WithGrantOption value: %s, should be either "true" or "false"`, parts[1]))
	}
	databaseRoleId.WithGrantOption = parts[1] == "true"

	if parts[2] != "false" && parts[2] != "true" {
		return databaseRoleId, sdk.NewError(fmt.Sprintf(`invalid AlwaysApply value: %s, should be either "true" or "false"`, parts[2]))
	}
	databaseRoleId.AlwaysApply = parts[2] == "true"

	privileges := strings.Split(parts[3], ",")
	if len(privileges) == 0 || (len(privileges) == 1 && privileges[0] == "") {
		return databaseRoleId, sdk.NewError(fmt.Sprintf(`invalid Privileges value: %s, should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all privileges`, parts[3]))
	}
	if len(privileges) == 1 && (privileges[0] == "ALL" || privileges[0] == "ALL PRIVILEGES") {
		databaseRoleId.AllPrivileges = true
	} else {
		databaseRoleId.Privileges = privileges
	}
	databaseRoleId.Kind = DatabaseRoleGrantKind(parts[4])

	switch databaseRoleId.Kind {
	case OnDatabaseDatabaseRoleGrantKind:
		databaseRoleId.Data = &OnDatabaseGrantData{
			DatabaseName: sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[5]),
		}
	case OnSchemaDatabaseRoleGrantKind:
		if len(parts) < 7 {
			return databaseRoleId, sdk.NewError(`database role identifier should hold at least 7 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_on_schema_type>|<on_schema_grant_data>..."`)
		}
		onSchemaGrantData := OnSchemaGrantData{
			Kind: OnSchemaGrantKind(parts[5]),
		}
		switch onSchemaGrantData.Kind {
		case OnSchemaSchemaGrantKind:
			onSchemaGrantData.SchemaName = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[6]))
		case OnAllSchemasInDatabaseSchemaGrantKind, OnFutureSchemasInDatabaseSchemaGrantKind:
			onSchemaGrantData.DatabaseName = sdk.Pointer(sdk.NewAccountObjectIdentifier(parts[6]))
		default:
			return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaGrantKind: %s", onSchemaGrantData.Kind))
		}
		databaseRoleId.Data = &onSchemaGrantData
	case OnSchemaObjectDatabaseRoleGrantKind:
		if len(parts) < 7 {
			return databaseRoleId, sdk.NewError(`database role identifier should hold at least 7 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_on_schema_object_type>|<on_schema_object_grant_data>..."`)
		}
		onSchemaObjectGrantData := OnSchemaObjectGrantData{
			Kind: OnSchemaObjectGrantKind(parts[5]),
		}
		switch onSchemaObjectGrantData.Kind {
		case OnObjectSchemaObjectGrantKind:
			if len(parts) != 8 {
				return databaseRoleId, sdk.NewError(`database role identifier should hold 8 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnObject|<object_type>|<object_name>"`)
			}
			onSchemaObjectGrantData.Object = &sdk.Object{
				ObjectType: sdk.ObjectType(parts[6]),
				Name:       sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(parts[7]),
			}
		case OnAllSchemaObjectGrantKind, OnFutureSchemaObjectGrantKind:
			bulkOperationGrantData := &BulkOperationGrantData{
				ObjectNamePlural: sdk.PluralObjectType(parts[6]),
			}
			if len(parts) > 7 {
				if len(parts) != 9 {
					return databaseRoleId, sdk.NewError(`database role identifier should hold 9 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|On[All or Future]|<object_type_plural>|In[Database or Schema]|<identifier>"`)
				}
				bulkOperationGrantData.Kind = BulkOperationGrantKind(parts[7])
				switch bulkOperationGrantData.Kind {
				case InDatabaseBulkOperationGrantKind:
					bulkOperationGrantData.Database = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[8]))
				case InSchemaBulkOperationGrantKind:
					bulkOperationGrantData.Schema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[8]))
				default:
					return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid BulkOperationGrantKind: %s", bulkOperationGrantData.Kind))
				}
			}
			onSchemaObjectGrantData.OnAllOrFuture = bulkOperationGrantData
		default:
			return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaObjectGrantKind: %s", onSchemaObjectGrantData.Kind))
		}
		databaseRoleId.Data = &onSchemaObjectGrantData
	default:
		return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid DatabaseRoleGrantKind: %s", databaseRoleId.Kind))
	}

	return databaseRoleId, nil
}
