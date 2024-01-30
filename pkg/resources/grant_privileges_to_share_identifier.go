package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ShareGrantKind string

const (
	OnDatabaseShareGrantKind ShareGrantKind = "OnDatabase"
	OnSchemaShareGrantKind   ShareGrantKind = "OnSchema"
	//	TODO(SNOW-1021686): Because function identifier contains arguments which are not supported right now
	// OnFunctionShareGrantKind          ShareGrantKind = "OnFunction"
	OnTableShareGrantKind             ShareGrantKind = "OnTable"
	OnAllTablesInSchemaShareGrantKind ShareGrantKind = "OnAllTablesInSchema"
	OnTagShareGrantKind               ShareGrantKind = "OnTag"
	OnViewShareGrantKind              ShareGrantKind = "OnView"
)

type GrantPrivilegesToShareId struct {
	ShareName  sdk.ExternalObjectIdentifier
	Privileges []string
	Kind       ShareGrantKind
	Identifier sdk.ObjectIdentifier
}

func (id *GrantPrivilegesToShareId) String() string {
	return strings.Join([]string{
		id.ShareName.FullyQualifiedName(),
		strings.Join(id.Privileges, ","),
		string(id.Kind),
		id.Identifier.FullyQualifiedName(),
	}, helpers.IDDelimiter)
}

func ParseGrantPrivilegesToShareId(idString string) (GrantPrivilegesToShareId, error) {
	var grantPrivilegesToShareId GrantPrivilegesToShareId

	parts := strings.Split(idString, helpers.IDDelimiter)
	if len(parts) != 4 {
		return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf(`snowflake_grant_privileges_to_share id is composed out of 4 parts "<account_name>.<share_name>|<privileges>|<grant_on_type>|<grant_on_identifier>", but got %d parts: %v`, len(parts), parts))
	}

	grantPrivilegesToShareId.ShareName = sdk.NewExternalObjectIdentifierFromFullyQualifiedName(parts[0])
	privileges := strings.Split(parts[1], ",")
	if len(privileges) == 0 || (len(privileges) == 1 && privileges[0] == "") {
		return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf(`invalid Privileges value: %s, should be comma separated list of privileges`, privileges))
	}
	grantPrivilegesToShareId.Privileges = privileges
	grantPrivilegesToShareId.Kind = ShareGrantKind(parts[2])

	id, err := helpers.DecodeSnowflakeParameterID(parts[3])
	if err != nil {
		return grantPrivilegesToShareId, err
	}

	switch grantPrivilegesToShareId.Kind {
	case OnDatabaseShareGrantKind:
		if typedIdentifier, ok := id.(sdk.AccountObjectIdentifier); ok {
			grantPrivilegesToShareId.Identifier = typedIdentifier
		} else {
			return grantPrivilegesToShareId, fmt.Errorf(
				"invalid identifier, expected fully qualified name of account object: %s, but instead got: %s",
				getExpectedIdentifierRepresentationFromGeneric[sdk.AccountObjectIdentifier](),
				getExpectedIdentifierRepresentationFromParam(id),
			)
		}
	case OnSchemaShareGrantKind, OnAllTablesInSchemaShareGrantKind:
		if typedIdentifier, ok := id.(sdk.DatabaseObjectIdentifier); ok {
			grantPrivilegesToShareId.Identifier = typedIdentifier
		} else {
			return grantPrivilegesToShareId, fmt.Errorf(
				"invalid identifier, expected fully qualified name of database object: %s, but instead got: %s",
				getExpectedIdentifierRepresentationFromGeneric[sdk.DatabaseObjectIdentifier](),
				getExpectedIdentifierRepresentationFromParam(id),
			)
		}
	case OnTableShareGrantKind, OnViewShareGrantKind, OnTagShareGrantKind: // , OnFunctionShareGrantKind:
		if typedIdentifier, ok := id.(sdk.SchemaObjectIdentifier); ok {
			grantPrivilegesToShareId.Identifier = typedIdentifier
		} else {
			return grantPrivilegesToShareId, fmt.Errorf(
				"invalid identifier, expected fully qualified name of schema object: %s, but instead got: %s",
				getExpectedIdentifierRepresentationFromGeneric[sdk.SchemaObjectIdentifier](),
				getExpectedIdentifierRepresentationFromParam(id),
			)
		}
	default:
		return grantPrivilegesToShareId, fmt.Errorf("unexpected share grant kind: %v", grantPrivilegesToShareId.Kind)
	}

	return grantPrivilegesToShareId, nil
}
