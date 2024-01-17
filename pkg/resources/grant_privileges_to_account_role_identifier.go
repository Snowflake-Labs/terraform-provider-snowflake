package resources

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strconv"
	"strings"
)

type AccountRoleGrantKind string

const (
	OnAccountAccountRoleGrantKind       AccountRoleGrantKind = "OnAccount"
	OnAccountObjectAccountRoleGrantKind AccountRoleGrantKind = "OnAccountObject"
	OnSchemaAccountRoleGrantKind        AccountRoleGrantKind = "OnSchema"
	OnSchemaObjectAccountRoleGrantKind  AccountRoleGrantKind = "OnSchemaObject"
)

type OnAccountGrantData struct {
}

func (d *OnAccountGrantData) String() string {
	return ""
}

type OnAccountObjectGrantData struct {
	ObjectType sdk.ObjectType
	ObjectName sdk.AccountObjectIdentifier
}

func (d *OnAccountObjectGrantData) String() string {
	return strings.Join([]string{
		d.ObjectType.String(),
		d.ObjectName.FullyQualifiedName(),
	}, helpers.IDDelimiter)
}

type GrantPrivilegesToAccountRoleId struct {
	RoleName        sdk.AccountObjectIdentifier
	WithGrantOption bool
	AlwaysApply     bool
	AllPrivileges   bool
	Privileges      []string
	Kind            AccountRoleGrantKind
	Data            fmt.Stringer
}

func (g *GrantPrivilegesToAccountRoleId) String() string {
	var parts []string
	parts = append(parts, g.RoleName.FullyQualifiedName())
	parts = append(parts, strconv.FormatBool(g.WithGrantOption))
	parts = append(parts, strconv.FormatBool(g.AlwaysApply))
	if g.AllPrivileges {
		parts = append(parts, "ALL")
	} else {
		parts = append(parts, strings.Join(g.Privileges, ","))
	}
	parts = append(parts, string(g.Kind))
	data := g.Data.String()
	if len(data) > 0 {
		parts = append(parts, data)
	}
	return strings.Join(parts, helpers.IDDelimiter)
}

func ParseGrantPrivilegesToAccountRoleId(id string) (GrantPrivilegesToAccountRoleId, error) {
	var accountRoleId GrantPrivilegesToAccountRoleId

	parts := strings.Split(id, helpers.IDDelimiter)
	if len(parts) < 5 {
		return accountRoleId, sdk.NewError(`account role identifier should hold at least 5 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>"`)
	}

	// TODO: Identifier parsing should be replaced with better version introduced in SNOW-999049.
	// Right now, it's same as sdk.NewAccountObjectIdentifierFromFullyQualifiedName, but with error handling.
	name := strings.Trim(parts[0], `"`)
	if len(name) == 0 {
		return accountRoleId, sdk.NewError(fmt.Sprintf(`invalid (empty) AccountRoleName value: %s, should be a fully qualified name of account object <name>`, parts[0]))
	}
	accountRoleId.RoleName = sdk.NewAccountObjectIdentifier(name)

	if parts[1] != "false" && parts[1] != "true" {
		return accountRoleId, sdk.NewError(fmt.Sprintf(`invalid WithGrantOption value: %s, should be either "true" or "false"`, parts[1]))
	}
	accountRoleId.WithGrantOption = parts[1] == "true"

	if parts[2] != "false" && parts[2] != "true" {
		return accountRoleId, sdk.NewError(fmt.Sprintf(`invalid AlwaysApply value: %s, should be either "true" or "false"`, parts[2]))
	}
	accountRoleId.AlwaysApply = parts[2] == "true"

	privileges := strings.Split(parts[3], ",")
	if len(privileges) == 0 || (len(privileges) == 1 && privileges[0] == "") {
		return accountRoleId, sdk.NewError(fmt.Sprintf(`invalid Privileges value: %s, should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all privileges`, parts[3]))
	}
	if len(privileges) == 1 && (privileges[0] == "ALL" || privileges[0] == "ALL PRIVILEGES") {
		accountRoleId.AllPrivileges = true
	} else {
		accountRoleId.Privileges = privileges
	}

	accountRoleId.Kind = AccountRoleGrantKind(parts[4])
	switch accountRoleId.Kind {
	case OnAccountAccountRoleGrantKind:
		accountRoleId.Data = new(OnAccountGrantData)
	case OnAccountObjectAccountRoleGrantKind:
		if len(parts) != 7 {
			return accountRoleId, sdk.NewError(`account role identifier should hold at least 7 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnAccountObject|<object_type>|<object_name>"`)
		}
		accountRoleId.Data = &OnAccountObjectGrantData{
			ObjectType: sdk.ObjectType(parts[5]),
			ObjectName: sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[6]),
		}
	case OnSchemaAccountRoleGrantKind:
		if len(parts) < 7 {
			return accountRoleId, sdk.NewError(`account role identifier should hold at least 7 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|<grant_on_schema_type>|<on_schema_grant_data>..."`)
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
			return accountRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaGrantKind: %s", onSchemaGrantData.Kind))
		}
		accountRoleId.Data = &onSchemaGrantData
	case OnSchemaObjectAccountRoleGrantKind:
		if len(parts) < 7 {
			return accountRoleId, sdk.NewError(`account role identifier should hold at least 7 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|<grant_on_schema_object_type>|<on_schema_object_grant_data>..."`)
		}
		onSchemaObjectGrantData := OnSchemaObjectGrantData{
			Kind: OnSchemaObjectGrantKind(parts[5]),
		}
		switch onSchemaObjectGrantData.Kind {
		case OnObjectSchemaObjectGrantKind:
			if len(parts) != 8 {
				return accountRoleId, sdk.NewError(`account role identifier should hold 8 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnObject|<object_type>|<object_name>"`)
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
					return accountRoleId, sdk.NewError(`account role identifier should hold 9 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|On[All or Future]|<object_type_plural>|In[Database or Schema]|<identifier>"`)
				}
				bulkOperationGrantData.Kind = BulkOperationGrantKind(parts[7])
				switch bulkOperationGrantData.Kind {
				case InDatabaseBulkOperationGrantKind:
					bulkOperationGrantData.Database = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[8]))
				case InSchemaBulkOperationGrantKind:
					bulkOperationGrantData.Schema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[8]))
				default:
					return accountRoleId, sdk.NewError(fmt.Sprintf("invalid BulkOperationGrantKind: %s", bulkOperationGrantData.Kind))
				}
			}
			onSchemaObjectGrantData.OnAllOrFuture = bulkOperationGrantData
		default:
			return accountRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaObjectGrantKind: %s", onSchemaObjectGrantData.Kind))
		}
		accountRoleId.Data = &onSchemaObjectGrantData
	default:
		return accountRoleId, sdk.NewError(fmt.Sprintf("invalid AccountRoleGrantKind: %s", accountRoleId.Kind))
	}

	return accountRoleId, nil
}
